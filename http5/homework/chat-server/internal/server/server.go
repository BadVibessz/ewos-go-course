package server

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
)

type Server struct {
	UserRepo repository.UserRepo
}
