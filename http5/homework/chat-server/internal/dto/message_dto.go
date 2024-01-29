package dto

type CreatePrivateMessageDTO struct {
	FromID  int
	ToID    int
	Content string
}

type CreatePublicMessageDTO struct {
	FromID  int
	Content string
}
