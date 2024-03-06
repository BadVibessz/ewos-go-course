package entity

import "time"

type User struct {
	ID             int       `db:"id"`
	Email          string    `db:"email"`
	Username       string    `db:"username"`
	HashedPassword string    `db:"hashed_password"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func (u *User) Equal(other User) bool { return u.Username == other.Username }
