package public

import (
	"context"
	"errors"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/mocks"
	testingutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/utils/testing"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"

	repoerrors "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
)

func TestPublicMessageService_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	now := time.Now()

	msgRepoMock := mocks.NewMockPublicMessageRepo(ctrl)
	userRepoMock := mocks.NewMockUserRepo(ctrl)

	service := New(msgRepoMock, userRepoMock)

	type inputArgs = entity.PublicMessage
	type outputArg = *entity.PublicMessage

	tests := []struct {
		name          string
		mockBehaviour func()
		input         inputArgs
		want          outputArg
		wantErr       bool
	}{
		{
			name: "ok, valid username",
			mockBehaviour: func() {
				userRepoMock.
					EXPECT().
					GetUserByUsername(ctx, "username").
					Return(&entity.User{
						ID:             1,
						Email:          "email@mail.com",
						Username:       "username",
						HashedPassword: "hashed_password",
						CreatedAt:      time.Time{},
						UpdatedAt:      time.Time{},
					}, nil)

				msgRepoMock.
					EXPECT().
					AddPublicMessage(ctx, entity.PublicMessage{
						FromUsername: "username",
						Content:      "content",
					}).
					Return(&entity.PublicMessage{
						ID:           1,
						FromUsername: "username",
						Content:      "content",
						SentAt:       now,
						EditedAt:     now,
					}, nil)

			},

			input: inputArgs{
				FromUsername: "username",
				Content:      "content",
			},
			want: &inputArgs{
				ID:           1,
				FromUsername: "username",
				Content:      "content",
				SentAt:       now,
				EditedAt:     now,
			},
		},
		{
			name: "err, invalid username (no such)",
			mockBehaviour: func() {
				userRepoMock.
					EXPECT().
					GetUserByUsername(ctx, "username").
					Return(nil, repoerrors.ErrNoSuchUser)
			},

			input: inputArgs{
				FromUsername: "username",
				Content:      "content",
			},
			wantErr: true,
		},
		{
			name: "err, empty content",
			mockBehaviour: func() {
				userRepoMock.
					EXPECT().
					GetUserByUsername(ctx, "username").
					Return(&entity.User{
						ID:             1,
						Email:          "email@mail.com",
						Username:       "username",
						HashedPassword: "hashed_password",
						CreatedAt:      time.Time{},
						UpdatedAt:      time.Time{},
					}, nil)

				msgRepoMock.
					EXPECT().
					AddPublicMessage(ctx, entity.PublicMessage{
						FromUsername: "username",
						Content:      "",
					}).
					Return(nil, errors.New("empty content not acceptable"))

			},

			input: inputArgs{
				FromUsername: "username",
				Content:      "",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got, err := service.SendPublicMessage(ctx, test.input)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, testingutils.PublicMessagesEquals(*test.want, *got))
			}
		})
	}
}

func TestPublicMessageService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	now := time.Now()

	msgRepoMock := mocks.NewMockPublicMessageRepo(ctrl)
	userRepoMock := mocks.NewMockUserRepo(ctrl)

	service := New(msgRepoMock, userRepoMock)

	type inputArgs = int
	type outputArg = *entity.PublicMessage

	tests := []struct {
		name          string
		mockBehaviour func()
		input         inputArgs
		want          outputArg
		wantErr       bool
	}{
		{
			name: "ok, valid id",
			mockBehaviour: func() {
				msgRepoMock.
					EXPECT().
					GetPublicMessage(ctx, 1).
					Return(&entity.PublicMessage{
						ID:           1,
						FromUsername: "username",
						Content:      "content",
						SentAt:       now,
						EditedAt:     now,
					}, nil)

			},
			input: 1,
			want: &entity.PublicMessage{
				ID:           1,
				FromUsername: "username",
				Content:      "content",
				SentAt:       now,
				EditedAt:     now,
			},
		},
		{
			name: "err, invalid id (no such)",
			mockBehaviour: func() {
				msgRepoMock.
					EXPECT().
					GetPublicMessage(ctx, 1).
					Return(nil, repoerrors.ErrNoSuchUser)

			},
			input:   1,
			wantErr: true,
		},
		{
			name: "err, invalid id (negative value)",
			mockBehaviour: func() {
				msgRepoMock.
					EXPECT().
					GetPublicMessage(ctx, -1).
					Return(nil, repoerrors.ErrNoSuchUser)

			},
			input:   -1,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got, err := service.GetPublicMessage(ctx, test.input)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, testingutils.PublicMessagesEquals(*test.want, *got))
			}
		})
	}
}

func TestPublicMessageService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	now := time.Now()

	msgRepoMock := mocks.NewMockPublicMessageRepo(ctrl)
	userRepoMock := mocks.NewMockUserRepo(ctrl)

	service := New(msgRepoMock, userRepoMock)

	type outputArg = []entity.PublicMessage

	tests := []struct {
		name          string
		mockBehaviour func()
		offset        int
		limit         int
		want          outputArg
		wantErr       bool
	}{
		{
			name: "ok, no offset, no limit",
			mockBehaviour: func() {
				msgRepoMock.
					EXPECT().
					GetAllPublicMessages(ctx, 0, math.MaxInt64).
					Return([]*entity.PublicMessage{
						{
							ID:           1,
							FromUsername: "username",
							Content:      "content",
							SentAt:       now,
							EditedAt:     now,
						},
						{
							ID:           2,
							FromUsername: "username2",
							Content:      "content2",
							SentAt:       now,
							EditedAt:     now,
						},
						{
							ID:           3,
							FromUsername: "username3",
							Content:      "content3",
							SentAt:       now,
							EditedAt:     now,
						},
					})

			},
			offset: 0,
			limit:  math.MaxInt64,
			want: []entity.PublicMessage{
				{
					ID:           1,
					FromUsername: "username",
					Content:      "content",
					SentAt:       now,
					EditedAt:     now,
				},
				{
					ID:           2,
					FromUsername: "username2",
					Content:      "content2",
					SentAt:       now,
					EditedAt:     now,
				},
				{
					ID:           3,
					FromUsername: "username3",
					Content:      "content3",
					SentAt:       now,
					EditedAt:     now,
				},
			},
		},
		{
			name: "ok, offset 1, no limit",
			mockBehaviour: func() {
				msgRepoMock.
					EXPECT().
					GetAllPublicMessages(ctx, 1, math.MaxInt64).
					Return([]*entity.PublicMessage{
						{
							ID:           2,
							FromUsername: "username2",
							Content:      "content2",
							SentAt:       now,
							EditedAt:     now,
						},
						{
							ID:           3,
							FromUsername: "username3",
							Content:      "content3",
							SentAt:       now,
							EditedAt:     now,
						},
					})

			},
			offset: 1,
			limit:  math.MaxInt64,
			want: []entity.PublicMessage{
				{
					ID:           2,
					FromUsername: "username2",
					Content:      "content2",
					SentAt:       now,
					EditedAt:     now,
				},
				{
					ID:           3,
					FromUsername: "username3",
					Content:      "content3",
					SentAt:       now,
					EditedAt:     now,
				},
			},
		},
		{
			name: "ok, offset 1, limit greater than data length",
			mockBehaviour: func() {
				msgRepoMock.
					EXPECT().
					GetAllPublicMessages(ctx, 1, 10).
					Return([]*entity.PublicMessage{
						{
							ID:           2,
							FromUsername: "username2",
							Content:      "content2",
							SentAt:       now,
							EditedAt:     now,
						},
						{
							ID:           3,
							FromUsername: "username3",
							Content:      "content3",
							SentAt:       now,
							EditedAt:     now,
						},
					})

			},
			offset: 1,
			limit:  10,
			want: []entity.PublicMessage{
				{
					ID:           2,
					FromUsername: "username2",
					Content:      "content2",
					SentAt:       now,
					EditedAt:     now,
				},
				{
					ID:           3,
					FromUsername: "username3",
					Content:      "content3",
					SentAt:       now,
					EditedAt:     now,
				},
			},
		},
		{
			name: "ok, offset 1, limit 1",
			mockBehaviour: func() {
				msgRepoMock.
					EXPECT().
					GetAllPublicMessages(ctx, 1, 1).
					Return([]*entity.PublicMessage{
						{
							ID:           2,
							FromUsername: "username2",
							Content:      "content2",
							SentAt:       now,
							EditedAt:     now,
						},
					})

			},
			offset: 1,
			limit:  1,
			want: []entity.PublicMessage{
				{
					ID:           2,
					FromUsername: "username2",
					Content:      "content2",
					SentAt:       now,
					EditedAt:     now,
				},
			},
		},
		{
			name: "ok, offset greater than data length, no limit",
			mockBehaviour: func() {
				msgRepoMock.
					EXPECT().
					GetAllPublicMessages(ctx, 10, math.MaxInt64).
					Return(nil)

			},
			offset: 10,
			limit:  math.MaxInt64,
			want:   nil,
		},
		{
			name: "ok, no offset, limit 0",
			mockBehaviour: func() {
				msgRepoMock.
					EXPECT().
					GetAllPublicMessages(ctx, 0, 0).
					Return(nil)

			},
			offset: 0,
			limit:  0,
			want:   nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got := service.GetAllPublicMessages(ctx, test.offset, test.limit)

			assert.True(t, sliceutils.PointerAndValueSlicesEquals(got, test.want))
		})
	}
}
