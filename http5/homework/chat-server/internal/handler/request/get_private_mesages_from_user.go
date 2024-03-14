package request

import "github.com/go-playground/validator/v10"

type GetPrivateMessagesFromUserRequest struct {
	FromUsername string `json:"from_username" validate:"required,min=1"`
}

func (sm *GetPrivateMessagesFromUserRequest) Validate(valid *validator.Validate) error {
	return valid.Struct(sm)
}
