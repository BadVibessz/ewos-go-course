package postgres

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
)

type PrivateMessageRepo struct {
	DB *sqlx.DB
}

func NewPrivateMessageRepo(db *sqlx.DB) *PrivateMessageRepo {
	return &PrivateMessageRepo{
		DB: db,
	}
}

func (pr *PrivateMessageRepo) AddPrivateMessage(ctx context.Context, msg entity.PrivateMessage) (*entity.PrivateMessage, error) {
	now := time.Now()

	msg.SentAt = now
	msg.EditedAt = now

	result, err := pr.DB.NamedQueryContext(ctx,
		`INSERT INTO private_message (from_username, to_username, content, sent_at, edited_at) 
VALUES (:from_username, :to_username, :content, :sent_at, :edited_at) 
RETURNING id, from_username, to_username, content, sent_at, edited_at`,
		&msg)
	if err != nil {
		return nil, err
	}

	var resMsg entity.PrivateMessage

	if result.Next() {
		if err = result.StructScan(&resMsg); err != nil {
			return nil, err
		}
	}

	return &resMsg, nil
}

func (pr *PrivateMessageRepo) GetAllPrivateMessages(ctx context.Context, offset, limit int) []*entity.PrivateMessage {
	var query string

	if limit == math.MaxInt64 {
		query = fmt.Sprintf("SELECT * FROM private_message ORDER BY sent_at OFFSET %v", offset)
	} else {
		query = fmt.Sprintf("SELECT * FROM private_message ORDER BY sent_at LIMIT %v OFFSET %v", limit, offset)
	}

	rows, err := pr.DB.QueryxContext(ctx, query)
	if err != nil {
		return nil
	}

	var users []*entity.PrivateMessage

	for rows.Next() {
		var msg entity.PrivateMessage

		err = rows.StructScan(&msg)
		if err != nil {
			return nil
		}

		users = append(users, &msg)
	}

	return users
}

func (pr *PrivateMessageRepo) GetPrivateMessage(ctx context.Context, id int) (*entity.PrivateMessage, error) {
	row := pr.DB.QueryRowxContext(ctx, "SELECT * FROM private_message WHERE id = $1", id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var msg entity.PrivateMessage

	err := row.StructScan(&msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}
