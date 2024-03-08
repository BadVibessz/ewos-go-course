package entity

import "time"

type User struct {
	ID             int
	Email          string
	Username       string
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (u *User) Equal(other User) bool {
	return u.Username == other.Username
}
