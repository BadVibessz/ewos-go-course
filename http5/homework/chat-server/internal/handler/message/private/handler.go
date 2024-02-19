// nolint
package private

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/mapper"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/middleware"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/request"

	messageservice "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service/message"

	handlerinternalutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/utils/handler"
	handlerutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/handler"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
)

type PrivateMessageService interface {
	SendPrivateMessage(ctx context.Context, fromID, toID int, content string) (*entity.PrivateMessage, error)
	GetPrivateMessage(ctx context.Context, id int) (*entity.PrivateMessage, error)
	GetAllPrivateMessages(ctx context.Context, userToID int, offset, limit int) []*entity.PrivateMessage
	GetAllPrivateMessagesFromUser(ctx context.Context, toID, fromID int, offset, limit int) ([]*entity.PrivateMessage, error)
}

type UserService interface {
	RegisterUser(ctx context.Context, user entity.User) (*entity.User, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	GetAllUsers(ctx context.Context, offset, limit int) []*entity.User
	UpdateUser(ctx context.Context, id int, updateModel entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, id int) (*entity.User, error)
}

type AuthService interface {
	Login(ctx context.Context, loginReq request.LoginRequest) (*entity.User, error)
}

type Handler struct {
	MessageService PrivateMessageService
	UserService    UserService
	AuthService    AuthService
	logger         *logrus.Logger
	valid          *validator.Validate
}

func New(privateMessageService PrivateMessageService, userService UserService, authService AuthService,
	logger *logrus.Logger, valid *validator.Validate) *Handler {
	return &Handler{
		MessageService: privateMessageService,
		UserService:    userService,
		AuthService:    authService,
		logger:         logger,
		valid:          valid,
	}
}

func (h *Handler) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(h.AuthService, h.logger, h.valid))

		r.Get("/", h.GetAllPrivateMessages)
		r.Post("/", h.SendPrivateMessage)

		r.Get("/user/{id}", h.GetAllPrivateMessagesFromUser)
	})

	return router
}

// SendPrivateMessage godoc
//
//	@Summary		Send private message to user
//	@Description	Send private message to user
//	@Security		BasicAuth
//	@Tags			Message
//	@Accept			json
//	@Produce		json
//	@Param			input	body		request.SendPrivateMessageRequest	true	"private message schema"
//	@Success		200		{object}	[]response.PrivateMessageResponse
//	@Failure		401		{string}	Unauthorized
//	@Failure		400		{string}	invalid		message	provided
//	@Failure		500		{string}	internal	error
//	@Router			/api/v1/messages/private [post]
func (h *Handler) SendPrivateMessage(rw http.ResponseWriter, req *http.Request) {
	id, err := handlerutils.GetIntHeaderByKey(req, "id")
	if err != nil {
		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusUnauthorized, "", err.Error())
		return
	}

	var privMsgReq request.SendPrivateMessageRequest

	if err = render.DecodeJSON(req.Body, &privMsgReq); err != nil {
		logMsg := fmt.Sprintf("error occurred validating PrivateMessageRequest struct: %v", err)
		respMsg := fmt.Sprintf("invalid message provided: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	privMsgReq.FromID = id

	if err = privMsgReq.Validate(h.valid); err != nil {
		logMsg := fmt.Sprintf("error occurred validating PrivateMessageRequest struct: %v", err)
		respMsg := fmt.Sprintf("invalid message provided: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	message, err := h.MessageService.SendPrivateMessage(req.Context(), privMsgReq.FromID, privMsgReq.ToID, privMsgReq.Content)

	if err != nil {
		switch {
		case errors.Is(err, messageservice.ErrNoSuchReceiver):
			errMsg := fmt.Sprintf("error occurred sending private message: %s", err)

			handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, errMsg, errMsg)

			return

		default:
			errMsg := fmt.Sprintf("error occurred saving private message: %s", err)

			handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusInternalServerError, errMsg, errMsg)

			return
		}
	}

	render.JSON(rw, req, mapper.MapPrivateMessageToResponse(message))
	rw.WriteHeader(http.StatusCreated)
}

// GetAllPrivateMessages godoc
//
//	@Summary		Get all private messages
//	@Description	Get all private messages that were sent to chat
//	@Security		BasicAuth
//	@Tags			Message
//	@Produce		json
//	@Param			offset	query		int	true	"Offset"
//	@Param			limit	query		int	true	"Limit"
//	@Success		200		{object}	[]response.PrivateMessageResponse
//	@Failure		401		{string}	Unauthorized
//	@Router			/api/v1/messages/private [get]
func (h *Handler) GetAllPrivateMessages(rw http.ResponseWriter, req *http.Request) {
	id, err := handlerutils.GetIntHeaderByKey(req, "id")
	if err != nil {
		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusUnauthorized, "", err.Error())
		return
	}

	paginationOpts := handlerinternalutils.GetPaginationOptsFromQuery(req, handler.DefaultOffset, handler.DefaultLimit)

	if err = paginationOpts.Validate(h.valid); err != nil {
		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, "", err.Error())

		return
	}

	messages := h.MessageService.GetAllPrivateMessages(req.Context(), id, paginationOpts.Offset, paginationOpts.Limit)

	render.JSON(rw, req, sliceutils.Map(messages, mapper.MapPrivateMessageToResponse))
	rw.WriteHeader(http.StatusOK)
}

// GetAllPrivateMessagesFromUser godoc
//
//	@Summary		Get all private messages from user
//	@Description	Get all private messages from user
//	@Security		BasicAuth
//	@Tags			Message
//	@Produce		json
//	@Param			offset	query	int	true	"Offset"
//	@Param			limit	query	int	true	"Limit"
//	@Param			user_id	path	int	true	"User FromID"
//	@Para			page query int true "page"
//	@Success		200	{object}	[]response.PrivateMessageResponse
//	@Failure		401	{string}	Unauthorized
//	@Router			/api/v1/messages/private/user/{user_id} [get]
func (h *Handler) GetAllPrivateMessagesFromUser(rw http.ResponseWriter, req *http.Request) {
	id, err := handlerutils.GetIntHeaderByKey(req, "id")
	if err != nil {
		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusUnauthorized, "", err.Error())
		return
	}

	ctx := req.Context()

	fromID, err := strconv.Atoi(chi.URLParam(req, "id"))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
	}

	paginationOpts := handlerinternalutils.GetPaginationOptsFromQuery(req, handler.DefaultOffset, handler.DefaultLimit)

	if err = paginationOpts.Validate(h.valid); err != nil {
		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, "", err.Error())

		return
	}

	messages, err := h.MessageService.GetAllPrivateMessagesFromUser(ctx, id, fromID, paginationOpts.Offset, paginationOpts.Limit)
	if err != nil {
		msg := fmt.Sprintf("error occurred getting private messages from user: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)

		return
	}

	render.JSON(rw, req, sliceutils.Map(messages, mapper.MapPrivateMessageToResponse))
	rw.WriteHeader(http.StatusOK)
}
