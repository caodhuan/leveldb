package leveldb

type TableBuilder struct {
	options *Options
	indexBlockOptions *Options
	file *WritableFile
	offset uint64
	s Status
}

// Create a builder that will store the contents of the table it is
// building in *file.  Does not close the file.  It is up to the
// caller to close the file after calling Finish().
func newTableBuilder(options *Options, file *WritableFile) *TableBuilder {
	return &TableBuilder{}
}

