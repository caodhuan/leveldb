package leveldb

import (
	"strings"
	"./utilties"
)

type Compartor interface {
	Compare( a string, b string) int
	Name() string
	FindShortestSeparator(start *string, limit string)
	FindShortSuccessor(key *string)
}

type InternalKeyComparator struct {

}

func BytewiseComparator() Compartor {
	return &bytewiseComparator{}
}


type bytewiseComparator struct {
}

func (this *bytewiseComparator) Compare(a string, b string) int {
	return strings.Compare(a, b)
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