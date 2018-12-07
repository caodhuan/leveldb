package leveldb

import "fmt"

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
	env *Env
	internalComparator  InternalKeyComparator
	internalFilterPolicy InternalFilterPolicy
	options Options
	ownsInfoLog bool
	ownsCache bool
	dbName string

	tableCache TableCache
}

type keyRange struct {
	start string
	limit string
}

func init() {
	fmt.Printf("Hello world!\n")
}

func Open(options Options, name string, db *DB) Status {
	*db = nil
	impl := &dbImpl {
		options: options,
	}
	s := OK()

	if s.OK() {
		*db = impl
	}
	var defaultE defaultEnv
	if defaultE == DefaultEnv() {
		fmt.Printf("Hello world!``1111\n")
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