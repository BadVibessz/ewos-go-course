package request

import "github.com/go-playground/validator/v10"

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=1"`
}

func (lr *LoginRequest) Validate() error {
	return validator.New(validator.WithRequiredStructEnabled()).Struct(lr)
}
