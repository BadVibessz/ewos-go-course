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

type PrivateMessageInMemRepo struct {
	DB inmemory.InMemoryDB
}

const privateMessageTableName = "private_messages"

var ErrNoSuchPrivateMessage = errors.New("no such private message")

func NewInMemPrivateMessageRepo(db inmemory.InMemoryDB) *PrivateMessageInMemRepo {
	repo := PrivateMessageInMemRepo{DB: db}

	repo.DB.CreateTable(privateMessageTableName)

	return &repo
}

func (pr *PrivateMessageInMemRepo) AddPrivateMessage(_ context.Context, msg model.PrivateMessage) (*model.PrivateMessage, error) {
	idOffset, err := pr.DB.GetRowsCount(privateMessageTableName)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	toCreate := model.PrivateMessage{
		ID:       idOffset + 1,
		From:     msg.From,
		To:       msg.To,
		Content:  msg.Content,
		SentAt:   now,
		EditedAt: now,
	}

	err = pr.DB.AddRow(privateMessageTableName, strconv.Itoa(toCreate.ID), toCreate)
	if err != nil {
		return nil, err // todo: maybe return custom err of this layer? (ErrNoSuchUser?)
	}

	return &toCreate, nil
}

func (pr *PrivateMessageInMemRepo) GetAllPrivateMessages(_ context.Context, offset, limit int) []*model.PrivateMessage {
	rows, err := pr.DB.GetAllRows(privateMessageTableName, offset, limit)
	if err != nil {
		return nil
	}

	res := make([]*model.PrivateMessage, 0, len(rows))

	for _, row := range rows {
		msg, ok := row.(model.PrivateMessage)
		if ok {
			res = append(res, &msg)
		}
	}

	sort.Slice(res, func(i, j int) bool { return res[i].SentAt.Before(res[j].SentAt) })

	return res
}

func (pr *PrivateMessageInMemRepo) GetPrivateMessage(_ context.Context, id int) (*model.PrivateMessage, error) {
	row, err := pr.DB.GetRow(privateMessageTableName, strconv.Itoa(id))
	if err != nil {
		return nil, ErrNoSuchPrivateMessage
	}

	msg, ok := row.(model.PrivateMessage)
	if !ok {
		return nil, ErrNoSuchPrivateMessage
	}

	return &msg, nil
}

func (pr *PrivateMessageInMemRepo) UpdatePrivateMessage(_ context.Context, id int, newContent string) (*model.PrivateMessage, error) {
	// TODO implement me
	panic("implement me")
}

func (pr *PrivateMessageInMemRepo) DeletePrivateMessage(_ context.Context, id int) (*model.PrivateMessage, error) {
	// TODO implement me
	panic("implement me")
}
