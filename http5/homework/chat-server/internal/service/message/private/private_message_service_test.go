package private

import (
	"context"
	"errors"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/mocks"
	repoerrors "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
	testingutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/testing"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

func TestPrivateMessageService_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	now := time.Now()

	msgRepoMock := mocks.NewMockPrivateMessageRepo(ctrl)
	userRepoMock := mocks.NewMockUserRepo(ctrl)

	service := New(msgRepoMock, userRepoMock)

	type inputArgs = entity.PrivateMessage
	type outputArg = *entity.PrivateMessage

	tests := []struct {
		name          string
		mockBehaviour func()
		input         inputArgs
		want          outputArg
		wantErr       bool
	}{
		{
			name: "ok, valid from username and to username",
			mockBehaviour: func() {
				userRepoMock.
					EXPECT().
					GetUserByUsername(ctx, "from_username").
					Return(&entity.User{
						ID:             1,
						Email:          "from_email@mail.com",
						Username:       "from_username",
						HashedPassword: "hashed_password",
						CreatedAt:      time.Time{},
						UpdatedAt:      time.Time{},
					}, nil)

				userRepoMock.
					EXPECT().
					GetUserByUsername(ctx, "to_username").
					Return(&entity.User{
						ID:             1,
						Email:          "to_email@mail.com",
						Username:       "to_username",
						HashedPassword: "hashed_password",
						CreatedAt:      time.Time{},
						UpdatedAt:      time.Time{},
					}, nil)

				msgRepoMock.
					EXPECT().
					AddPrivateMessage(ctx, entity.PrivateMessage{
						FromUsername: "from_username",
						ToUsername:   "to_username",
						Content:      "content",
					}).
					Return(&entity.PrivateMessage{
						ID:           1,
						FromUsername: "from_username",
						ToUsername:   "to_username",
						Content:      "content",
						SentAt:       now,
						EditedAt:     now,
					}, nil)

			},

			input: inputArgs{
				FromUsername: "from_username",
				ToUsername:   "to_username",
				Content:      "content",
			},
			want: &entity.PrivateMessage{
				ID:           1,
				FromUsername: "from_username",
				ToUsername:   "to_username",
				Content:      "content",
				SentAt:       now,
				EditedAt:     now,
			},
		},
		{
			name: "err, invalid from username (no such)",
			mockBehaviour: func() {
				userRepoMock.
					EXPECT().
					GetUserByUsername(ctx, "from_username").
					Return(nil, repoerrors.ErrNoSuchUser)
			},

			input: inputArgs{
				FromUsername: "from_username",
				ToUsername:   "to_username",
				Content:      "content",
			},
			wantErr: true,
		},
		{
			name: "err, invalid to username (no such)",
			mockBehaviour: func() {
				userRepoMock.
					EXPECT().
					GetUserByUsername(ctx, "from_username").
					Return(&entity.User{
						ID:             1,
						Email:          "from_email@mail.com",
						Username:       "from_username",
						HashedPassword: "hashed_password",
						CreatedAt:      time.Time{},
						UpdatedAt:      time.Time{},
					}, nil)

				userRepoMock.
					EXPECT().
					GetUserByUsername(ctx, "to_username").
					Return(nil, repoerrors.ErrNoSuchUser)
			},

			input: inputArgs{
				FromUsername: "from_username",
				ToUsername:   "to_username",
				Content:      "content",
			},
			wantErr: true,
		},
		{
			name: "err, empty content",
			mockBehaviour: func() {
				userRepoMock.
					EXPECT().
					GetUserByUsername(ctx, "from_username").
					Return(&entity.User{
						ID:             1,
						Email:          "from_email@mail.com",
						Username:       "from_username",
						HashedPassword: "hashed_password",
						CreatedAt:      time.Time{},
						UpdatedAt:      time.Time{},
					}, nil)

				userRepoMock.
					EXPECT().
					GetUserByUsername(ctx, "to_username").
					Return(&entity.User{
						ID:             1,
						Email:          "to_email@mail.com",
						Username:       "to_username",
						HashedPassword: "hashed_password",
						CreatedAt:      time.Time{},
						UpdatedAt:      time.Time{},
					}, nil)

				msgRepoMock.
					EXPECT().
					AddPrivateMessage(ctx, entity.PrivateMessage{
						FromUsername: "from_username",
						ToUsername:   "to_username",
						Content:      "",
					}).
					Return(nil, errors.New("empty content not acceptable"))

			},

			input: inputArgs{
				FromUsername: "from_username",
				ToUsername:   "to_username",
				Content:      "",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got, err := service.SendPrivateMessage(ctx, test.input)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, testingutils.PrivateMessagesEquals(*test.want, *got))
			}
		})
	}
}

func TestPrivateMessageService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	now := time.Now()

	msgRepoMock := mocks.NewMockPrivateMessageRepo(ctrl)
	userRepoMock := mocks.NewMockUserRepo(ctrl)

	service := New(msgRepoMock, userRepoMock)

	type inputArgs = int
	type outputArg = *entity.PrivateMessage

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
					GetPrivateMessage(ctx, 1).
					Return(&entity.PrivateMessage{
						ID:           1,
						FromUsername: "from_username",
						ToUsername:   "to_username",
						Content:      "content",
						SentAt:       now,
						EditedAt:     now,
					}, nil)

			},
			input: 1,
			want: &entity.PrivateMessage{
				ID:           1,
				FromUsername: "from_username",
				ToUsername:   "to_username",
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
					GetPrivateMessage(ctx, 1).
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
					GetPrivateMessage(ctx, -1).
					Return(nil, repoerrors.ErrNoSuchUser)

			},
			input:   -1,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got, err := service.GetPrivateMessage(ctx, test.input)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, testingutils.PrivateMessagesEquals(*test.want, *got))
			}
		})
	}
}

func TestPrivateMessageService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	now := time.Now()

	msgRepoMock := mocks.NewMockPrivateMessageRepo(ctrl)
	userRepoMock := mocks.NewMockUserRepo(ctrl)

	service := New(msgRepoMock, userRepoMock)

	type outputArg = []entity.PrivateMessage

	tests := []struct {
		name          string
		mockBehaviour func()
		toUsername    string
		offset        int
		limit         int
		want          outputArg
		wantErr       bool
	}{
		{
			name: "ok, no offset, no limit", // TODO:
			mockBehaviour: func() {
				msgRepoMock.
					EXPECT().
					GetAllPrivateMessages(ctx, 0, math.MaxInt64).
					Return([]*entity.PrivateMessage{
						{
							ID:           1,
							FromUsername: "from_username1",
							ToUsername:   "to_username",
							Content:      "content",
							SentAt:       now,
							EditedAt:     now,
						},
						{
							ID:           2,
							FromUsername: "from_username",
							ToUsername:   "to_username2",
							Content:      "content2",
							SentAt:       now,
							EditedAt:     now,
						},
						{
							ID:           3,
							FromUsername: "from_username3",
							ToUsername:   "from_username",
							Content:      "content3",
							SentAt:       now,
							EditedAt:     now,
						},
					})

			},
			toUsername: "from_username",
			offset:     0,
			limit:      math.MaxInt64,
			want: []entity.PrivateMessage{
				{
					ID:           2,
					FromUsername: "from_username",
					ToUsername:   "to_username2",
					Content:      "content2",
					SentAt:       now,
					EditedAt:     now,
				},
				{
					ID:           3,
					FromUsername: "from_username3",
					ToUsername:   "from_username",
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
					GetAllPrivateMessages(ctx, 1, math.MaxInt64).
					Return([]*entity.PrivateMessage{
						{
							ID:           2,
							FromUsername: "from_username2",
							ToUsername:   "to_username2",
							Content:      "content2",
							SentAt:       now,
							EditedAt:     now,
						},
						{
							ID:           3,
							FromUsername: "from_username3",
							ToUsername:   "to_username3",
							Content:      "content3",
							SentAt:       now,
							EditedAt:     now,
						},
					})

			},
			offset: 1,
			limit:  math.MaxInt64,
			want: []entity.PrivateMessage{
				{
					ID:           2,
					FromUsername: "from_username2",
					ToUsername:   "to_username2",
					Content:      "content2",
					SentAt:       now,
					EditedAt:     now,
				},
				{
					ID:           3,
					FromUsername: "from_username3",
					ToUsername:   "to_username3",
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
					GetAllPrivateMessages(ctx, 1, 10).
					Return([]*entity.PrivateMessage{
						{
							ID:           2,
							FromUsername: "from_username2",
							ToUsername:   "to_username2",
							Content:      "content2",
							SentAt:       now,
							EditedAt:     now,
						},
						{
							ID:           3,
							FromUsername: "from_username3",
							ToUsername:   "to_username3",
							Content:      "content3",
							SentAt:       now,
							EditedAt:     now,
						},
					})

			},
			offset: 1,
			limit:  10,
			want: []entity.PrivateMessage{
				{
					ID:           2,
					FromUsername: "from_username2",
					ToUsername:   "to_username2",
					Content:      "content2",
					SentAt:       now,
					EditedAt:     now,
				},
				{
					ID:           3,
					FromUsername: "from_username3",
					ToUsername:   "to_username3",
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
					GetAllPrivateMessages(ctx, 1, 1).
					Return([]*entity.PrivateMessage{
						{
							ID:           2,
							FromUsername: "from_username2",
							ToUsername:   "to_username2",
							Content:      "content2",
							SentAt:       now,
							EditedAt:     now,
						}})

			},
			offset: 1,
			limit:  1,
			want: []entity.PrivateMessage{
				{
					ID:           2,
					FromUsername: "from_username2",
					ToUsername:   "to_username2",
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
					GetAllPrivateMessages(ctx, 10, math.MaxInt64).
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
					GetAllPrivateMessages(ctx, 0, 0).
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

			got := service.GetAllPrivateMessages(ctx, test.toUsername, test.offset, test.limit)

			assert.True(t, sliceutils.PointerAndValueSlicesEquals(got, test.want))
		})
	}
}
