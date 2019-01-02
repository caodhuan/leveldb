package leveldb

import "unsafe"

type Block struct {
	data []byte
	size uint
	restartOffset uint32
	owned bool
}

func (this *Block) NewIterator(comparator Comparator) Iterator {
	var s uint32
	if (this.size < uint(unsafe.Sizeof(s)) ) {
		return NewErrorIterator(Corruption("bad block contents"))
	}
	numRestarts := 1 // NumRestarts();

	if (numRestarts == 0) {
		return NewEmptyIterator()
	} else {
		// new Iter(cmp, data_, restart_offset_, num_restarts);
	}

	return NewEmptyIterator()
}