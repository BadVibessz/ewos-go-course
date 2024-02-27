package repository

import "errors"

var (
	ErrNoSuchPrivateMessage = errors.New("no such private message")
	ErrNoSuchPublicMessage  = errors.New("no such public message")
)
