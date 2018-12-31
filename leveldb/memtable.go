package leveldb

import "./structure"


type keyComparator struct{
	internalKeyComparator
}

type MemTable struct {
	table structure.SkipList
	arena Arena
}

func (this *keyComparator) Compare(aKey string, bKey string) int {
	aKey = getlengthPrefixedSlice(aKey)
	bKey = getlengthPrefixedSlice(bKey)
	return this.internalKeyComparator.Compare(aKey, bKey)
}


func getlengthPrefixedSlice(str string) string {
	remainStr, _ := getVarint32(str)

	return remainStr
}