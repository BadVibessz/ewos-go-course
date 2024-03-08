package request

import "errors"

var (
	ErrInvalidOffset = errors.New("invalid offset provided")
	ErrInvalidLimit  = errors.New("invalid limit provided")
)
