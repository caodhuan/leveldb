package leveldb

import (
	"time"
)

type Env interface {
	// Create a brand new sequentially-readable file with the specified name.
  	// On success, stores a pointer to the new file in *result and returns OK.
  	// On failure stores NULL in *result and returns non-OK.  If the file does
  	// not exist, returns a non-OK status.
  	//
  	// The returned file will only be accessed by one thread at a time.
	NewSequentialFile(fname string, result *SequentialFile) Status

	// Create a brand new random access read-only file with the
  	// specified name.  On success, stores a pointer to the new file in
  	// *result and returns OK.  On failure stores NULL in *result and
  	// returns non-OK.  If the file does not exist, returns a non-OK
  	// status.
  	//
  	// The returned file may be concurrently accessed by multiple threads.
	NewRandomAccessFile(fname string, result *RandomAccessFile) Status

	// Create an object that writes to a new file with the specified
	// name.  Deletes any existing file with the same name and creates a
	// new file.  On success, stores a pointer to the new file in
	// *result and returns OK.  On failure stores NULL in *result and
	// returns non-OK.
	//
	// The returned file will only be accessed by one thread at a time.
	NewWritableFile(fname string, result *WritableFile) Status
	  
	FileExists(fname string) bool

	GetChildren(dir string, result []string) Status

	DeleteFile(fname string) Status

	CreateDir(dirname string) Status

	GetFileSize(fnmae string, fileSize uint64) Status

	RenameFile(src string, target string) Status

	LockFile(fname string, lock *FileLock) Status

	UnlockFile(lock FileLock) Status

	Schedule(f func())

	NewLogger(fname string, result *Logger) Status

	NowMicros() uint64

	SleepForMicroseconds(micros uint32) 
}


type defaultEnv struct {
}

func DefaultEnv() Env {
	return &defaultEnv{}
}

type SequentialFile interface {

}

type RandomAccessFile interface {
	// Read up to "n" bytes from the file starting at "offset".
	// "scratch[0..n-1]" may be written by this routine.  Sets "*result"
	// to the data that was read (including if fewer than "n" bytes were
	// successfully read).  May set "*result" to point at data in
	// "scratch[0..n-1]", so "scratch[0..n-1]" must be live when
	// "*result" is used.  If an error was encountered, returns a non-OK
	// status.
	//
	// Safe for concurrent use by multiple threads.
	Read(offset uint64, n uint64, result *string, scratch string) Status
}

type WritableFile interface {
	Append(data string) Status 
	Close() Status 
	Flush() Status 
	Sync() Status 
}

type FileLock interface {

}

func (this *defaultEnv) NewSequentialFile(fname string, result *SequentialFile) Status {
	return OK()
}

func (this *defaultEnv) NewRandomAccessFile(fname string, result *RandomAccessFile) Status {
	return OK()
}

func (this *defaultEnv) NewWritableFile(fname string, result *WritableFile) Status {
	return OK()
}
	  
func (this *defaultEnv) FileExists(fname string) bool {
	return true
}

func (this *defaultEnv) GetChildren(dir string, result []string) Status {
	return OK()
}

func (this *defaultEnv) DeleteFile(fname string) Status {
	return OK()
}

func (this *defaultEnv) CreateDir(dirname string) Status {
	return OK()
}

func (this *defaultEnv) GetFileSize(fnmae string, fileSize uint64) Status {
	return OK()
}

func (this *defaultEnv) RenameFile(src string, target string) Status {
	return OK()
}

func (this *defaultEnv) LockFile(fname string, lock *FileLock) Status {
	return OK()
}

func (this *defaultEnv) UnlockFile(lock FileLock) Status {
	return OK()
}

func (this *defaultEnv) Schedule(f func()) {

}

func (this *defaultEnv) NewLogger(fname string, result *Logger) Status {
	return OK()
}

func (this *defaultEnv) NowMicros() uint64 {
	return uint64(time.Now().UnixNano() / 1000)
}

func (this *defaultEnv) SleepForMicroseconds(micros uint32) {
	nanoseconds := 1000 * time.Duration(micros)
	
	time.Sleep(nanoseconds)
	
} 

type defaultWritableFile struct {
	
}