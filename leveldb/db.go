package leveldb

import "fmt"

type DB interface {
	Put(writeOptions WriteOptions, key string, value string) Status
	Delete(writeOptions WriteOptions, key string) Status
	Write(writeOptions WriteOptions, updates *WriteBatch) Status
	Get(readOptions ReadOptions, key string, value *string) Status
	NewIterator( readOptions ReadOptions)   *Iterator
	GetSnapshot() *Snapshot
	ReleaseSnapshot(snapshot *Snapshot)
	GetProperty(property string, value *string) bool
	GetApproximateSizes(kr *keyRange, n int32, sizes *uint64)
	CompactRange(begin string, end string)
}

type dbImpl struct {
}

type keyRange struct {
	start string
	limit string
}

func init() {
	fmt.Printf("Hello world!\n")
}

func Open(options Options, name string, db **DB) Status {
	*db = nil
	
	s := OK()

	return s
}