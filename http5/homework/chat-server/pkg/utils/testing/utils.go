package testing

import (
	"database/sql/driver"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"time"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func timesAlmostEquals(tim1, tim2 time.Time) bool {
	return tim1.Sub(tim2) <= 1*time.Second
}

func UsersEquals(usr1, usr2 entity.User) bool {
	return usr1.ID == usr2.ID &&
		usr1.Username == usr2.Username &&
		usr1.Email == usr2.Email &&
		usr1.HashedPassword == usr2.HashedPassword &&
		timesAlmostEquals(usr1.CreatedAt, usr2.CreatedAt) &&
		timesAlmostEquals(usr1.UpdatedAt, usr2.UpdatedAt)
}

func PublicMessagesEquals(msg1, msg2 entity.PublicMessage) bool {
	return msg1.ID == msg2.ID &&
		msg1.FromUsername == msg2.FromUsername &&
		msg1.Content == msg2.Content &&
		timesAlmostEquals(msg1.SentAt, msg2.SentAt) &&
		timesAlmostEquals(msg1.EditedAt, msg2.EditedAt)
}

func PrivateMessagesEquals(msg1, msg2 entity.PrivateMessage) bool {
	return msg1.ID == msg2.ID &&
		msg1.FromUsername == msg2.FromUsername &&
		msg1.ToUsername == msg2.ToUsername &&
		msg1.Content == msg2.Content &&
		timesAlmostEquals(msg1.SentAt, msg2.SentAt) &&
		timesAlmostEquals(msg1.EditedAt, msg2.EditedAt)
}
