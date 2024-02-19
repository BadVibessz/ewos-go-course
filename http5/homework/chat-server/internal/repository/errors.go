package repository

import "errors"

var (
	ErrNoSuchPrivateMessage = errors.New("no such private message")
	ErrNoSuchPublicMessage  = errors.New("no such public message")

	ErrNoSuchUser     = errors.New("no such user")
	ErrEmailExists    = errors.New("user with this email already exists")
	ErrUsernameExists = errors.New("user with this username already exists")
)
