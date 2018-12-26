package leveldb

import (
	"sync"
	"sync/atomic"
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
	internalKeyComparator
	internalFilterPolicy
	options Options
	ownsInfoLog bool
	ownsCache bool
	dbName string

	// table_cache_ provides its own synchronization
	tableCache TableCache
	
	// Lock over the persistent DB state.  Non-NULL iff successfully acquired.
	dbLock FileLock

	mutex sync.Mutex
	shuttingDown atomic.Value
	bgCV sync.Cond
	
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
	s := OK()

	if s.OK() {
		*db = impl
	}

	return s
}

func (this dbImpl) Put(writeOptions WriteOptions, key string, value string) Status {
	s := OK()

	return s
}

func (this dbImpl) Delete(writeOptions WriteOptions, key string) Status {
	s := OK()

	return s
}

func (this dbImpl) Write(writeOptions WriteOptions, updates *WriteBatch) Status {
	return OK()
}

func (this dbImpl) Get(readOptions ReadOptions, key string, value *string) Status {
	return OK()
}

func (this dbImpl) NewIterator( readOptions ReadOptions) *Iterator {
	return nil
}
func (this dbImpl) GetSnapshot() *Snapshot {
	return nil
}

func (this dbImpl) ReleaseSnapshot(snapshot *Snapshot) {

}

func (this dbImpl) GetProperty(property string, value *string) bool {
	return true
}

func (this dbImpl) GetApproximateSizes(kr *keyRange, n int32, sizes *uint64) {

}

func (this dbImpl) CompactRange(begin string, end string) {

}


func makeDBImpl(options *Options, name string) dbImpl {
	return dbImpl{
		env: options.Env,
		internalKeyComparator: makeInternalKeyComparator(options.Comparator),
		internalFilterPolicy: makeInternalFilterPolicy(options.FilterPolicy),
		dbName:name,
	}
}


func sanitizeOptions(dbName string, icmp internalKeyComparator, ipolicy internalFilterPolicy, options *Options) Options {
	result := *options
	result.Comparator = &icmp
	if options.FilterPolicy != nil {
		result.FilterPolicy = &ipolicy
	} else {
		result.FilterPolicy = nil
	}

	ClipToRangeInt32(&result.MaxOpenFiles, 64 + kNumNonTableCacheFiles, 50000)
	ClipToRangeUint32(&result.WriteBufferSize, 64<<10, 1<<30)
	ClipToRangeUint32(&result.BlockSize, 1<<10, 4<<20)

	if result.InfoLog != nil {
		// Open a log file in the same directory as the db
		// options.CreateDir(dbName)
		// options.RenameFile()
	}

	return result
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