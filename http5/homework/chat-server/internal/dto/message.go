package dto

type PrivateMessageDTO struct {
	FromID  int
	ToID    int
	Content string
}

type PublicMessageDTO struct {
	FromID  int
	Content string
}
