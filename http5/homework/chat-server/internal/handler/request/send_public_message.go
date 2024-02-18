package request

import "github.com/go-playground/validator/v10"

type SendPublicMessageRequest struct {
	FromID  int    `json:"from_id" validate:"required,min=1"`
	Content string `json:"content" validate:"required,min=1,max=2000"`
}

func (sm *SendPublicMessageRequest) Validate() error {
	return validator.New(validator.WithRequiredStructEnabled()).Struct(sm)
}
