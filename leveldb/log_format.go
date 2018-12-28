package leveldb

const (
	kZeroType = iota
	kFullType
	kFirstType
	kMiddleType
	kLastType
)

const kMaxRecordType = kLastType

const kLogBlockSize = 32768

const kHeaderSize = 4 + 2 + 1