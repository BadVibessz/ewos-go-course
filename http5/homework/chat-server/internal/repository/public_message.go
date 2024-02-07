package repository

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"time"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"

	inmemory "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/db/in-memory"
)

type PublicMessageInMemRepo struct {
	DB inmemory.InMemoryDB
}

const publicMessageTableName = "public_messages"

var ErrNoSuchPublicMessage = errors.New("no such public message")

func NewInMemPublicMessageRepo(db inmemory.InMemoryDB) *PublicMessageInMemRepo {
	repo := PublicMessageInMemRepo{DB: db}

	repo.DB.CreateTable(publicMessageTableName)

	return &repo
}

func (pr *PublicMessageInMemRepo) AddPublicMessage(_ context.Context, msg model.PublicMessage) (*model.PublicMessage, error) {
	idOffset, err := pr.DB.GetRowsCount(publicMessageTableName)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	toCreate := model.PublicMessage{
		ID:       idOffset + 1,
		From:     msg.From,
		Content:  msg.Content,
		SentAt:   now,
		EditedAt: now,
	}

	err = pr.DB.AddRow(publicMessageTableName, strconv.Itoa(toCreate.ID), toCreate)
	if err != nil {
		return nil, err // todo: maybe return custom err of this layer? (ErrNoSuchUser?)
	}

	return &toCreate, nil
}

func (pr *PublicMessageInMemRepo) GetAllPublicMessages(_ context.Context, offset, limit int) []*model.PublicMessage {
	rows, err := pr.DB.GetAllRows(publicMessageTableName, offset, limit)
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
	row, err := pr.DB.GetRow(publicMessageTableName, strconv.Itoa(id))
	if err != nil {
		return nil, ErrNoSuchPublicMessage
	}

	msg, ok := row.(model.PublicMessage)
	if !ok {
		return nil, ErrNoSuchPublicMessage
	}

	return &msg, nil
}
