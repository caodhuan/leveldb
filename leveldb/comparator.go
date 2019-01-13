package leveldb

import (
	"strings"
	"./utilties"
)

type Comparator interface {
	// Three-way comparison.  Returns value:
  	//   < 0 iff "a" < "b",
  	//   == 0 iff "a" == "b",
  	//   > 0 iff "a" > "b"
	Compare(aKey string, bKey string) int
	Name() string
	
	// Advanced functions: these are used to reduce the space requirements
  	// for internal data structures like index blocks.

  	// If *start < limit, changes *start to a short string in [start,limit).
  	// Simple comparator implementations may return with *start unchanged,
  	// i.e., an implementation of this method that does nothing is correct.
	FindShortestSeparator(start *string, limit string)

	// Changes *key to a short string >= *key.
  	// Simple comparator implementations may return with *key unchanged,
  	// i.e., an implementation of this method that does nothing is correct.
	FindShortSuccessor(key *string)
}

type internalKeyComparator struct {
	Comparator
}

func BytewiseComparator() Comparator {
	return &bytewiseComparator{}
}


type bytewiseComparator struct {
}

func (this *bytewiseComparator) Compare(aKey string, bKey string) int {
	return strings.Compare(aKey, bKey)
}

func (this *bytewiseComparator) Name() string {
	return "leveldb.BytewiseComparator"
}

func (this *bytewiseComparator) FindShortestSeparator(start *string, limit string) {
	minLength := utilties.Min(len(*start), len(limit))
	diffIndex := 0
	for((diffIndex < minLength) && (*start)[diffIndex] == limit[diffIndex]) {
		diffIndex += 1
	}

	if diffIndex >= minLength {

	} else {
		diffByte := (*start)[diffIndex]
		if (diffByte < 0xFF) && (diffByte + 1) < limit[diffIndex] {
			char := diffByte + 1
			*start = (*start)[0 : diffIndex]
			*start += string(char +1)

			//assert(strings.Compare(start, limit) < 0)
		}
	}
}

func (this *bytewiseComparator) FindShortSuccessor(key *string) {
	for i, c := range *key {
		if c != 0xFF {
			(*key) = (*key)[0: i]
			(*key) += string(c + 1)
			break
		}
	}
}

func (this *internalKeyComparator) Name() string {
	return "leveldb.InternalKeyComparator"
}

func (this *internalKeyComparator) Compare(aKey string, bKey string) int {
	result := this.Comparator.Compare(extractUserKey(aKey), extractUserKey(bKey) )

	// Order by:
  	//    increasing user key (according to user-supplied comparator)
  	//    decreasing sequence number
  	//    decreasing type (though sequence# should be enough to disambiguate)
	if result == 0 {
		anum := decodeFixed64(aKey[ len(aKey) - kKeyHead: ])
		bnum := decodeFixed64(bKey[ len(bKey) - kKeyHead: ])

		if anum > bnum {
			result = -1
		} else if (anum < bnum) {
			result = +1
		}
	}

	return result
}

func (this *internalKeyComparator) FindShortestSeparator(start *string, limit string) {
	userStart := extractUserKey(*start)
	userLimit := extractUserKey(limit)

	tmp := userStart

	this.Comparator.FindShortestSeparator(&tmp, userLimit)

	if (len(tmp) < len(userStart) ) && (this.Comparator.Compare(userStart, tmp) < 0) {
		// User key has become shorter physically, but larger logically.
		// Tack on the earliest possible number to the shortened user key.
		b := make([]byte, 8)
		encodeFixed64(b, packSequenceAndType(uint64(kMaxSequenceNumber), kValueTypeForSeek) )
	
		*start = tmp + string(b)
	}
}

func (this *internalKeyComparator) FindShortSuccessor(key *string) {
	userKey := extractUserKey(*key)
	tmp := userKey

	this.Comparator.FindShortSuccessor(&tmp)

	if (len(tmp) < len(userKey)) && (this.Comparator.Compare(userKey, tmp) < 0) {
		// User key has become shorter physically, but larger logically.
		// Tack on the earliest possible number to the shortened user key.
		b := make([]byte, 8)
		encodeFixed64(b, packSequenceAndType(uint64(kMaxSequenceNumber), kValueTypeForSeek) )
	
		*key = tmp + string(b)
	}
}

func (this *internalKeyComparator) userComparator() Comparator {
	return this.Comparator
}

func makeInternalKeyComparator(c Comparator) *internalKeyComparator {
	return &internalKeyComparator {
		Comparator: c,
	}
}