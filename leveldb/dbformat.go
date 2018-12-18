package leveldb

import (
	"fmt"
	"encoding/binary"
	"./utilties"
)

// Grouping of constants.  We may want to make some of these
// parameters set via options.
const (
	kNumLevels = 7

	// Level-0 compaction is started when we hit this many files.
	kL0_CompactionTrigger = 4

)

type ValueType uint8
const (
	kTypeDeletion ValueType = 0
	kTypeValue ValueType = 1
)

type sequenceNumber uint64
const kMaxSequenceNumber sequenceNumber = (1 << 56) -1

const kKeyHead = 8

type internalKey struct {
	rep string
}

func makeInternalKey(userKey string, number sequenceNumber, t ValueType) internalKey {
	iKey := new(internalKey)
	
	parsediKey := makeParsedInternalKey(userKey, number, t)

	appendInternalKey(&iKey.rep, &parsediKey)

	return *iKey
}

func clear(internalKey internalKey) {
	internalKey.rep = ""
}

func userKey(internalKey internalKey) string {
	return extractUserKey(internalKey.rep)
}

type parsedInternalKey struct {
	userKey string
	sequence sequenceNumber
	vt ValueType
}

func makeParsedInternalKey(key string, sequence sequenceNumber, vt ValueType) parsedInternalKey {
	return parsedInternalKey{
		userKey: key,
		sequence: sequence,
		vt: vt,
	}
}

func (this *parsedInternalKey) String() string {
	return fmt.Sprintf("'%s' @ %v : %d", utilties.EscapeString(this.userKey), this.sequence, this.vt)
}

func extractUserKey(internalKey string) string {
	// assert(len(internalKey) > 8)
	return internalKey[:len(internalKey) - kKeyHead]
}

func appendInternalKey(result *string, key *parsedInternalKey) {
	*result += key.userKey

	// looks like go had a encoding/binary packge to process these operations

	buf := make([]byte, kKeyHead)
	binary.PutUvarint(buf, packSequenceAndType(uint64(key.sequence), key.vt))

	*result += string(buf)
}

func packSequenceAndType(seq uint64, t ValueType) uint64 {
	// assert(seq <= kMaxSequenceNumber)
	// assert(t <= kValueTypeForSeek)

	return (seq << 8) | uint64(t)
}

func parseInternalKey(internalKey string, result *parsedInternalKey) bool {
	if(len(internalKey) < kKeyHead) {
		return false
	}

	buf := []byte(internalKey[len(internalKey) - kKeyHead : ])

	val, _ := binary.Uvarint(buf)

	c := val & 0xFF

	result.sequence = sequenceNumber(val >> 8)
	result.vt = ValueType(c)
	result.userKey = internalKey[ : len(internalKey) - kKeyHead]

	return ValueType(c) <= kTypeValue
}