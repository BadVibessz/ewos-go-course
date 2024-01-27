package models

import "time"

type User struct { // todo: add fields!
	ID             int
	Email          string
	Username       string
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreateUserModel struct {
	Email          string
	Username       string
	HashedPassword string
}

type UpdateUserModel struct {
	NewEmail          string
	NewUsername       string
	NewHashedPassword string
}
