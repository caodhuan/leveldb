package leveldb

import (
	"sync"
)

// LRU cache implementation

const kNumShardBits = 4

const kNumShards = 1 << kNumShardBits

type Cache interface {
	Value(handle *interface{}) *interface{}
}


// An entry is a variable length heap-allocated structure.  Entries
// are kept in a circular doubly linked list ordered by access time.
type LRUHandle struct {
	value *interface{}
	deleter func(string, *interface{})
	nextHash *LRUHandle
	next *LRUHandle
	prev *LRUHandle
	charge uint
	hash uint32
	keyData []byte
}

func (this *LRUHandle) key() string {
	if this.next == this {
		value, _ := (*this.value).(string)
		return value
	}

	return string(this.keyData)
}

type HandleTable struct {
	// The table consists of an array of buckets where each bucket is
  	// a linked list of cache entries that hash into the bucket.
	length uint32
	elems uint32
	list **LRUHandle

}

type LRUCache struct {
	capacity uint
	mutex *sync.Mutex

	lru LRUHandle

	table HandleTable
}

func (this *LRUCache) SetCapacity(capacity uint) {
	this.capacity = capacity
}


type ShardedLRUCache struct {
	shared [kNumShards]*LRUCache
	idMutex *sync.Mutex
	lastID	uint64
}

func (this *ShardedLRUCache) Value(handle *interface{}) *interface{} {
	lruHandle, ok := (*handle).(LRUHandle)
	if !ok {
		return lruHandle.value
	}

	return nil
}

func NewLRUCache(capacity int) Cache {
	var result ShardedLRUCache

	perShared := (capacity + (kNumShards - 1) ) / kNumShards

	for i := 0; i < kNumShards; i++ {
		result.shared[i].SetCapacity(uint(perShared) )
	}

	return result
}

