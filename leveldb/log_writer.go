package leveldb

type LogWriter struct {
	dest *WritableFile
	blockOffset int
}