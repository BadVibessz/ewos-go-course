package private

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/mocks"
	repoerrors "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
	testingutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/testing"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var now = time.Now()

func TestPrivateMessageService_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

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

	msgRepoMock := mocks.NewMockPrivateMessageRepo(ctrl)
	userRepoMock := mocks.NewMockUserRepo(ctrl)

	service := New(msgRepoMock, userRepoMock)

	messages := []*entity.PrivateMessage{
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
		{
			ID:           4,
			FromUsername: "from_username4",
			ToUsername:   "from_username4",
			Content:      "content4",
			SentAt:       now,
			EditedAt:     now,
		},
		{
			ID:           5,
			FromUsername: "from_username5",
			ToUsername:   "from_username",
			Content:      "content5",
			SentAt:       now,
			EditedAt:     now,
		},
		{
			ID:           6,
			FromUsername: "from_username6",
			ToUsername:   "from_username6",
			Content:      "content6",
			SentAt:       now,
			EditedAt:     now,
		},
		{
			ID:           7,
			FromUsername: "from_username7",
			ToUsername:   "from_username7",
			Content:      "content7",
			SentAt:       now,
			EditedAt:     now,
		},
	}

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
			name: "ok, no offset, no limit", // TODO: SOMETHING WRONG WITH SLICEUTILS.FILTER FUNC???
			mockBehaviour: func() {
				msgRepoMock.
					EXPECT().
					GetAllPrivateMessages(ctx, 0, math.MaxInt64).
					Return(messages)

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
				{
					ID:           5,
					FromUsername: "from_username5",
					ToUsername:   "from_username",
					Content:      "content5",
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
					GetAllPrivateMessages(ctx, 0, math.MaxInt64).
					Return(messages)

			},
			toUsername: "from_username",
			offset:     1,
			limit:      math.MaxInt64,
			want: []entity.PrivateMessage{
				{
					ID:           3,
					FromUsername: "from_username3",
					ToUsername:   "from_username",
					Content:      "content3",
					SentAt:       now,
					EditedAt:     now,
				},
				{
					ID:           5,
					FromUsername: "from_username5",
					ToUsername:   "from_username",
					Content:      "content5",
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
					GetAllPrivateMessages(ctx, 0, math.MaxInt64).
					Return(messages)

			},
			toUsername: "from_username",
			offset:     1,
			limit:      10,
			want: []entity.PrivateMessage{
				{
					ID:           3,
					FromUsername: "from_username3",
					ToUsername:   "from_username",
					Content:      "content3",
					SentAt:       now,
					EditedAt:     now,
				},
				{
					ID:           5,
					FromUsername: "from_username5",
					ToUsername:   "from_username",
					Content:      "content5",
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
					GetAllPrivateMessages(ctx, 0, math.MaxInt64).
					Return(messages)

			},
			toUsername: "from_username",
			offset:     1,
			limit:      1,
			want: []entity.PrivateMessage{
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
			name: "ok, offset greater than data length, no limit",
			mockBehaviour: func() {
				msgRepoMock.
					EXPECT().
					GetAllPrivateMessages(ctx, 0, math.MaxInt64).
					Return(messages)

			},
			toUsername: "from_username",
			offset:     10,
			limit:      math.MaxInt64,
			want:       nil,
		},
		{
			name: "ok, no offset, limit 0",
			mockBehaviour: func() {
				msgRepoMock.
					EXPECT().
					GetAllPrivateMessages(ctx, 0, math.MaxInt64).
					Return(messages)

			},
			toUsername: "from_username",
			offset:     0,
			limit:      0,
			want:       nil,
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

func TestPrivateMessageService_GetAllFromUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	msgRepoMock := mocks.NewMockPrivateMessageRepo(ctrl)
	userRepoMock := mocks.NewMockUserRepo(ctrl)

	service := New(msgRepoMock, userRepoMock)

	messages := []*entity.PrivateMessage{
		{
			ID:           1,
			FromUsername: "from_username",
			ToUsername:   "to_username",
			Content:      "content",
			SentAt:       now,
			EditedAt:     now,
		},
		{
			ID:           2,
			FromUsername: "to_username",
			ToUsername:   "to_username2",
			Content:      "content2",
			SentAt:       now,
			EditedAt:     now,
		},
		{
			ID:           3,
			FromUsername: "username3",
			ToUsername:   "to_username",
			Content:      "content3",
			SentAt:       now,
			EditedAt:     now,
		},
		{
			ID:           4,
			FromUsername: "from_username",
			ToUsername:   "to_username",
			Content:      "content4",
			SentAt:       now,
			EditedAt:     now,
		},
		{
			ID:           5,
			FromUsername: "from_username",
			ToUsername:   "to_username",
			Content:      "content5",
			SentAt:       now,
			EditedAt:     now,
		},
		{
			ID:           6,
			FromUsername: "from_username6",
			ToUsername:   "from_username6",
			Content:      "content6",
			SentAt:       now,
			EditedAt:     now,
		},
		{
			ID:           7,
			FromUsername: "from_username7",
			ToUsername:   "from_username7",
			Content:      "content7",
			SentAt:       now,
			EditedAt:     now,
		},
	}

	type outputArg = []entity.PrivateMessage

	tests := []struct {
		name          string
		mockBehaviour func()
		toUsername    string
		fromUsername  string
		offset        int
		limit         int
		want          outputArg
		wantErr       bool
	}{
		{
			name: "ok, valid to username, valid form username, no offset, no limit",
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
					GetAllPrivateMessages(ctx, 0, math.MaxInt64).
					Return(messages)

			},
			fromUsername: "from_username",
			toUsername:   "to_username",
			offset:       0,
			limit:        math.MaxInt64,
			want: []entity.PrivateMessage{
				{
					ID:           1,
					FromUsername: "from_username",
					ToUsername:   "to_username",
					Content:      "content",
					SentAt:       now,
					EditedAt:     now,
				},
				{
					ID:           4,
					FromUsername: "from_username",
					ToUsername:   "to_username",
					Content:      "content4",
					SentAt:       now,
					EditedAt:     now,
				},
				{
					ID:           5,
					FromUsername: "from_username",
					ToUsername:   "to_username",
					Content:      "content5",
					SentAt:       now,
					EditedAt:     now,
				},
			},
		},
		{
			name: "err, invalid to username, valid form username, no offset, no limit",
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
					GetUserByUsername(ctx, "invalid_to_username").
					Return(nil, repoerrors.ErrNoSuchUser)
			},
			fromUsername: "from_username",
			toUsername:   "invalid_to_username",
			offset:       0,
			limit:        math.MaxInt64,
			wantErr:      true,
		},
		{
			name: "err, valid to username, invalid form username, no offset, no limit",
			mockBehaviour: func() {
				userRepoMock.
					EXPECT().
					GetUserByUsername(ctx, "invalid_from_username").
					Return(nil, repoerrors.ErrNoSuchUser)
			},
			fromUsername: "invalid_from_username",
			toUsername:   "to_username",
			offset:       0,
			limit:        math.MaxInt64,
			wantErr:      true,
		},
		{
			name: "ok, valid to username, valid form username, offset 1, no limit",
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
					GetAllPrivateMessages(ctx, 0, math.MaxInt64).
					Return(messages)

			},
			fromUsername: "from_username",
			toUsername:   "to_username",
			offset:       1,
			limit:        math.MaxInt64,
			want: []entity.PrivateMessage{
				{
					ID:           4,
					FromUsername: "from_username",
					ToUsername:   "to_username",
					Content:      "content4",
					SentAt:       now,
					EditedAt:     now,
				},
				{
					ID:           5,
					FromUsername: "from_username",
					ToUsername:   "to_username",
					Content:      "content5",
					SentAt:       now,
					EditedAt:     now,
				},
			},
		},
		{
			name: "ok, valid to username, valid form username, offset 1, limit 1",
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
					GetAllPrivateMessages(ctx, 0, math.MaxInt64).
					Return(messages)

			},
			fromUsername: "from_username",
			toUsername:   "to_username",
			offset:       1,
			limit:        1,
			want: []entity.PrivateMessage{
				{
					ID:           4,
					FromUsername: "from_username",
					ToUsername:   "to_username",
					Content:      "content4",
					SentAt:       now,
					EditedAt:     now,
				},
			},
		},
		{
			name: "nil, valid to username, valid form username, offset greater than data length, no limit",
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
					GetAllPrivateMessages(ctx, 0, math.MaxInt64).
					Return(messages)

			},
			fromUsername: "from_username",
			toUsername:   "to_username",
			offset:       100,
			limit:        math.MaxInt64,
			want:         nil,
		},
		{
			name: "nil, valid to username, valid form username, no offset, limit 0",
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
					GetAllPrivateMessages(ctx, 0, math.MaxInt64).
					Return(messages)

			},
			fromUsername: "from_username",
			toUsername:   "to_username",
			offset:       0,
			limit:        0,
			want:         nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got, err := service.GetAllPrivateMessagesFromUser(ctx, test.toUsername, test.fromUsername, test.offset, test.limit)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, sliceutils.PointerAndValueSlicesEquals(got, test.want))
			}

		})
	}
}

func TestPrivateMessageService_GetAllUsersThatSentMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	msgRepoMock := mocks.NewMockPrivateMessageRepo(ctrl)
	userRepoMock := mocks.NewMockUserRepo(ctrl)

	service := New(msgRepoMock, userRepoMock)

	messages := []*entity.PrivateMessage{
		{
			ID:           1,
			FromUsername: "username1",
			ToUsername:   "to_username",
			Content:      "content",
			SentAt:       now,
			EditedAt:     now,
		},
		{
			ID:           2,
			FromUsername: "to_username",
			ToUsername:   "to_username2",
			Content:      "content2",
			SentAt:       now,
			EditedAt:     now,
		},
		{
			ID:           3,
			FromUsername: "username3",
			ToUsername:   "to_username",
			Content:      "content3",
			SentAt:       now,
			EditedAt:     now,
		},
		{
			ID:           4,
			FromUsername: "to_username",
			ToUsername:   "to_username4",
			Content:      "content4",
			SentAt:       now,
			EditedAt:     now,
		},
		{
			ID:           5,
			FromUsername: "username5",
			ToUsername:   "to_username",
			Content:      "content5",
			SentAt:       now,
			EditedAt:     now,
		},
		{
			ID:           6,
			FromUsername: "from_username6",
			ToUsername:   "from_username6",
			Content:      "content6",
			SentAt:       now,
			EditedAt:     now,
		},
		{
			ID:           7,
			FromUsername: "from_username7",
			ToUsername:   "from_username7",
			Content:      "content7",
			SentAt:       now,
			EditedAt:     now,
		},
	}
	users := []*entity.User{
		{
			ID:             1,
			Email:          "email1@mail.com",
			Username:       "username1",
			HashedPassword: "hashed_password1",
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             2,
			Email:          "email2@mail.com",
			Username:       "username2",
			HashedPassword: "hashed_password2",
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             3,
			Email:          "email3@mail.com",
			Username:       "username3",
			HashedPassword: "hashed_password3",
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             4,
			Email:          "email4@mail.com",
			Username:       "username4",
			HashedPassword: "hashed_password4",
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             5,
			Email:          "email5@mail.com",
			Username:       "username5",
			HashedPassword: "hashed_password5",
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             6,
			Email:          "email6@mail.com",
			Username:       "username6",
			HashedPassword: "hashed_password6",
			CreatedAt:      now,
			UpdatedAt:      now,
		},
	}

	type outputArg = []entity.User

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
			name: "ok, valid to username, no offset, no limit",
			mockBehaviour: func() {
				msgRepoMock.
					EXPECT().
					GetAllPrivateMessages(ctx, 0, math.MaxInt64).
					Return(messages)

				userRepoMock.
					EXPECT().
					GetAllUsers(ctx, 0, math.MaxInt64).
					Return(users)

			},
			toUsername: "to_username",
			offset:     0,
			limit:      math.MaxInt64,
			want: []entity.User{
				{
					ID:             1,
					Email:          "email1@mail.com",
					Username:       "username1",
					HashedPassword: "hashed_password1",
					CreatedAt:      now,
					UpdatedAt:      now,
				},
				{
					ID:             3,
					Email:          "email3@mail.com",
					Username:       "username3",
					HashedPassword: "hashed_password3",
					CreatedAt:      now,
					UpdatedAt:      now,
				},
				{
					ID:             5,
					Email:          "email5@mail.com",
					Username:       "username5",
					HashedPassword: "hashed_password5",
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
		},
		{
			name: "nil, invalid to username, no offset, no limit",
			mockBehaviour: func() {
				msgRepoMock.
					EXPECT().
					GetAllPrivateMessages(ctx, 0, math.MaxInt64).
					Return(messages)
			},
			toUsername: "username",
			offset:     0,
			limit:      math.MaxInt64,
			want:       nil,
		},
		{
			name: "ok, valid to username, offset 1, no limit",
			mockBehaviour: func() {
				msgRepoMock.
					EXPECT().
					GetAllPrivateMessages(ctx, 0, math.MaxInt64).
					Return(messages)

				userRepoMock.
					EXPECT().
					GetAllUsers(ctx, 0, math.MaxInt64).
					Return(users)

			},
			toUsername: "to_username",
			offset:     1,
			limit:      math.MaxInt64,
			want: []entity.User{
				{
					ID:             3,
					Email:          "email3@mail.com",
					Username:       "username3",
					HashedPassword: "hashed_password3",
					CreatedAt:      now,
					UpdatedAt:      now,
				},
				{
					ID:             5,
					Email:          "email5@mail.com",
					Username:       "username5",
					HashedPassword: "hashed_password5",
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
		},
		{
			name: "ok, valid to username, offset 1, limit 1",
			mockBehaviour: func() {
				msgRepoMock.
					EXPECT().
					GetAllPrivateMessages(ctx, 0, math.MaxInt64).
					Return(messages)

				userRepoMock.
					EXPECT().
					GetAllUsers(ctx, 0, math.MaxInt64).
					Return(users)

			},
			toUsername: "to_username",
			offset:     1,
			limit:      1,
			want: []entity.User{
				{
					ID:             3,
					Email:          "email3@mail.com",
					Username:       "username3",
					HashedPassword: "hashed_password3",
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
		},
		{
			name: "nil, valid to username, offset greater than data length, no limit",
			mockBehaviour: func() {
				msgRepoMock.
					EXPECT().
					GetAllPrivateMessages(ctx, 0, math.MaxInt64).
					Return(messages)
			},
			toUsername: "to_username",
			offset:     100,
			limit:      math.MaxInt64,
			want:       nil,
		},
		{
			name:          "nil, valid to username, no offset, limit 0",
			mockBehaviour: func() {},
			toUsername:    "to_username",
			offset:        0,
			limit:         0,
			want:          nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour()

			got := service.GetAllUsersThatSentMessage(ctx, test.toUsername, test.offset, test.limit)

			assert.True(t, sliceutils.PointerAndValueSlicesEquals(got, test.want))
		})
	}
}
