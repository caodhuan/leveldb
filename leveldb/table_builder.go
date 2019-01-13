package leveldb

import "./utilties"

type TableBuilder struct {
	options *Options
	indexBlockOptions *Options
	file WritableFile
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

	compressedOutput []byte
}

// Create a builder that will store the contents of the table it is
// building in *file.  Does not close the file.  It is up to the
// caller to close the file after calling Finish().
func newTableBuilder(options *Options, file WritableFile) *TableBuilder {

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

// Change the options used by this builder. Note: only some of the
// option fields can be changed after construction. If a field is
// not allowed to be change dynamically and its value in the structure
// passed to the constructor is different from its value in the structure
// passed to this method, this method will return an error without
// changing any fields
func (this *TableBuilder) ChangeOptions(options *Options) Status {

	// Note: if more fields are added to Options, update this function to
	// catch changes that should not be allowed to change in the middle of
	// building a Table.
	if options.Comparator != this.options.Comparator {
		return InvalidArgument("changing comparator while building table")
	} 

	this.options = &(*options)
	this.indexBlockOptions = &(*options)
	this.indexBlockOptions.BlockRestartInterval = 1

	return OK()
}

// Add key, value to the table being constructed.Add
// REQUIRES: key is after any previously added key according to comparator.
// REQUIRES: Finish(), Abandon() has not been called
func (this *TableBuilder) Add(key string, value string) {
	if (!this.ok()) {
		return
	}
	keyByte := []byte(key)

	if (this.numEntries > 0) {

	}

	if (this.pendingIndexEntry) {
		this.options.Comparator.FindShortestSeparator(&this.lastKey, key)
		var handleEncoding []byte
		this.pendingHandle.EncodeTo(&handleEncoding)
		this.indexBlock.Add(keyByte, handleEncoding)
		this.pendingIndexEntry = false
	}

	if (this.filterBlock != nil) {
		this.filterBlock.AddKey(keyByte)
	}

	this.lastKey = key
	this.numEntries++
	this.dataBlock.Add(keyByte, []byte(value) )
	
	estimatedBlockSize := this.dataBlock.CurrentSizeEstimate()

	if (estimatedBlockSize >= this.options.BlockSize) {
		this.Flush()
	}
}

// Advanced operation: flush any buffered key/value pairs to file.
// Can be used to ensure that two adjacent entries never live in
// the same data block.  Most clients should not need to use this method.
// REQUIRES: Finish(), Abandon() have not been called
func (this *TableBuilder) Flush() {
	if (!this.ok()) {
		return
	}

	if (this.dataBlock.Empty() ) {
		return 
	}

	this.writeBlock(this.dataBlock, this.pendingHandle)

	if (this.ok() ) {
		this.pendingIndexEntry = true
		this.s = this.file.Flush()
	}

	if (this.filterBlock != nil) {
		this.filterBlock.StartBlock(this.offset)
	}

}

func (this *TableBuilder) ok() bool {
	return this.s.OK()
}

func (this *TableBuilder) writeBlock(block *BlockBuilder, handle *BlockHandle) {
	// File format contains a sequence of blocks where each block has:
	//    block_data: uint8[n]
	//    type: uint8
	//    crc: uint32
	raw := block.Finish()

	var blockContents []byte

	compressionType := this.options.Compression

	switch compressionType {
	case NoCompression:
		blockContents = raw
	case SnappyCompression:
		// this would be done in the future~~~
		// for now, just do the same as NoCompression 
		blockContents = raw
	}

	this.writeRawBlock(blockContents, compressionType, handle)

	this.compressedOutput = this.compressedOutput[ :0]
	
	block.Reset()
}

func (this *TableBuilder) writeRawBlock(blockContents []byte, compressionType CompressionType, handle *BlockHandle) {
	handle.offset = this.offset
	handle.size = uint64(len(blockContents))
	this.s = this.file.Append(blockContents)

	if (this.s.OK()) {
		trailer := make([]byte, kBlockTrailerSize)
		trailer = append(trailer, byte(compressionType))
		crc := utilties.Value(blockContents)
		crc = utilties.Extend(crc, trailer)
		encodeFixed32(trailer[1:], utilties.Mask(crc))
		this.s = this.file.Append(trailer)

		if this.s.OK() {
			this.offset += uint64(len(blockContents) + len(trailer))
		}
	}
}
