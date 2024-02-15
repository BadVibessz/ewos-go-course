package request

import "github.com/go-playground/validator/v10"

type RegisterRequest struct {
	Username        string `json:"username" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,max=128"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

func (rr *RegisterRequest) Validate() error {
	return validator.New(validator.WithRequiredStructEnabled()).Struct(rr)
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=1"`
}

func (lr *LoginRequest) Validate() error {
	return validator.New(validator.WithRequiredStructEnabled()).Struct(lr)
}
