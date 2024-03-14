package request

import "github.com/go-playground/validator/v10"

type SendPublicMessageRequest struct {
	Content string `json:"content" validate:"required,min=1,max=2000"`
}

func (sm *SendPublicMessageRequest) Validate(valid *validator.Validate) error {
	return valid.Struct(sm)
}
