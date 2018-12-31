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

func (this *TableCache) NewIterator() Iterator {

	
}