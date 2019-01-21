package leveldb

import "./utilties"

type LogWriter struct {
	dest *WritableFile
	blockOffset int

	typeCRC [kMaxRecordType + 1]uint32
}

func newLogWriter(dst *WritableFile) *LogWriter {
	var result LogWriter
	result.dest = dst
	result.blockOffset = 0

	a := make([]byte, 1)
	for i:= 0; i <= kMaxRecordType; i += 1 {
		t := byte(i)
		a[0] = t
		result.typeCRC[i] = utilties.Value(a)
	}

	return &result
}