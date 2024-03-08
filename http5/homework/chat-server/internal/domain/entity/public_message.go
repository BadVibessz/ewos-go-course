package entity

import "time"

type PublicMessage struct {
	ID       int
	From     *User
	Content  string
	SentAt   time.Time
	EditedAt time.Time
}
