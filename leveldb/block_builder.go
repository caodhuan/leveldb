package leveldb

import (
	"unsafe"
	"encoding/binary"
	"./utilties"
	"fmt"
)

type BlockBuilder struct {
	options *Options
	buffer []byte	// Destination buffer
	restarts []int // Restart points
	counter int	// Number of entries emitted since restart
	finished bool // has Finish() been called?
	lastKey string
}

func newBlockBuilder(options *Options) *BlockBuilder {
	return &BlockBuilder{
		options: options,
		counter: 0,
		finished: false,
		lastKey: "",
	}
}

// Reset the contents as if the BlockBuilder was just constructed.
func (this *BlockBuilder) Reset() {
	this.buffer = this.buffer[:0]
	this.restarts = this.restarts[ :0]
	this.counter = 0
	this.finished = false
	this.lastKey = ""
}

// REQUIRES: Finish() has not been called since tha last call to Reset().
// REQUIRES: key is larger than any previously added key
func (this *BlockBuilder) Add(key string, value string) {
	lastKeyPiece := this.lastKey
	shared := 0
	if this.counter < this.options.BlockRestartInterval {
		// See how much sharing to do with previous string
		minLength := utilties.Min(len(lastKeyPiece), len(key))
		for (shared < minLength) && (lastKeyPiece[shared] == key[shared]) {
			shared++
		}
	} else {
		// Restart compression
		this.restarts = append(this.restarts, len(this.buffer) )
		this.counter = 0
	}

	nonShared := len(key) - shared

	// Add "shared nonShared valueSize" to buffer
	tmpBuffer := make([]byte, 4 * 3)

	binary.PutUvarint(tmpBuffer[0 * 4: ], uint64(shared) )
	binary.PutUvarint(tmpBuffer[1 * 4: ], uint64(nonShared) )
	binary.PutUvarint(tmpBuffer[2 * 4: ], uint64(len(value)) )

	this.buffer = append(this.buffer, tmpBuffer ...)
	
	// Add string delta to buffer followed by value
	this.buffer = append(this.buffer, []byte(key[shared:nonShared]) ...) 
	this.buffer = append(this.buffer, []byte(value) ...)
	
	// update status
	this.lastKey = this.lastKey[:shared]
	this.lastKey += key[shared: nonShared]

	if (this.lastKey != key) {
		panic( fmt.Sprintf("inequal key: %s %s\n", this.lastKey, key ) )
	}
	
	this.counter++
}

// Finish building the block and return a []byte that refers to the
// block contents.
func (this *BlockBuilder) Finish() []byte {
	if this.restarts == nil {
		return nil
	}

	if this.finished {
		return nil
	}

	restartNum := len(this.restarts)

	tmpBuffer := make([]byte, restartNum + 1)

	for index, key := range this.restarts {
		binary.PutUvarint(tmpBuffer[index* 4: ], uint64(key) )
	}

	binary.PutUvarint(tmpBuffer[restartNum * 4 : ], uint64(restartNum) )

	this.finished = true

	return append(this.buffer, tmpBuffer ...)
}

// Returns an estimate of the current size (uncompressed) of the block
// we are building
func (this *BlockBuilder) CurrentSizeEstimate() uint {
	return uint(len(this.buffer)) + uint(len(this.restarts) * 4 ) + uint(unsafe.Sizeof(this.counter) )
}

// Return true if no entries have been added since the last Reset()
func (this *BlockBuilder) Empty() bool {
	return len(this.buffer) == 0
}
