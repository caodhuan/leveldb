package leveldb

// A Table is a sorted map from strings to strings.  Tables are
// immutable and persistent.  A Table may be safely accessed from
// multiple threads without external synchronization.
type Table struct {
	options *Options
	s Status
	file *RandomAccessFile
	filter *FilterBlockReader
	
}