package entity

import "time"

type PublicMessage struct {
	ID           int       `db:"id"`
	FromUsername string    `db:"from_username"`
	Content      string    `db:"content"`
	SentAt       time.Time `db:"sent_at"`
	EditedAt     time.Time `db:"edited_at"`
}
