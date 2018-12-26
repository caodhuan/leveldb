package leveldb

import "./structure"

type keyComparator struct{
	internalKeyComparator
}


type Table structure.SkipList
type MemTable struct {
	table Table
	arena Arena
}

func (this *keyComparator) Compare(aKey string, bKey string) int {

}


func getlengthPrefixedSlice(str string) string {
	
}