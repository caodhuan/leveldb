package leveldb

type Status struct {
	// OK status has a NULL state_.  Otherwise, state_ is a new[] array
  	// of the following form:
  	//    state_[0..3] == length of message
  	//    state_[4]    == code
  	//    state_[5..]  == message
	state int32
}

const(
	kOK	= iota
	kNotFound
	kCorruption
	kNotSupported
	kInvalidArgument
	kIOError
)

func (this Status) OK() bool {
	return this.state == kOK
}

func (this Status) IsNotFound() bool {
	return this.state == kNotFound
}
