// nolint
package mapper

import (
	"errors"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/request"
)

var ErrCannotRetrieveUsernameAndPass = errors.New("cannot retrieve username and password from basic auth header")

func MapBasicAuthToLoginRequest(username, pass string, ok bool) (*request.LoginRequest, error) {
	if !ok {
		return nil, ErrCannotRetrieveUsernameAndPass
	}

	loginReq := request.LoginRequest{
		Username: username,
		Password: pass,
	}

	return &loginReq, nil
}
