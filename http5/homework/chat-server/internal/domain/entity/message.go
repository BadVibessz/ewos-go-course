package entity

type PrivateMessage struct {
	FromID  int
	ToID    int
	Content string
}

type PublicMessage struct {
	FromID  int
	Content string
}
