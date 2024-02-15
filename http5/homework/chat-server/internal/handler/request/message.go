package request

import "github.com/go-playground/validator/v10"

type SendPublicMessageRequest struct {
	FromID  int    `json:"from_id" validate:"required,min=1"`
	Content string `json:"content" validate:"required,min=1,max=2000"`
}

func (sm *SendPublicMessageRequest) Validate() error {
	return validator.New(validator.WithRequiredStructEnabled()).Struct(sm) // todo: maybe global var?
}

type SendPrivateMessageRequest struct {
	ToID    int    `json:"to_id" validate:"required,min=1"`
	FromID  int    `json:"from_id" validate:"required,min=1"`
	Content string `json:"content" validate:"required,min=1,max=2000"`
}

func (sm *SendPrivateMessageRequest) Validate() error {
	return validator.New(validator.WithRequiredStructEnabled()).Struct(sm)
}

type PaginationOptions struct {
	Offset int
	Limit  int
}

func NewPaginationOptions(offset, limit int) PaginationOptions {
	return PaginationOptions{
		Offset: offset,
		Limit:  limit,
	}
}

func (po *PaginationOptions) Validate() error {
	if po.Offset < 0 {
		return ErrInvalidOffset
	}

	if po.Limit < 0 {
		return ErrInvalidLimit
	}

	return nil
}
