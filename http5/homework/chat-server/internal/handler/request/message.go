package request

type SendPublicMessageRequest struct {
	Content string `json:"content" validate:"required,min=1,max=2000"`
}

type SendPrivateMessageRequest struct {
	ToID    int    `json:"to_id"`
	Content string `json:"content" validate:"required,min=1,max=2000"`
}
