// nolint
package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/mapper"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/request"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/middleware"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"

	messageservice "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service/message"

	handlerinternalutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/utils/handler"
	handlerutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/handler"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
)

type PrivateMessageService interface {
	SendPrivateMessage(ctx context.Context, createModel entity.PrivateMessage) (*model.PrivateMessage, error)
	GetPrivateMessage(ctx context.Context, id int) (*model.PrivateMessage, error)
	GetAllPrivateMessages(ctx context.Context, userToID int, paginationOpts request.PaginationOptions) []*model.PrivateMessage
	GetAllPrivateMessagesFromUser(ctx context.Context, toID, fromID int, paginationOpts request.PaginationOptions) ([]*model.PrivateMessage, error)
}

type PrivateMessageHandler struct {
	MessageService PrivateMessageService
	UserService    UserService
	AuthService    middleware.AuthService
	logger         *logrus.Logger
}

func NewPrivateMessageHandler(ms PrivateMessageService, us UserService, as middleware.AuthService, logger *logrus.Logger) *PrivateMessageHandler {
	return &PrivateMessageHandler{
		MessageService: ms,
		UserService:    us,
		AuthService:    as,
		logger:         logger,
	}
}

func (mh *PrivateMessageHandler) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(mh.AuthService, mh.logger))

		r.Get("/", mh.GetAllPrivateMessages)
		r.Post("/", mh.SendPrivateMessage)

		r.Get("/user/{id}", mh.GetAllPrivateMessagesFromUser)
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
//	@Router			/messages/private [post]
func (mh *PrivateMessageHandler) SendPrivateMessage(rw http.ResponseWriter, req *http.Request) {
	id, err := handlerutils.GetIntHeaderByKey(req, "id")
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	var privMsgReq request.SendPrivateMessageRequest

	err = render.DecodeJSON(req.Body, &privMsgReq)
	if err != nil {
		logMsg := fmt.Sprintf("error occurred validating PrivateMessageRequest struct: %s", err)
		respMsg := fmt.Sprintf("invalid message provided: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, mh.logger, http.StatusBadRequest, logMsg, respMsg)

		rw.WriteHeader(http.StatusBadRequest)

		return
	}

	privMsgReq.FromID = id

	if err = privMsgReq.Validate(); err != nil {
		logMsg := fmt.Sprintf("error occurred validating PrivateMessageRequest struct: %s", err)
		respMsg := fmt.Sprintf("invalid message provided: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, mh.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	msg := mapper.MapPrivateMessageRequestToEntity(&privMsgReq)

	message, err := mh.MessageService.SendPrivateMessage(req.Context(), msg)
	if errors.Is(err, messageservice.ErrNoSuchReceiver) {
		logMsg := fmt.Sprintf("error occurred sending private message: %s", err)
		respMsg := fmt.Sprintf("error occurred sending private message: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, mh.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	} else if err != nil {
		logMsg := fmt.Sprintf("error occurred saving private message: %s", err)
		respMsg := fmt.Sprintf("error occurred saving private message: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, mh.logger, http.StatusInternalServerError, logMsg, respMsg)

		return
	}

	resp := mapper.MapPrivateMessageToResponse(message)

	render.JSON(rw, req, resp)
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
//	@Router			/messages/private [get]
func (mh *PrivateMessageHandler) GetAllPrivateMessages(rw http.ResponseWriter, req *http.Request) {
	id, err := handlerutils.GetIntHeaderByKey(req, "id")
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	paginationOpts := handlerinternalutils.GetPaginationOptsFromQuery(req, defaultOffset, defaultLimit)

	err = paginationOpts.Validate()
	if err != nil {
		respMsg := fmt.Sprintf("%v", err)

		handlerutils.WriteErrResponseAndLog(rw, mh.logger, http.StatusBadRequest, "", respMsg)

		return
	}

	messages := mh.MessageService.GetAllPrivateMessages(req.Context(), id, paginationOpts)

	resp := sliceutils.Map(messages, mapper.MapPrivateMessageToResponse)

	render.JSON(rw, req, resp)
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
//	@Router			/messages/private/user/{user_id} [get]
func (mh *PrivateMessageHandler) GetAllPrivateMessagesFromUser(rw http.ResponseWriter, req *http.Request) {
	id, err := handlerutils.GetIntHeaderByKey(req, "id")
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx := req.Context()

	fromID, err := strconv.Atoi(chi.URLParam(req, "id"))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
	}

	paginationOpts := handlerinternalutils.GetPaginationOptsFromQuery(req, defaultOffset, defaultLimit)

	err = paginationOpts.Validate()
	if err != nil {
		respMsg := fmt.Sprintf("%v", err)

		handlerutils.WriteErrResponseAndLog(rw, mh.logger, http.StatusBadRequest, "", respMsg)

		return
	}

	messages, err := mh.MessageService.GetAllPrivateMessagesFromUser(ctx, id, fromID, paginationOpts)
	if err != nil {
		logMsg := fmt.Sprintf("error occurred getting private messages from user: %s", err)
		respMsg := fmt.Sprintf("error occurred getting private messages from user: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, mh.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	resp := sliceutils.Map(messages, mapper.MapPrivateMessageToResponse)

	render.JSON(rw, req, resp)
	rw.WriteHeader(http.StatusOK)
}
