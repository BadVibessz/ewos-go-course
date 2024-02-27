// nolint
package repository

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"

	inmemory "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/db/in-memory"
)

type PrivateMessageInMemRepo struct {
	DB      inmemory.InMemoryDB
	counter int
	mutex   sync.RWMutex
}

func NewInMemPrivateMessageRepo(db inmemory.InMemoryDB) *PrivateMessageInMemRepo {
	repo := PrivateMessageInMemRepo{
		DB:      db,
		counter: 1,
		mutex:   sync.RWMutex{},
	}

	_, err := repo.DB.GetTable(PrivateMessageTableName)
	if err != nil && errors.Is(err, inmemory.ErrNotExistedTable) {
		repo.DB.CreateTable(PrivateMessageTableName)
	}

	return &repo
}

func (pr *PrivateMessageInMemRepo) AddPrivateMessage(_ context.Context, msg entity.PrivateMessage) (*entity.PrivateMessage, error) {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	now := time.Now()

	msg.ID = pr.counter
	msg.SentAt = now
	msg.EditedAt = now

	pr.counter++

	err := pr.DB.AddRow(PrivateMessageTableName, strconv.Itoa(msg.ID), msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (pr *PrivateMessageInMemRepo) GetAllPrivateMessages(_ context.Context, offset, limit int) []*entity.PrivateMessage {
	rows, err := pr.DB.GetAllRows(PrivateMessageTableName, offset, limit)
	if err != nil {
		return nil
	}

	res := make([]*entity.PrivateMessage, 0, len(rows))

	for _, row := range rows {
		msg, ok := row.(entity.PrivateMessage)
		if ok {
			res = append(res, &msg)
		}
	}

	sort.Slice(res, func(i, j int) bool { return res[i].SentAt.Before(res[j].SentAt) })

	return res
}

func (pr *PrivateMessageInMemRepo) GetPrivateMessage(_ context.Context, id int) (*entity.PrivateMessage, error) {
	row, err := pr.DB.GetRow(PrivateMessageTableName, strconv.Itoa(id))
	if err != nil {
		return nil, ErrNoSuchPrivateMessage
	}

	msg, ok := row.(entity.PrivateMessage)
	if !ok {
		return nil, ErrNoSuchPrivateMessage
	}

	return &msg, nil
}
