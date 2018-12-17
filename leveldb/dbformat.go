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

type internalKey struct {
	rep string
}

func makeInternalKey(userKey string, number sequenceNumber, t ValueType) internalKey {
	return internalKey {
		rep: userKey,
	}
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
	return internalKey[8:]
}

func appendInternalKey(result *string, key *parsedInternalKey) {
	*result += key.userKey
	// looks like go had a encoding/binary packge to process these operations
	// todo
}