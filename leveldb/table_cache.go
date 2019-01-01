package leveldb

type TableCache struct {
	Env
	dbName string
	options *Options
	Cache
}

func newTableCache(dbName string, options *Options, entries int) *TableCache {
	return &TableCache{
		Env: options.Env,
		dbName: dbName,
		options: options,
		Cache: NewLRUCache(entries),
	}
}

func (this *TableCache) NewIterator(options ReadOptions, fileNumber uint64, fileSize uint64, tablePtr **Table) Iterator {
	if  tablePtr != nil {
		*tablePtr = nil
	}

	var handle interface{}

	s := this.FindTable(fileNumber, fileSize, &handle)
	if !s.OK() {
		return NewErrorIterator(s)
	}

	tableAndFile, ok := (*this.Cache.Value(&handle)).(TableAndFile)
	table := tableAndFile.table
	result := table.NewIterator()
}

func (this *TableCache) FindTable(fileNumber, fileSize uint64, handle *interface{}) Status {

	return OK()
}

type TableAndFile struct {
	file *RandomAccessFile
	table *Table
}