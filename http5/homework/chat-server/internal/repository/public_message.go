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

type PublicMessageInMemRepo struct {
	DB    inmemory.InMemoryDB
	mutex sync.RWMutex
}

func NewInMemPublicMessageRepo(db inmemory.InMemoryDB) *PublicMessageInMemRepo {
	repo := PublicMessageInMemRepo{
		DB:    db,
		mutex: sync.RWMutex{},
	}

	_, err := repo.DB.GetTable(PublicMessageTableName)
	if err != nil && errors.Is(err, inmemory.ErrNotExistedTable) {
		repo.DB.CreateTable(PublicMessageTableName)
	}

	return &repo
}

func (pr *PublicMessageInMemRepo) AddPublicMessage(_ context.Context, msg entity.PublicMessage) (*entity.PublicMessage, error) {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	idOffset, err := pr.DB.GetTableCounter(PublicMessageTableName)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	msg.ID = idOffset + 1
	msg.SentAt = now
	msg.EditedAt = now

	if err = pr.DB.AddRow(PublicMessageTableName, strconv.Itoa(msg.ID), msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

func (pr *PublicMessageInMemRepo) GetAllPublicMessages(_ context.Context, offset, limit int) []*entity.PublicMessage {
	rows, err := pr.DB.GetAllRows(PublicMessageTableName, offset, limit)
	if err != nil {
		return nil
	}

	res := make([]*entity.PublicMessage, 0, len(rows))

	for _, row := range rows {
		msg, ok := row.(entity.PublicMessage)
		if ok {
			res = append(res, &msg)
		}
	}

	sort.Slice(res, func(i, j int) bool { return res[i].SentAt.Before(res[j].SentAt) })

	return res
}

func (pr *PublicMessageInMemRepo) GetPublicMessage(_ context.Context, id int) (*entity.PublicMessage, error) {
	row, err := pr.DB.GetRow(PublicMessageTableName, strconv.Itoa(id))
	if err != nil {
		return nil, ErrNoSuchPublicMessage
	}

	msg, ok := row.(entity.PublicMessage)
	if !ok {
		return nil, ErrNoSuchPublicMessage
	}

	return &msg, nil
}
