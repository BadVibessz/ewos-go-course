package model

import "time"

type PrivateMessage struct {
	ID       int
	From     *User
	To       *User
	Content  string
	SentAt   time.Time
	EditedAt time.Time
}

type PublicMessage struct {
	ID       int
	From     *User
	Content  string
	SentAt   time.Time
	EditedAt time.Time
}
