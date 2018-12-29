package leveldb


type FileMetaData struct {
	// Seeks allowed until compaction
	allowedSeeks int
	number uint64
	fileSize uint64
	smallest *internalKey
	largest *internalKey
}

func newFileMetaData() *FileMetaData {
	return &FileMetaData{
		allowedSeeks: 1 << 30,
		fileSize: 0,
	}
}

type VersionEdit struct {
	
}

