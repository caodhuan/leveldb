package leveldb

import (
	"sync"
	//"sync/atomic"
)

const kNumNonTableCacheFiles = 10

type DB interface {
	Put(writeOptions WriteOptions, key string, value string) Status
	Delete(writeOptions WriteOptions, key string) Status
	Write(writeOptions WriteOptions, updates *WriteBatch) Status
	Get(readOptions ReadOptions, key string, value *string) Status
	NewIterator(readOptions ReadOptions) *Iterator
	GetSnapshot() *Snapshot
	ReleaseSnapshot(snapshot *Snapshot)
	GetProperty(property string, value *string) bool
	GetApproximateSizes(kr *keyRange, n int32, sizes *uint64)
	CompactRange(begin string, end string)
}

type dbImpl struct {
	env Env
	*internalKeyComparator
	*internalFilterPolicy
	options *Options
	ownsInfoLog bool
	ownsCache bool
	dbName string

	// table_cache_ provides its own synchronization
	tableCache *TableCache
	
	// Lock over the persistent DB state.  Non-NULL iff successfully acquired.
	dbLock *FileLock

	// State below is protected by mutex
	mutex sync.Mutex
	shuttingDown int32
	bgCV *sync.Cond  	// Signalled when background work finishes
	mem *MemTable	
	imm *MemTable 		// Memtable being compacted
	hasImm int32	// So bg thread can detect non-NULL imm
	logFile *WritableFile
	logFileNumber uint64
	log *LogWriter
	seed uint32 // for sampling

	// Queue of writers.
	writers []*Writer
	tmpBatch *WriteBatch

	snapshots *SnapshotList

	pendingOutputs map[uint64]bool

	bgCompactionScheduled bool

	manualCompaction *ManualCompaction

	versions *VersionSet

	bgError *Status

	status [kNumLevels]*CompactionStats
}

// Information for a manual compaction
type ManualCompaction struct {
	level int
	done bool
	begin *internalKey	 // NULL means beginning of key range
	end *internalKey	 // NULL means end of key range
	tmpStore internalKey	// Used to keep track of compaction progress
}

type Writer struct {	

}

type CompactionStats struct {
	micros int64
	bytesRead int64
	bytesWritten int64
}

func (this *CompactionStats) Add(other *CompactionStats) {
	this.micros += other.micros
	this.bytesRead += other.bytesRead
	this.bytesWritten += other.bytesWritten
}


type keyRange struct {
	start string
	limit string
}

func init() {
}

func Open(options *Options, name string, db *DB) Status {
	*db = nil
	impl := makeDBImpl(options, name)

	impl.mutex.Lock()
	edit := newVersionEdit()

	s := impl.recover(edit)

	if s.OK() {
		newLogNumber := impl.versions.NewFileNumber()
		var lfile WritableFile
		s = options.Env.NewWritableFile(LogFileName(name, newLogNumber), &lfile)

		if s.OK() {
			edit.SetLogNumber(newLogNumber)
			impl.logFile = &lfile
			impl.log = newLogWriter(&lfile)

			s = impl.versions.LogAndApply(edit, &impl.mutex)
		}

		if s.OK() {
			//impl.DeleteObsoleteFiles()
			//impl.MaybeScheduleCompaction()
		}
	}

	impl.mutex.Unlock()

	if s.OK() {
		*db = impl
	} 

	return s
}

func (this *dbImpl) Put(writeOptions WriteOptions, key string, value string) Status {
	s := OK()

	return s
}

func (this *dbImpl) Delete(writeOptions WriteOptions, key string) Status {
	s := OK()

	return s
}

func (this *dbImpl) Write(writeOptions WriteOptions, updates *WriteBatch) Status {
	return OK()
}

func (this *dbImpl) Get(readOptions ReadOptions, key string, value *string) Status {
	return OK()
}

func (this *dbImpl) NewIterator( readOptions ReadOptions) *Iterator {
	return nil
}
func (this *dbImpl) GetSnapshot() *Snapshot {
	return nil
}

func (this *dbImpl) ReleaseSnapshot(snapshot *Snapshot) {

}

func (this *dbImpl) GetProperty(property string, value *string) bool {
	return true
}

func (this *dbImpl) GetApproximateSizes(kr *keyRange, n int32, sizes *uint64) {

}

func (this *dbImpl) CompactRange(begin string, end string) {

}


func makeDBImpl(options *Options, name string) *dbImpl {
	var impl dbImpl

	impl.env = options.Env
	impl.internalKeyComparator = makeInternalKeyComparator(options.Comparator)
	impl.internalFilterPolicy = makeInternalFilterPolicy(options.FilterPolicy)
	impl.options = sanitizeOptions(name, *impl.internalKeyComparator, *impl.internalFilterPolicy, options)
	impl.dbName = name
	impl.ownsInfoLog = options.InfoLog != impl.options.InfoLog
	impl.ownsCache = options.BlockCache != impl.options.BlockCache
	impl.dbLock = nil
	impl.shuttingDown = 0
	impl.bgCV = sync.NewCond(&impl.mutex)
	impl.mem = newMemTable(*impl.internalKeyComparator)
	impl.imm = nil
	impl.logFile = nil
	impl.logFileNumber = 0
	impl.log = nil
	impl.seed = 0
	impl.tmpBatch = newWriteBatch()
	impl.bgCompactionScheduled = false
	impl.manualCompaction = nil
	impl.hasImm = 0
	
	tableCacheSize := impl.options.MaxOpenFiles - kNumNonTableCacheFiles
	impl.tableCache = newTableCache(impl.dbName, impl.options, int(tableCacheSize) )

	impl.versions = newVersionSet(impl.dbName, impl.options, impl.tableCache, impl.internalKeyComparator)

	return &impl
}


func sanitizeOptions(dbName string, icmp internalKeyComparator, ipolicy internalFilterPolicy, options *Options) *Options {
	result := *options
	result.Comparator = &icmp
	if options.FilterPolicy != nil {
		result.FilterPolicy = &ipolicy
	} else {
		result.FilterPolicy = nil
	}

	ClipToRangeInt32(&result.MaxOpenFiles, 64 + kNumNonTableCacheFiles, 50000)
	ClipToRangeUint32(&result.WriteBufferSize, 64<<10, 1<<30)
	ClipToRangeUint(&result.BlockSize, 1<<10, 4<<20)

	if result.InfoLog != nil {
		// Open a log file in the same directory as the db
		options.CreateDir(dbName)
		options.RenameFile(InfoLogFileName(dbName), OldInfoLogFileName(dbName))
		s := options.NewLogger(InfoLogFileName(dbName), &result.InfoLog)

		if !s.OK() {
			// No place suitable for logging
			result.InfoLog = nil
		}
	}

	if result.BlockCache == nil {
		result.BlockCache = NewLRUCache( 8 << 20)
	}

	return &result
}

func (this *dbImpl) recover(edit *VersionEdit) Status {
	return OK()
}

func ClipToRangeUint64(value *uint64, minValue uint64, maxValue uint64) {
	if *value > maxValue {
		*value = maxValue
	}
	if *value < minValue {
		*value = minValue
	}
}

func ClipToRangeInt64(value *int64, minValue int64, maxValue int64) {
	if *value > maxValue {
		*value = maxValue
	}
	if *value < minValue {
		*value = minValue
	}
}

func ClipToRangeUint32(value *uint32, minValue uint32, maxValue uint32) {
	if *value > maxValue {
		*value = maxValue
	}
	if *value < minValue {
		*value = minValue
	}
}

func ClipToRangeInt32(value *int32, minValue int32, maxValue int32) {
	if *value > maxValue {
		*value = maxValue
	}
	if *value < minValue {
		*value = minValue
	}
}

func ClipToRangeUint(value *uint, minValue uint, maxValue uint) {
	if *value > maxValue {
		*value = maxValue
	}
	if *value < minValue {
		*value = minValue
	}
}