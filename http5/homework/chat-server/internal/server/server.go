package server

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repositories"
)

type Server struct {
	UserRepo repositories.UserRepo
}
