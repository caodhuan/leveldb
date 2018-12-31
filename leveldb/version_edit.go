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


type FileMetaDataSort struct {
	fileMetaData []*FileMetaData
	less func(a, b *FileMetaData) bool
}

func (this *FileMetaDataSort) Len() int {
	return len(this.fileMetaData)
}

func (this *FileMetaDataSort) Less(i, j int) bool {
	return this.less(this.fileMetaData[i], this.fileMetaData[j])
}

func (this *FileMetaDataSort) Swap(i, j int) {
	this.fileMetaData[i], this.fileMetaData[j] = this.fileMetaData[j], this.fileMetaData[i]
}