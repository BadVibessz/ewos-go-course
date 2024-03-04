// nolint
package postgres

import (
	"errors"
	inmemory "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/db/in-memory"
	"sync"
)

type UserPostgresRepo struct {
	mutex sync.RWMutex
	DB    inmemory.InMemoryDB
}

func New() *UserPostgresRepo {
	repo := UserPostgresRepo{
		DB:    db,
		mutex: sync.RWMutex{},
	}

	_, err := repo.DB.GetTable(UserTableName)
	if errors.Is(err, inmemory.ErrNotExistedTable) {
		repo.DB.CreateTable(UserTableName)
	}

	return &repo
}
