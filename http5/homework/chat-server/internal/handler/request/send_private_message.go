package request

import "github.com/go-playground/validator/v10"

type SendPrivateMessageRequest struct {
	ToUsername string `json:"to_username" validate:"required,min=1"`
	Content    string `json:"content" validate:"required,min=1,max=2000"`
}

func (sm *SendPrivateMessageRequest) Validate(valid *validator.Validate) error {
	return valid.Struct(sm)
}
