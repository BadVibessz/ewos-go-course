package model

import "time"

type User struct {
	ID             int
	Email          string
	Username       string
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
