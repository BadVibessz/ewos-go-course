package entity

import "time"

type PrivateMessage struct {
	ID           int
	FromUsername string
	ToUsername   string
	Content      string
	SentAt       time.Time
	EditedAt     time.Time
}
