// nolint
package fixtures

import (
	"context"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/dto"

	messageservice "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service/message"
	userservice "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service/user"
)

func LoadFixtures(ctx context.Context, us *userservice.UserService, ms *messageservice.MessageService) {
	users := []dto.UserDTO{
		{
			Username:       "test",
			Email:          "test@mail.ru",
			HashedPassword: "$2a$10$n1ZupQQL9NBnIDHShSIfwut3wf2cUMtsmzBo/7r29oRo4tYRrmoLS",
		},
		{
			Username:       "test2",
			Email:          "test2@mail.ru",
			HashedPassword: "$2a$10$O3bRPhNaWgVibnpkUFL.K.xXwmYnDKKMJ1Ak4iavFrSnn8wAsgYPW",
		},
		{
			Username:       "test3",
			Email:          "test3@mail.ru",
			HashedPassword: "$2a$10$lgQ9a71CwJQkAF1yUcKKl..RGDT4OaGRjyBAVFgGupkdMclmS7wMS",
		},
	}

	for _, user := range users {
		_, err := us.RegisterUser(ctx, user)
		if err != nil {
			return
		}
	}

	pubMessages := []dto.PublicMessageDTO{
		{
			FromID:  1,
			Content: "Hello everyone, I'm Test!",
		},
		{
			FromID:  2,
			Content: "Hi Test, I'm Test2 ;)",
		},
		{
			FromID:  3,
			Content: "What's up! I'm Test3",
		},
	}

	for _, pubMsg := range pubMessages {
		_, err := ms.SendPublicMessage(ctx, pubMsg)
		if err != nil {
			return
		}
	}

	privMessages := []dto.PrivateMessageDTO{
		{
			FromID:  1,
			ToID:    2,
			Content: "Excuse me, where am I?",
		},
		{
			FromID:  2,
			ToID:    1,
			Content: "Ohh.. You are being tested too!",
		},

		{
			FromID:  3,
			ToID:    2,
			Content: "Have something?",
		},
		{
			FromID:  2,
			ToID:    3,
			Content: "What??.. Get off me!",
		},
	}

	for _, privMsg := range privMessages {
		_, err := ms.SendPrivateMessage(ctx, privMsg)
		if err != nil {
			return
		}
	}
}
