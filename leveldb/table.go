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
	table, _ := arg.(*Table)
	blockCache := table.options.BlockCache
	var (block Block
		 Handle interface{}
		 handle BlockHandle
		)
	input := indexValue

	s := handle.DecodeFrom([]byte(input) )
	// We intentionally allow extra stuff in index_value so that we
	// can add more features in the future.
	if s.OK() {
		// var contents BlockContents
		// if blockCache != nil {
		// 	var cacheKeyBuffer [16]byte
		// }
	}
	// todo

	return NewEmptyIterator()
}

func (this *Table) NewIterator(readOptions *ReadOptions) Iterator {
	return NewTwoLevelIterator(this.indexBlock.NewIterator(this.options.Comparator), BlockReader, this, readOptions)
}