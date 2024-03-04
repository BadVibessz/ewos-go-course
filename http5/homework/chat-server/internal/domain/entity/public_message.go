package entity

import "time"

type PublicMessage struct {
	ID           int
	FromUsername string
	Content      string
	SentAt       time.Time
	EditedAt     time.Time
}
