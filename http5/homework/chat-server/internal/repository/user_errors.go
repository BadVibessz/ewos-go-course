package repository

import "errors"

var (
	ErrNoSuchUser     = errors.New("no such user")
	ErrEmailExists    = errors.New("user with this email already exists")
	ErrUsernameExists = errors.New("user with this username already exists")
)
