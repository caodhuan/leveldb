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

// kValueTypeForSeek defines the ValueType that should be passed when
// constructing a ParsedInternalKey object for seeking to a particular
// sequence number (since we sort sequence numbers in decreasing order
// and the value type is embedded as the low 8 bits in the sequence
// number in internal keys, we need to use the highest-numbered
// ValueType, not the lowest).
const kValueTypeForSeek = kTypeValue

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

func (this *internalKey) clear() {
	this.rep = ""
}

func (this *internalKey) userKey() string {
	return extractUserKey(this.rep)
}

func (this *internalKey) decodeFrom(s string) {
	this.rep = s
}

func (this *internalKey) encode() string {
	return this.rep
}

func (this *internalKey) setFrom(p *parsedInternalKey) {
	this.rep = ""

	appendInternalKey(&this.rep, p)
}

type LookupKey struct {
	// We construct a char array of the form:
	//    klength  varint32               <-- 0
	//    userkey  char[klength]          <-- kStart
	//    tag      uint64
	// The array is a suitable MemTable key.
	// The suffix starting with "userkey" can be used as an InternalKey.
	kStart int
	space []byte
}

func newLookupKey(userKey string, sequence sequenceNumber) *LookupKey {
	var lookupKey LookupKey

	userKeyByte := []byte(userKey)

	lookupKey.space = make([]byte, len(userKeyByte) + kKeyHead +  (kKeyHead / 2) )

	lookupKey.kStart = encodeVarint32(lookupKey.space, uint32(len(userKey) + kKeyHead) )
	
	userKeyLen := copy(lookupKey.space[lookupKey.kStart:], userKeyByte)

	encodeFixed64(lookupKey.space[lookupKey.kStart + userKeyLen:], packSequenceAndType(uint64(sequence), kValueTypeForSeek) )

	return &lookupKey
}

// Return a key suitable for lookup in a MemTable.
func (this *LookupKey) memtableKey() string {
	return string(this.space)
}

// Return an internal key (suitable for passing to an internal iterator)
func (this *LookupKey) internalKey() string {
	return string(this.space[ this.kStart : ] )
}

// Return the user key
func (this *LookupKey) userKey() string {
	return string(this.space[ this.kStart : len(this.space) - 8] )
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

	b := make([]byte, 8)
	// looks like go had a encoding/binary packge to process these operations
	encodeFixed64(b, packSequenceAndType(uint64(key.sequence), key.vt))

	*result += string(b)
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

	val := decodeFixed64(internalKey[len(internalKey) - kKeyHead : ]) 

	c := val & 0xFF

	result.sequence = sequenceNumber(val >> 8)
	result.vt = ValueType(c)
	result.userKey = internalKey[ : len(internalKey) - kKeyHead]

	return ValueType(c) <= kTypeValue
}

func encodeFixed64(buf []byte, value uint64) int {
	binary.PutUvarint(buf, value)

	return kKeyHead
}

func encodeFixed32(buf []byte, value uint32) int {
	binary.PutUvarint(buf, uint64(value) )
	return kKeyHead / 2
}

func encodeVarint64(buf []byte, value uint64) int {
	return binary.PutUvarint(buf, value)
}

func encodeVarint32(buf []byte, value uint32) int {
	return binary.PutUvarint(buf, uint64(value) )
}

func decodeFixed64(value string) uint64 {
	if len(value) != kKeyHead {
		return 0
	}

	buf := []byte(value)

	result, _ := binary.Uvarint(buf)

	return result
}

func decodeFix32(value string) uint32 {
	if len(value) != (kKeyHead / 2) {
		return 0
	}

	buf := []byte(value)

	result, _ := binary.Uvarint(buf)

	return uint32(result)
}

// Standard Get... routines parse a value from the beginning of a Slice
// and advance the slice past the parsed value.
func getVarint64(value string) (string, uint64) {
	if len(value) != kKeyHead {
		return value, 0
	}

	buf := []byte(value)

	result, l := binary.Uvarint(buf)

	return value[l:], result
}

// Standard Get... routines parse a value from the beginning of a Slice
// and advance the slice past the parsed value.
func getVarint32(value string) (string, uint32) {
	if len(value) != (kKeyHead / 2) {
		return value, 0
	}

	buf := []byte(value)

	result, l := binary.Uvarint(buf)

	return value[l:], uint32(result)
}
