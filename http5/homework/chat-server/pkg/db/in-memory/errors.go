package in_memory

import "errors"

var (
	ErrNotExistedRow   = errors.New("no such row")
	ErrNotExistedTable = errors.New("no such table")
	ErrExistingKey     = errors.New("key already exists")
)
