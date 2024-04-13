// nolint
package private

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/mapper"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/request"

	messageservice "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service/message"

	handlerinternalutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/utils/handler"
	handlerutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/handler"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
)

type MessageService interface {
	SendPrivateMessage(ctx context.Context, msg entity.PrivateMessage) (*entity.PrivateMessage, error)
	GetAllPrivateMessages(ctx context.Context, toUsername string, offset, limit int) []*entity.PrivateMessage
	GetAllPrivateMessagesFromUser(ctx context.Context, toUsername, fromUsername string, offset, limit int) ([]*entity.PrivateMessage, error)
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

type Middleware = func(http.Handler) http.Handler

type Handler struct {
	MessageService MessageService
	UserService    UserService
	Middlewares    []Middleware
	logger         *logrus.Logger
	validator      *validator.Validate
}

func New(
	privateMessageService MessageService,
	userService UserService,
	logger *logrus.Logger,
	validator *validator.Validate,
	middlewares ...Middleware,
) *Handler {
	return &Handler{
		MessageService: privateMessageService,
		UserService:    userService,
		Middlewares:    middlewares,
		logger:         logger,
		validator:      validator,
	}
}

func (h *Handler) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(h.Middlewares...)

		r.Get("/", h.GetAllPrivateMessages)
		r.Post("/", h.SendPrivateMessage)

		r.Get("/user", h.GetAllPrivateMessagesFromUser)
	})

	return router
}

func switchByErrorAndWriteResponse(err error, rw http.ResponseWriter, logger *logrus.Logger) {
	switch {
	case errors.Is(err, messageservice.ErrNoSuchReceiver):
		errMsg := fmt.Sprintf("error occurred sending private message: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusBadRequest, errMsg, errMsg)

	default:
		errMsg := fmt.Sprintf("error occurred saving private message: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, logger, http.StatusInternalServerError, errMsg, errMsg)
	}
}

// SendPrivateMessage godoc
//
//	@Summary		Send private message to user
//	@Description	Send private message to user
//	@Security		BasicAuth
//	@Security		JWT
//	@Tags			Message
//	@Accept			json
//	@Produce		json
//	@Param			input	body		request.SendPrivateMessageRequest	true	"private message schema"
//	@Success		200		{object}	[]response.GetPrivateMessageResponse
//	@Failure		401		{string}	Unauthorized
//	@Failure		400		{string}	invalid		message	provided
//	@Failure		500		{string}	internal	error
//	@Router			/api/v1/messages/private [post]
func (h *Handler) SendPrivateMessage(rw http.ResponseWriter, req *http.Request) {
	username, err := handlerutils.GetStringHeaderByKey(req, "username")
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

	if err = privMsgReq.Validate(h.validator); err != nil {
		logMsg := fmt.Sprintf("error occurred validating PrivateMessageRequest struct: %v", err)
		respMsg := fmt.Sprintf("invalid message provided: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	message, err := h.MessageService.SendPrivateMessage(req.Context(), mapper.MapSendPrivateMessageRequestToEntity(privMsgReq, username))
	if err != nil {
		switchByErrorAndWriteResponse(err, rw, h.logger)
	}

	render.JSON(rw, req, mapper.MapPrivateMessageToResponse(message))
	rw.WriteHeader(http.StatusCreated)
}

// GetAllPrivateMessages godoc
//
//	@Summary		Get all private messages
//	@Description	Get all private messages that were sent to chat
//	@Security		BasicAuth
//	@Security		JWT
//	@Tags			Message
//	@Produce		json
//	@Param			offset	query		int	true	"Offset"
//	@Param			limit	query		int	true	"Limit"
//	@Success		200		{object}	[]response.GetPrivateMessageResponse
//	@Failure		401		{string}	Unauthorized
//	@Router			/api/v1/messages/private [get]
func (h *Handler) GetAllPrivateMessages(rw http.ResponseWriter, req *http.Request) {
	username, err := handlerutils.GetStringHeaderByKey(req, "username")
	if err != nil {
		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusUnauthorized, "", err.Error())
		return
	}

	paginationOpts := handlerinternalutils.GetPaginationOptsFromQuery(req, handler.DefaultOffset, handler.DefaultLimit)

	if err = paginationOpts.Validate(h.validator); err != nil {
		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, "", err.Error())

		return
	}

	messages := h.MessageService.GetAllPrivateMessages(req.Context(), username, paginationOpts.Offset, paginationOpts.Limit)

	render.JSON(rw, req, sliceutils.Map(messages, mapper.MapPrivateMessageToResponse))
	rw.WriteHeader(http.StatusOK)
}

// GetAllPrivateMessagesFromUser godoc
//
//	@Summary		Get all private messages from user
//	@Description	Get all private messages from user
//	@Security		BasicAuth
//	@Security		JWT
//	@Tags			Message
//	@Produce		json
//	@Param			offset			query		int		true	"Offset"
//	@Param			limit			query		int		true	"Limit"
//	@Param			from_username	query		string	true	"from_username"
//	@Success		200				{object}	[]response.GetPrivateMessageResponse
//	@Failure		401				{string}	Unauthorized
//	@Router			/api/v1/messages/private/user [get]
func (h *Handler) GetAllPrivateMessagesFromUser(rw http.ResponseWriter, req *http.Request) {
	username, err := handlerutils.GetStringHeaderByKey(req, "username")
	if err != nil {
		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusUnauthorized, "", err.Error())
		return
	}

	ctx := req.Context() // TODO: method not working!

	fromUsername, err := handlerutils.GetStringParamFromQuery(req, "from_username")
	if err != nil {
		msg := fmt.Sprintf("error occurred retrieving from_username query param from request: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)
		return
	}

	paginationOpts := handlerinternalutils.GetPaginationOptsFromQuery(req, handler.DefaultOffset, handler.DefaultLimit)

	if err = paginationOpts.Validate(h.validator); err != nil {
		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, "", err.Error())
		return
	}

	messages, err := h.MessageService.GetAllPrivateMessagesFromUser(ctx, username, fromUsername, paginationOpts.Offset, paginationOpts.Limit)
	if err != nil {
		msg := fmt.Sprintf("error occurred getting private messages from user: %v", err)

		handlerutils.WriteErrResponseAndLog(rw, h.logger, http.StatusBadRequest, msg, msg)

		return
	}

	render.JSON(rw, req, sliceutils.Map(messages, mapper.MapPrivateMessageToResponse))
	rw.WriteHeader(http.StatusOK)
}
