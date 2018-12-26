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

	// Store in result the names of the children of the specified directory.
  	// The names are relative to "dir".
  	// Original contents of results are dropped.
	GetChildren(dir string) ([]string, Status)

	DeleteFile(fname string) Status

	CreateDir(dirname string) Status

	GetFileSize(fname string) (fileSize int64, status Status)

	RenameFile(src string, target string) Status

	// Lock the specified file.  Used to prevent concurrent access to
	// the same db by multiple processes.  On failure, stores NULL in
	// *lock and returns non-OK.
	//
	// On success, stores a pointer to the object that represents the
	// acquired lock in *lock and returns OK.  The caller should call
	// UnlockFile(*lock) to release the lock.  If the process exits,
	// the lock will be automatically released.
	//
	// If somebody else already holds the lock, finishes immediately
	// with a failure.  I.e., this call does not wait for existing locks
	// to go away.
	//
	// May create the named file if it does not already exist.
	LockFile(fname string, lock *FileLock) Status

	// Release the lock acquired by a previous successful call to LockFile.
	// REQUIRES: lock was returned by a successful LockFile() call
	// REQUIRES: lock has not already been unlocked.
	UnlockFile(lock FileLock) Status

	// Arrange to run "f()" once in a background thread.
	//
	// "function" may run in an unspecified thread.  Multiple functions
	// added to the same Env may run concurrently in different threads.
	// I.e., the caller may not assume that background work items are
	// serialized.
	Schedule(f func())

	// Create and return a log file for storing informational messages.
	NewLogger(fname string, result *Logger) Status
	  
	// Returns the number of micro-seconds since some fixed point in time. Only
  	// useful for computing deltas of time.
	NowMicros() uint64
	  
	// Sleep/delay the thread for the prescribed number of micro-seconds.
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
	// f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0755)
	// if err != nil {
	// 	return IOError(fmt.Sprintf("%v", err) )
	// }
	// todo:先省略不实现
	return OK()
}

func (this *defaultFileLock) Unlock() Status {
	// todo: 同上
	return OK()
}

func (this *defaultEnv) NewSequentialFile(fname string, result *SequentialFile) Status {
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		return IOError(fmt.Sprintf("%v", err) )
	}

	*result = &defaultSequentialFile{
		File: f,
	}

	return OK()
}

func (this *defaultEnv) NewRandomAccessFile(fname string, result *RandomAccessFile) Status {
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		return IOError(fmt.Sprintf("%v", err) )
	}

	*result = &defaultRandomAccessFile{
		File: f,
	}

	return OK()
}

func (this *defaultEnv) NewWritableFile(fname string, result *WritableFile) Status {

	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		return IOError(fmt.Sprintf("%v", err) )
	}

	*result = &defaultWritableFile{
		File: f,
	}

	return OK()
}
	  
func (this *defaultEnv) FileExists(fname string) bool {
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		return false
	}
	
	return true
}

func (this *defaultEnv) GetChildren(dir string) ([]string, Status) {
	f, err := os.Open(dir)

	if err != nil {
		return  nil, IOError(fmt.Sprintf("%v", err) )
	}

	names, err := f.Readdirnames(-1)

	if err != nil {
		return  nil, IOError(fmt.Sprintf("%v", err) )
	}

	return names, OK()
}

func (this *defaultEnv) DeleteFile(fname string) Status {
	err := os.Remove(fname)
	if err != nil {
		return  IOError(fmt.Sprintf("%v", err) )
	}

	return OK()
}

func (this *defaultEnv) CreateDir(dirname string) Status {
	if this.FileExists(dirname) {
		return OK()
	}
	err := os.Mkdir(dirname, os.ModePerm)

	if err != nil {
		return IOError(fmt.Sprintf("%v", err) )
	}

	return OK()
}

func (this *defaultEnv) GetFileSize(fname string) (fileSize int64, status Status) {
	f, err := os.Stat(fname)
	if os.IsNotExist(err) {
		return 0, IOError(fmt.Sprintf("%v", err) )
	}
	return f.Size(), OK()
}

func (this *defaultEnv) RenameFile(src string, target string) Status {
	err := os.Rename(src, target)
	if err != nil {
		return IOError(fmt.Sprintf("%v", err) )
	}

	return OK()
}

func (this *defaultEnv) LockFile(fname string, lock *FileLock) Status {
	l := &defaultFileLock{}
	*lock = l

	return (*lock).Lock(fname)
}

func (this *defaultEnv) UnlockFile(lock FileLock) Status {
	return lock.Unlock()
}

func (this *defaultEnv) Schedule(f func()) {
	go f()
}

func (this *defaultEnv) NewLogger(fname string, result *Logger) Status {
	dl := new(defaultLogger)
	
	*result = dl

	return this.NewWritableFile(fname, &(dl.WritableFile))
}

func (this *defaultEnv) NowMicros() uint64 {
	return uint64(time.Now().UnixNano() / 1000)
}

func (this *defaultEnv) SleepForMicroseconds(micros uint32) {
	nanoseconds := 1000 * time.Duration(micros)
	
	time.Sleep(nanoseconds)
	
} 