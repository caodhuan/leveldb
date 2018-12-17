package leveldb

import "fmt"

type Status struct {
	code uint8
	msg string
}

const(
	kOK	= iota
	kNotFound
	kCorruption
	kNotSupported
	kInvalidArgument
	kIOError
)

func (this *Status) OK() bool {
	return this.code == kOK
}

func (this *Status) IsNotFound() bool {
	return this.code == kNotFound
}

func (this *Status) IsCorruption() bool {
	return this.code == kCorruption
}

func (this *Status) IsIOError() bool {
	return this.code == kIOError
}

func (this *Status) String() string {
	var sCode string
	switch this.code {
	case kOK:
		sCode = "OK"
	case kNotFound:
		sCode = "NotFound: "
	case kCorruption:
		sCode = "Corrupution: "
	case kNotSupported:
		sCode = "Not implemented: "
	case kInvalidArgument:
		sCode = "Invalid argument: "
	case kIOError:
		sCode = "IO error: "
	default:
		sCode = fmt.Sprintf("unknow code (%v): ", this.code)
	}

	return sCode + this.msg
}

func makeStatus(code uint8, msg string) Status {
	return Status {
		code: code,
		msg: msg,
	}
}

func OK() Status {
	return makeStatus(kOK, "")
}

func NotFound(msg string) Status {
	return makeStatus(kNotFound, msg)
}

func Corruption(msg string) Status {
	return makeStatus(kCorruption, msg)
}

func NotSupported(msg string) Status {
	return makeStatus(kNotSupported, msg)
}

func InvalidArgument(msg string) Status {
	return makeStatus(kInvalidArgument, msg)
}

func IOError(msg string) Status {
	return makeStatus(kIOError, msg)
}
