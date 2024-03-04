package entity

import "time"

type PublicMessage struct {
	ID       int
	From     *User // TODO: CHANGE TO FromID
	Content  string
	SentAt   time.Time
	EditedAt time.Time
}
