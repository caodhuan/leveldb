package leveldb

// An iterator yields a sequence of key/value pairs from a source.
// The following class defines the interface.  Multiple implementations
// are provided by this library.  In particular, iterators are provided
// to access the contents of a Table or a DB.
//
// Multiple threads can invoke const methods on an Iterator without
// external synchronization, but if any of the threads may call a
// non-const method, all threads accessing the same Iterator must use
// external synchronization.
type Iterator interface {
 	// An iterator is either positioned at a key/value pair, or
  	// not valid.  This method returns true iff the iterator is valid.
	Valid() bool

	// Position at the first key in the source.  The iterator is Valid()
  	// after this call iff the source is not empty.
	SeekToFirst()
	  
	// Position at the last key in the source.  The iterator is
	// Valid() after this call iff the source is not empty.
	SeekToLast()

  	// Position at the first key in the source that at or past target
  	// The iterator is Valid() after this call iff the source contains
  	// an entry that comes at or past target.
	Seek(target string)

	// Moves to the next entry in the source.  After this call, Valid() is
  	// true iff the iterator was not positioned at the last entry in the source.
  	// REQUIRES: Valid()
	Next()

	// Moves to the previous entry in the source.  After this call, Valid() is
  	// true iff the iterator was not positioned at the first entry in source.
  	// REQUIRES: Valid()
	Prev()

	// Return the key for the current entry.  The underlying storage for
  	// the returned slice is valid only until the next modification of
  	// the iterator.
  	// REQUIRES: Valid()
	Key() string

	// Return the value for the current entry.  The underlying storage for
  	// the returned slice is valid only until the next modification of
  	// the iterator.
  	// REQUIRES: Valid()
	Value() string

	// If an error has occurred, return it.  Else return an ok status.
	Status() Status
	
}

type EmptyIterator struct {
	s Status
}

func (this *EmptyIterator) Valid() bool {
	return false
}

func (this *EmptyIterator) SeekToFirst() {
}

func (this *EmptyIterator) SeekToLast() {
}

func (this *EmptyIterator) Seek(target string) {
}

func (this *EmptyIterator) Next() {
}

func (this *EmptyIterator) Prev() {
}

func (this *EmptyIterator) Key() string {
	return ""
}

func (this *EmptyIterator) Value() string {
	return ""
}

func (this *EmptyIterator) Status() Status {
	return this.s
}


type TwoLevelIterator struct {
	blockFunction func(arg interface{}, options *ReadOptions, indexValue string) Iterator
	arg *interface{}
	options *ReadOptions
	s Status
	
	indexIter *IteratorWrapper
	dataIter *IteratorWrapper

	// If dataIter is non-NULL, then "dataBlockHandle" holds the
  	// "index_value" passed to blockFunction to create the dataIter.
	dataBlockHandle string
}

func (this *TwoLevelIterator) Valid() bool {
	return false
}

func (this *TwoLevelIterator) SeekToFirst() {
}

func (this *TwoLevelIterator) SeekToLast() {
}

func (this *TwoLevelIterator) Seek(target string) {
}

func (this *TwoLevelIterator) Next() {
}

func (this *TwoLevelIterator) Prev() {
}

func (this *TwoLevelIterator) Key() string {
	return ""
}

func (this *TwoLevelIterator) Value() string {
	return ""
}

func (this *TwoLevelIterator) Status() Status {
	return OK()
}


func NewTwoLevelIterator(indexIter Iterator, blockFunction func(arg interface{}, options *ReadOptions, indexValue string) Iterator, arg *interface{}, options *ReadOptions ) Iterator {
	return &TwoLevelIterator{
		blockFunction: blockFunction,
		arg: arg,
		options: options,
		indexIter: IteratorToIteratorWrapper(indexIter),
		dataIter: nil,
	}
}

func NewEmptyIterator() Iterator {
	return &EmptyIterator{
		s: OK(),
	}
}

func NewErrorIterator(s Status) Iterator {
	return &EmptyIterator{
		s: s,
	}
}

func IteratorToIteratorWrapper(iter Iterator) *IteratorWrapper {
	var result IteratorWrapper
	result.Set(iter)
	return &result
}

type IteratorWrapper struct {
	iter Iterator
	valid bool
	key string
}

func (this *IteratorWrapper) Set(iter Iterator) {
	this.iter = iter
	if this.iter == nil {
		this.valid = false
	} else {
		this.Update()
	}
}

func (this *IteratorWrapper) Update() {
	this.valid = this.Valid()
	if this.valid {
		this.key = this.Key()
	}
}

func (this *IteratorWrapper) Key() string {
	return this.key
}

func (this *IteratorWrapper) Valid() bool {
	return this.valid
}
