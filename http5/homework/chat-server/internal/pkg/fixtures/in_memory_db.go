// nolint
package fixtures

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository/in-memory"
	"strconv"
	"time"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	inmemory "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/db/in-memory"
)

func LoadFixtures(db inmemory.InMemoryDB) {
	now := time.Now()

	users := []entity.User{
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

	db.CreateTable(in_memory.UserTableName)

	for _, user := range users {
		err := db.AddRow(in_memory.UserTableName, strconv.Itoa(user.ID), user)
		if err != nil {
			return
		}
	}

	pubMessages := []entity.PublicMessage{
		{
			ID:           1,
			FromUsername: "test",
			Content:      "Hello everyone, I'm Test!",
			SentAt:       now,
			EditedAt:     now,
		},
		{
			ID:           2,
			FromUsername: "test2",
			Content:      "Hello everyone, I'm Test2 ;)",
			SentAt:       now,
			EditedAt:     now,
		},
		{
			ID:           3,
			FromUsername: "test3",
			Content:      "What's up! I'm Test3",
			SentAt:       now,
			EditedAt:     now,
		},
	}

	db.CreateTable(in_memory.PublicMessageTableName)

	for _, pubMsg := range pubMessages {
		err := db.AddRow(in_memory.PublicMessageTableName, strconv.Itoa(pubMsg.ID), pubMsg)
		if err != nil {
			return
		}
	}

	privMessages := []entity.PrivateMessage{
		{
			ID:           1,
			FromUsername: "test",
			ToUsername:   "test2",
			Content:      "Excuse me, where am I?",
			SentAt:       now,
			EditedAt:     now,
		},
		{
			ID:           2,
			FromUsername: "test2",
			ToUsername:   "test",
			Content:      "Ohh.. You are being tested too!",
			SentAt:       now,
			EditedAt:     now,
		},
		{
			ID:           3,
			FromUsername: "test3",
			ToUsername:   "test2",
			Content:      "Have something?",
			SentAt:       now,
			EditedAt:     now,
		},
		{
			ID:           4,
			FromUsername: "test2",
			ToUsername:   "test3",
			Content:      "What??.. Get off me!",
			SentAt:       now,
			EditedAt:     now,
		},
	}

	db.CreateTable(in_memory.PrivateMessageTableName)

	for _, privMsg := range privMessages {
		err := db.AddRow(in_memory.PrivateMessageTableName, strconv.Itoa(privMsg.ID), privMsg)
		if err != nil {
			return
		}
	}
}
