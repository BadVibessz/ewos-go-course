package repository

import (
	"context"
	"errors"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"
	inmemory "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/db/in-memory"
	"strconv"
	"time"
)

type PrivateMessageInMemRepo struct {
	DB *inmemory.InMemDB
}

const privateMessageTableName = "private_messages"

var (
	NoSuchPrivateMessageErr = errors.New("no such private message")
)

func NewInMemPrivateMessageRepo(db *inmemory.InMemDB) *PrivateMessageInMemRepo {
	repo := PrivateMessageInMemRepo{DB: db}

	repo.DB.CreateTable(privateMessageTableName)

	return &repo
}

func (pr *PrivateMessageInMemRepo) AddPrivateMessage(ctx context.Context, msg model.PrivateMessage) (*model.PrivateMessage, error) {
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
		return nil, err // todo: maybe return custom err of this layer? (NoSuchUserErr?)
	}

	return &toCreate, nil
}

func (pr *PrivateMessageInMemRepo) GetAllPrivateMessages(ctx context.Context) []*model.PrivateMessage {
	rows, err := pr.DB.GetAllRows(privateMessageTableName)
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

	return res
}

func (pr *PrivateMessageInMemRepo) GetPrivateMessage(ctx context.Context, id int) (*model.PrivateMessage, error) {
	row, err := pr.DB.GetRow(privateMessageTableName, strconv.Itoa(id))
	if err != nil {
		return nil, NoSuchPrivateMessageErr
	}

	msg, ok := row.(model.PrivateMessage)
	if !ok {
		return nil, NoSuchPrivateMessageErr
	}

	return &msg, nil
}

func (pr *PrivateMessageInMemRepo) UpdatePrivateMessage(ctx context.Context, id int, newContent string) (*model.PrivateMessage, error) {
	//TODO implement me
	panic("implement me")
}

func (pr *PrivateMessageInMemRepo) DeletePrivateMessage(ctx context.Context, id int) (*model.PrivateMessage, error) {
	//TODO implement me
	panic("implement me")
}
