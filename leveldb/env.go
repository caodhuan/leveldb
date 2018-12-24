package leveldb

import (
	"time"
	"os"
	"fmt"
	//"syscall"
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
	
	Read(scratch []byte) Status

	// Skip "n" bytes from the file. This is guaranteed to be no
	// slower that reading the same data, but may be faster.
	//
	// If end of file is reached, skipping will stop at the end of the
	// file, and Skip will return OK.
	//
	// REQUIRES: External synchronization
	Skip(n int64) Status
}

type defaultSequentialFile struct {
	*os.File
}

func (this *defaultSequentialFile) Read(scratch []byte) Status {
	
	_, err := this.File.Read(scratch)
	if err != nil {
		return IOError(fmt.Sprintf("%v", err) )
	}

	return OK()
}

func (this *defaultSequentialFile) Skip(n int64) Status {
	_, err := this.File.Seek(n, 1)

	if err != nil {
		return IOError(fmt.Sprintf("%v", err) )
	}

	return OK()
}

type RandomAccessFile interface {
	// Read up to "n" bytes from the file starting at "offset".
	// "scratch[0..n-1]" may be written by this routine.  Sets "*result"
	// to the data that was read (including if fewer than "n" bytes were
	// successfully read).  May set "*result" to point at data in
	// "scratch[0..n-1]", so "scratch[0..n-1]" must be live when
	// "*result" is used.  If an error was encountered, returns a non-OK
	// status.
	Read(offset int64, scratch []byte) Status
}

type defaultRandomAccessFile struct {
	*os.File
}

func (this *defaultRandomAccessFile) Read(offset int64, scratch []byte) Status {
	_, err := this.File.ReadAt(scratch, offset)

	if err != nil {
		if err != nil {
			return IOError(fmt.Sprintf("%v", err) )
		}
	}

	return OK()
}


type WritableFile interface {
	Append(data string) Status 
	Close() Status 
	Flush() Status 
	Sync() Status 
}

type defaultWritableFile struct {
	*os.File
}

func (this *defaultWritableFile) Append(data string) Status {
	_, err := this.File.Write([]byte(data))
	if err != nil {
		return IOError(fmt.Sprintf("%v", err) )
	}

	return OK()
}

func (this *defaultWritableFile) Close() Status {
	err := this.File.Close()

	if err != nil {
		return IOError(fmt.Sprintf("%v", err) )
	}

	return OK()
}

func (this *defaultWritableFile) Flush() Status {
	return OK()
}

func (this *defaultWritableFile) Sync() Status {
	err := this.File.Sync()

	if err != nil {
		return IOError(fmt.Sprintf("%v", err) )
	}

	return OK()
}

type FileLock interface {
	Lock(fname string) Status
	Unlock() Status
}

type defaultFileLock struct {
	*os.File
}

func (this *defaultFileLock) Lock(fname string) Status {
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return IOError(fmt.Sprintf("%v", err) )
	}
	
	return OK()
}

func (this *defaultFileLock) Unlock() Status {
	return OK()
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