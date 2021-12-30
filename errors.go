package yafsm

import "errors"

var (
	ErrEventDuplicated      = errors.New("event duplicated")
	ErrEventIllegal         = errors.New("event illegal")
	ErrEventNotExist        = errors.New("event does not exist")
	ErrStateNotExist        = errors.New("state does not exist")
	ErrIllegalStateForEvent = errors.New("illegal state for event")
)
