package leveldb

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

func (this *defaultEnv) NewSequentialFile(fname string, result *SequentialFile) Status {
	return OK()
}

func (this *defaultEnv) NewRandomAccessFile(fname string, result *RandomAccessFile) Status {
	return OK()
}