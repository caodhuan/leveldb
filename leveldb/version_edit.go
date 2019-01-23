package leveldb

import "strconv"

const (
	kComparator = iota
	kLogNumber
	kNextFileNumber
	kLastSequence
	kCompactPointer
	kDeletedFile
	kNewFile
	// 8 was used for large value refs
	kPrevLogNumber = 9 
)
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

	compactPointers []internalKeyPair
	deletedFiles map[string]deletedFilePair
	newFiles []fileMetaDataPair
}

type internalKeyPair struct {
	level int
	key internalKey
}

type deletedFilePair struct {
	level int
	file uint64
}

type fileMetaDataPair struct {
	level int
	FileMetaData
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
	
	ve.newFiles = make([]fileMetaDataPair, 1)
	ve.compactPointers = make([]internalKeyPair, 1)
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
	this.deletedFiles =  make(map[string]deletedFilePair)
	this.newFiles = this.newFiles[:0]
	this.compactPointers = this.compactPointers[:0]
}

func (this *VersionEdit) SetComparatorName(name string) {
	this.hasComparator = true
	this.comparator = name
}

func (this *VersionEdit) SetLogNumber(num uint64) {
	this.hasLogNumber = true
	this.prevLogNumber = num
}

func (this *VersionEdit) SetPrevLogNumber(num uint64) {
	this.hasPrevLogNumber = true
	this.prevLogNumber = num
}

func (this *VersionEdit) SetNextFile(num uint64) {
	this.hasNextFileNumber = true
	this.nextFileNumber = num
}

func (this *VersionEdit) SetLastSequence(seq sequenceNumber) {
	this.hasLastSequence = true
	this.lastSequence = seq
}

func (this *VersionEdit) SetCompactPointer(level int, key internalKey) {
	this.compactPointers = append(this.compactPointers, internalKeyPair{
		level: level,
		key: key,
	})
}

// Add the specified file at the specified number.
// REQUIRES: This version has not been saved (see Version::SaveTo)
// REQUIRES: "smallest" and "largest" are smallest and largest keys in file
func (this *VersionEdit) AddFile(level int, file uint64, fileSize uint64, smallest *internalKey,
	largest *internalKey) {
	f := FileMetaData {
		number: file,
		fileSize: fileSize,
		smallest: smallest,
		largest: largest,
	}

	this.newFiles = append(this.newFiles, fileMetaDataPair {
		level: level,
		FileMetaData: f,
	})
}

// Delete the specified "file" from the specified "level".
func (this *VersionEdit)DeleteFile(level int, file uint64) {
	mapKey := strconv.Itoa(level) + "_" + strconv.FormatUint(file, 10)
	this.deletedFiles[mapKey] = deletedFilePair{
		level: level,
		file: file,
	}
}

func (this *VersionEdit) EncodeTo(dst* []byte, start int) int {
	if this.hasComparator {
		start += putVarint32(dst, start, kComparator)
		start += putLengthPrefixedSlice(dst, start, this.comparator)
	}

	if this.hasLogNumber {
		start += putVarint32(dst, start, kLogNumber)
		start += putVarint64(dst, start, this.logNumber)
	}

	if this.hasPrevLogNumber {
		start += putVarint32(dst, start, kPrevLogNumber)
		start += putVarint64(dst, start, this.prevLogNumber)
	}

	if this.hasNextFileNumber {
		start += putVarint32(dst, start, kNextFileNumber)
		start += putVarint64(dst, start, this.nextFileNumber)
	}

	if this.hasLastSequence {
		start += putVarint32(dst, start, kLastSequence)
		start += putVarint64(dst, start, uint64(this.lastSequence) )
	}

	for _, v := range this.compactPointers {
		start += putVarint32(dst, start, kCompactPointer)
		start += putVarint32(dst, start, uint32(v.level) )
		start += putLengthPrefixedSlice(dst, start, v.key.encode())
	}

	for _, v := range this.deletedFiles {
		start += putVarint32(dst, start, kDeletedFile)
		start += putVarint32(dst, start, uint32(v.level) )
		start += putVarint64(dst, start, v.file)
	}

	for _, v := range this.newFiles {
		f := v.FileMetaData
		start += putVarint32(dst, start, kNewFile)
		start += putVarint32(dst, start, uint32(v.level) )
		start += putVarint64(dst, start, f.number)
		start += putVarint64(dst, start, f.fileSize)
		start += putLengthPrefixedSlice(dst, start, f.smallest.encode())
		start += putLengthPrefixedSlice(dst, start, f.largest.encode())
	}
}

func (this *VersionEdit) DecodeFrom(src []byte) Status {
	this.Clear()


	// Temporary storage for parsing
	var level int 
	var number uint64
	var f FileMetaData
	var key internalKey

	//for len(src) >= 0 && getVarint32()

	return OK()
}

func (this *VersionEdit) DebugString() string {
	var result string
	
	result += "VersionEdit {"

	if this.hasComparator {
		result += "\n	Comparator: "
		result += this.comparator
	}



	return result
}