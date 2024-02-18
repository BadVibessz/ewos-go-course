package request

import (
	"github.com/go-playground/validator/v10"
	"math"
)

type PaginationOptions struct {
	Offset int `json:"offset" validate:"required,min=0"`
	Limit  int `json:"limit" validate:"required,min=0"`
}

func (po *PaginationOptions) Validate() error {
	return validator.New(validator.WithRequiredStructEnabled()).Struct(po) // TODO:
}

func GetUnlimitedPaginationOptions() PaginationOptions {
	return PaginationOptions{
		Offset: 0,
		Limit:  math.MaxInt64,
	}
}
