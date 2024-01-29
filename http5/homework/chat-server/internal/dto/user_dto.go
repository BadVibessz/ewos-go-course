package dto

type CreateUserDTO struct {
	Email          string
	Username       string
	HashedPassword string
}

type UpdateUserDTO struct {
	NewEmail          string
	NewUsername       string
	NewHashedPassword string
}
