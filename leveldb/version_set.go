package leveldb

import (
	"sort"
)

type Version struct {
	vSet *VersionSet // VersionSet to which this Version belongs
	next *Version	// Next version in linked list
	prev *Version	// Previous version in linked list

	files [kNumLevels][]*FileMetaData  // List of files per level

	// Next file to compact based on seek stats.
	fileToCompact *FileMetaData
	fileToCompactLevel int
  
	// Level that should be compacted next and its compaction score.
  	// Score < 1 means compaction is not strictly needed.  These fields
  	// are initialized by Finalize().
	compactionScore float64
	CompactionLevel	int
}

type VersionSet struct {
	env *Env
	dbName string
	options *Options
	tableCache *TableCache
	icmp *internalKeyComparator
	nextFileNumber uint64
	mainfestFileNumber uint64
	lastSequence uint64
	logNumber uint64
	prevLogNumber uint64 // 0 or backing store for memtable being compacted

	// opened lazily
	descriptorFile *WritableFile
	descriptorLog *LogWriter
	dummyVersions *Version // Head of circular doubly-linked list of versions.
	current *Version	// == dummy_versions_.prev_

	// Per-level key at which the next compaction at that level should start.
	// Either an empty string, or a valid InternalKey.
	CompactPointer [kNumLevels]string
}

// Append to *iters a sequence of iterators that will
// yield the contents of this Version when merged together.
// REQUIRES: This version has been saved (see VersionSet::SaveTo)
func (this *Version) AddIterators(readOptions *ReadOptions, iters []*Iterator) {
	// for key, value := range this.files[0] {

	// }
}

func (this *Version) Get(readOptions *ReadOptions, key LookupKey, val *string) (seekFile *FileMetaData, seekFileLevel int) {
	iKey := key.internalKey()
	userKey := key.userKey()

	ucmp := this.vSet.icmp.userComparator()

	seekFile = nil
	seekFileLevel = -1

	var lastFileRead *FileMetaData
	var lastFileReadLevel = -1

	var tmp []*FileMetaData
	var tmp2 *FileMetaData
	  
	// We can search level-by-level since entries never hop across
  	// levels.  Therefore we are guaranteed that if we find data
  	// in an smaller level, later levels are irrelevant.
	for level := 0; level < kNumLevels; level++ {
		numFiles := len(this.files[level])
		if numFiles == 0 {
			continue
		}

		// Get the list of files to search in this level
		files := this.files[level]
		if level == 0 {
			tmp = make([]*FileMetaData, numFiles)[:0]

			for i := 0; i < numFiles; i++ {
				f := files[i]
				if (ucmp.Compare(userKey, f.smallest.userKey() ) >= 0 &&
					ucmp.Compare(userKey, f.largest.userKey() ) <= 0) {
					
					tmp = append(tmp, f)
				}
			}

			if len(tmp) == 0 {
				continue
			}

		}

	}
	
	return seekFile, seekFileLevel
}