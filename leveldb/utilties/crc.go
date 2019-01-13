package utilties

import (
	"hash/crc32"
)

var table = crc32.MakeTable(crc32.Castagnoli)

// provide exactly the same semantics of leveldb's crc32 interface

func Value(data []byte) uint32 {
	return crc32.Update(0, table, data)
}

func Extend(initCRC uint32, data []byte) uint32 {
	return crc32.Update(initCRC, table, data)
}

const kMaskDelta = 0xa282ead8
// Return a masked representation of crc.
//
// Motivation: it is problematic to compute the CRC of a string that
// contains embedded CRCs.  Therefore we recommend that CRCs stored
// somewhere (e.g., in files) should be masked before being stored.
func Mask(crc uint32) uint32 {
	 // Rotate right by 15 bits and add a constant.
	return ((crc >> 15) | (crc << 17)) + kMaskDelta;
}

// Return the crc whose masked representation is masked_crc.
func Unmask(crc uint32) uint32 {
	rot := crc - kMaskDelta;
	return ((rot >> 17) | (rot << 15));
}

