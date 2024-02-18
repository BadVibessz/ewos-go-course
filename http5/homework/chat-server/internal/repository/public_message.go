// nolint
package repository

import (
	"context"
	"sort"
	"strconv"
	"time"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"

	inmemory "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/db/in-memory"
)

type PublicMessageInMemRepo struct {
	DB inmemory.InMemoryDB
}

func NewInMemPublicMessageRepo(db inmemory.InMemoryDB) *PublicMessageInMemRepo {
	repo := PublicMessageInMemRepo{DB: db}

	_, err := repo.DB.GetTable(PublicMessageTableName)
	if err != nil {
		repo.DB.CreateTable(PublicMessageTableName)
	}

	return &repo
}

func (pr *PublicMessageInMemRepo) AddPublicMessage(_ context.Context, msg model.PublicMessage) (*model.PublicMessage, error) {
	idOffset, err := pr.DB.GetRowsCount(PublicMessageTableName)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	msg.ID = idOffset + 1
	msg.SentAt = now
	msg.EditedAt = now

	err = pr.DB.AddRow(PublicMessageTableName, strconv.Itoa(msg.ID), msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (pr *PublicMessageInMemRepo) GetAllPublicMessages(_ context.Context, offset, limit int) []*model.PublicMessage {
	rows, err := pr.DB.GetAllRows(PublicMessageTableName, offset, limit)
	if err != nil {
		return nil
	}

	res := make([]*model.PublicMessage, 0, len(rows))

	for _, row := range rows {
		msg, ok := row.(model.PublicMessage)
		if ok {
			res = append(res, &msg)
		}
	}

	sort.Slice(res, func(i, j int) bool { return res[i].SentAt.Before(res[j].SentAt) })

	return res
}

func (pr *PublicMessageInMemRepo) GetPublicMessage(_ context.Context, id int) (*model.PublicMessage, error) {
	row, err := pr.DB.GetRow(PublicMessageTableName, strconv.Itoa(id))
	if err != nil {
		return nil, ErrNoSuchPublicMessage
	}

	msg, ok := row.(model.PublicMessage)
	if !ok {
		return nil, ErrNoSuchPublicMessage
	}

	return &msg, nil
}
