package leveldb

type TableBuilder struct {
	options *Options
	indexBlockOptions *Options
	file *WritableFile
	offset uint64
	s Status
	dataBlock *BlockBuilder
	indexBlock *BlockBuilder
	lastKey string
	numEntries int64
	closed bool						// Either Finish() or Abandon() has been called.
	filterBlock *FilterBlockBuilder
	
	// We do not emit the index entry for a block until we have seen the
	// first key for the next data block.  This allows us to use shorter
	// keys in the index block.  For example, consider a block boundary
	// between the keys "the quick brown fox" and "the who".  We can use
	// "the r" as the key for the index block entry since it is >= all
	// entries in the first block and < all entries in subsequent
	// blocks.
	//
	// Invariant: r->pending_index_entry is true only if data_block is empty.
	pendingIndexEntry bool
	pendingHandle *BlockHandle	 // Handle to add to index block
}

// Create a builder that will store the contents of the table it is
// building in *file.  Does not close the file.  It is up to the
// caller to close the file after calling Finish().
func newTableBuilder(options *Options, file *WritableFile) *TableBuilder {

	copyOptions := *options

	result := TableBuilder{
		options: options,
		indexBlockOptions: &copyOptions,
		file: file,
		offset: 0,
		dataBlock: newBlockBuilder(options),
		indexBlock: newBlockBuilder(&copyOptions),
		numEntries: 0,
		closed: false,
		pendingIndexEntry: false,
	}

	if result.options.FilterPolicy == nil {
		result.filterBlock = nil
	} else {
		result.filterBlock = newFliterBlockBuilder(result.options.FilterPolicy)
	}

	result.indexBlockOptions.BlockRestartInterval = 1
	
	if result.filterBlock != nil {
		result.filterBlock.StartBlock(0)
	}

	return &result
}

