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
	comparator string
	logNumber uint64
	prevLogNumber uint64
	nextFileNumber uint64
	lastSequence sequenceNumber
	hasComparator bool
	hasLogNumber bool
	hasPrevLogNumber bool
	hasNextFileNumber bool
	hasLastSequence bool

	compactPointers []map[int]internalKey
	deletedFiles map[int]map[int]uint64
	newFiles []map[int]FileMetaData
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

func newVersionEdit() *VersionEdit {
	var ve VersionEdit
	
	ve.newFiles = make([]map[int]FileMetaData, 1)
	ve.compactPointers = make([]map[int]internalKey, 1)
	ve.Clear()
	
	return &ve
}

func (this *VersionEdit) Clear() {
	this.comparator = ""
	this.logNumber = 0
	this.prevLogNumber = 0
	this.lastSequence = 0
	this.nextFileNumber = 0
	this.hasComparator = false
	this.hasLogNumber = false
	this.hasPrevLogNumber = false
	this.hasNextFileNumber = false
	this.deletedFiles =  make(map[int]map[int]uint64)
	this.newFiles = this.newFiles[:0]
	this.compactPointers = this.compactPointers[:0]
}

func (this *VersionEdit) SetLogNumber(num uint64) {
	this.hasLogNumber = true
	this.prevLogNumber = num
}