package repository

import (
	"context"
	"errors"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"
	inmemory "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/db/in-memory"
	"sort"
	"strconv"
	"time"
)

type PublicMessageInMemRepo struct {
	DB *inmemory.InMemDB
}

const publicMessageTableName = "public_messages"

var (
	NoSuchPublicMessageErr = errors.New("no such public message")
)

func NewInMemPublicMessageRepo(db *inmemory.InMemDB) *PublicMessageInMemRepo {
	repo := PublicMessageInMemRepo{DB: db}

	repo.DB.CreateTable(publicMessageTableName)

	return &repo
}

func (pr *PublicMessageInMemRepo) AddPublicMessage(ctx context.Context, msg model.PublicMessage) (*model.PublicMessage, error) {
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
		return nil, err // todo: maybe return custom err of this layer? (NoSuchUserErr?)
	}

	return &toCreate, nil
}

func (pr *PublicMessageInMemRepo) GetAllPublicMessages(ctx context.Context) []*model.PublicMessage {
	rows, err := pr.DB.GetAllRows(publicMessageTableName)
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

func (pr *PublicMessageInMemRepo) GetPublicMessage(ctx context.Context, id int) (*model.PublicMessage, error) {
	row, err := pr.DB.GetRow(publicMessageTableName, strconv.Itoa(id))
	if err != nil {
		return nil, NoSuchPublicMessageErr
	}

	msg, ok := row.(model.PublicMessage)
	if !ok {
		return nil, NoSuchPublicMessageErr
	}

	return &msg, nil
}

func (pr *PublicMessageInMemRepo) UpdatePublicMessage(ctx context.Context, id int, newContent string) (*model.PublicMessage, error) {
	//TODO implement me
	panic("implement me")
}

func (pr *PublicMessageInMemRepo) DeletePublicMessage(ctx context.Context, id int) (*model.PublicMessage, error) {
	//TODO implement me
	panic("implement me")
}
