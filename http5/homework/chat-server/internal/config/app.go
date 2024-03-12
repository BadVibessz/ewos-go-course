package config

type Config struct {
	Server
	Jwt
	DB string
	InMemoryDB
	Postgres
}
