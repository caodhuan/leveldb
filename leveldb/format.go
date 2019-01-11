package leveldb

import "encoding/binary"
// Maximum encoding length of a BlockHandle
const kMaxEncodedLength = 10 + 10

type BlockContents struct {
	data []byte
	cachable bool
	heapAllocated bool
}

type BlockHandle struct {
	offset uint64
	size uint64
}

func (this *BlockHandle) EncodeTo(dst *[]byte) {
	tmp := make([]byte, 8)

	l := binary.PutUvarint(tmp, this.offset)
	*dst = append(*dst, tmp[:l] ...)

	l = binary.PutUvarint(tmp[0:], this.size)
	*dst = append(*dst, tmp[:l] ...)
	
}

func (this *BlockHandle) DecodeFrom(input []byte) Status {
	value, n := binary.Uvarint(input)
	if n > 0 {
		this.offset = value
	} else {
		return Corruption("bad block handle")
	}

	value, n = binary.Uvarint(input[n:])

	if n > 0 {
		this.size = value
	} else {
		Corruption("bad block handle")
	}

	return OK()	
}