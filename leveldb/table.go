package leveldb

// A Table is a sorted map from strings to strings.  Tables are
// immutable and persistent.  A Table may be safely accessed from
// multiple threads without external synchronization.
type Table struct {
	options *Options
	s Status
	file *RandomAccessFile
	filter *FilterBlockReader
	filterData []byte
	metaIndexHandle *BlockHandle
	indexBlock *Block
}

func BlockReader(arg interface{}, options *ReadOptions, indexValue string) Iterator {
	// todo
	return NewEmptyIterator()
}

func (this *Table) NewIterator(readOptions *ReadOptions) Iterator {
	return NewTwoLevelIterator(this.indexBlock.NewIterator(this.options.Comparator), BlockReader, this, readOptions)
}