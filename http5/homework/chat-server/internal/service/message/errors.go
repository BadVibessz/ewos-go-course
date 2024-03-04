package message

import "errors"

var (
	ErrNoSuchReceiver = errors.New("no such receiver")
	ErrNoSuchSender   = errors.New("no such sender")
)
