// nolint
package fixtures

import (
	"strconv"
	"time"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"

	inmemory "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/db/in-memory"
)

func LoadFixtures(db inmemory.InMemoryDB) {
	now := time.Now()

	users := []model.User{
		{
			ID:             1,
			Username:       "test",
			Email:          "test@mail.ru",
			HashedPassword: "$2a$10$n1ZupQQL9NBnIDHShSIfwut3wf2cUMtsmzBo/7r29oRo4tYRrmoLS",
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             2,
			Username:       "test2",
			Email:          "test2@mail.ru",
			HashedPassword: "$2a$10$O3bRPhNaWgVibnpkUFL.K.xXwmYnDKKMJ1Ak4iavFrSnn8wAsgYPW",
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             3,
			Username:       "test3",
			Email:          "test3@mail.ru",
			HashedPassword: "$2a$10$lgQ9a71CwJQkAF1yUcKKl..RGDT4OaGRjyBAVFgGupkdMclmS7wMS",
			CreatedAt:      now,
			UpdatedAt:      now,
		},
	}

	for _, user := range users {
		err := db.AddRow(repository.UserTableName, strconv.Itoa(user.ID), user)
		if err != nil {
			return
		}
	}

	pubMessages := []model.PublicMessage{
		{
			ID:       1,
			From:     &users[0],
			Content:  "Hello everyone, I'm Test!",
			SentAt:   now,
			EditedAt: now,
		},
		{
			ID:       2,
			From:     &users[1],
			Content:  "Hello everyone, I'm Test2 ;)",
			SentAt:   now,
			EditedAt: now,
		},
		{
			ID:       3,
			From:     &users[2],
			Content:  "What's up! I'm Test3",
			SentAt:   now,
			EditedAt: now,
		},
	}

	for _, pubMsg := range pubMessages {
		err := db.AddRow(repository.PublicMessageTableName, strconv.Itoa(pubMsg.ID), pubMsg)
		if err != nil {
			return
		}
	}

	privMessages := []model.PrivateMessage{
		{
			ID:       1,
			From:     &users[0],
			To:       &users[1],
			Content:  "Excuse me, where am I?",
			SentAt:   now,
			EditedAt: now,
		},
		{
			ID:       2,
			From:     &users[1],
			To:       &users[0],
			Content:  "Ohh.. You are being tested too!",
			SentAt:   now,
			EditedAt: now,
		},
		{
			ID:       3,
			From:     &users[2],
			To:       &users[1],
			Content:  "Have something?",
			SentAt:   now,
			EditedAt: now,
		},
		{
			ID:       4,
			From:     &users[1],
			To:       &users[2],
			Content:  "What??.. Get off me!",
			SentAt:   now,
			EditedAt: now,
		},
	}

	for _, privMsg := range privMessages {
		err := db.AddRow(repository.PrivateMessageTableName, strconv.Itoa(privMsg.ID), privMsg)
		if err != nil {
			return
		}
	}
}
