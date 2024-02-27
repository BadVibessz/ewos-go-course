package entity

type PrivateMessage struct {
	PublicMessage
	To *User
}
