package leveldb

import "bytes"

const kBlockSize = 4096

type Arena struct {
	// unassigned
	allocBlock []byte

	bytes.Buffer
	// Bytes of memory in blocks allocated so far
	blocksMemory int

}

func NewArena() *Arena {
	return &Arena{
		allocBlock: nil,
		blocksMemory: 0,
	}
}

func (this *Arena) Allocate(size int) []byte {
	if (size < len(this.allocBlock) ) {
		result := this.allocBlock[:size]
		this.allocBlock = this.allocBlock[size:]

		return result
	}
	
	return this.allocateFallback(size)
}

func (this *Arena) AllocateAligned(size int) []byte {
	return this.allocateFallback(size)
}

func (this *Arena) MemoryUsage() int {
	return this.blocksMemory
}

func (this *Arena) allocateFallback(size int) []byte {
	if (size > kBlockSize / 4) {
		// Object is more than a quarter of our block size.  Allocate it separately
		// to avoid wasting too much space in leftover bytes.
		result := this.allocateNewBlock(size)
		return result
	}

	// We waste the remaining space in the current block.
	this.allocBlock = this.allocateNewBlock(kBlockSize)

	result := this.allocBlock[:size]

	this.allocBlock = this.allocBlock[size:]

	return result
}



func (this *Arena) allocateNewBlock(size int) []byte {
	result := make([]byte, size)
	this.blocksMemory += size

	return result
}

