// nolint
package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/dto"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/request"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/middleware"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"

	messagemapper "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/mapper/message"
	messageservice "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service/message"
	handlerutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/handler"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
)

type MessageService interface {
	SendPrivateMessage(ctx context.Context, createModel dto.PrivateMessageDTO) (*model.PrivateMessage, error)
	SendPublicMessage(ctx context.Context, createModel dto.PublicMessageDTO) (*model.PublicMessage, error)

	GetPrivateMessage(ctx context.Context, id int) (*model.PrivateMessage, error)
	GetPublicMessage(ctx context.Context, id int) (*model.PublicMessage, error)

	GetAllPrivateMessages(ctx context.Context, userFromID int, offset, limit int) []*model.PrivateMessage
	GetAllPublicMessages(ctx context.Context, offset, limit int) []*model.PublicMessage

	GetAllPrivateMessagesFromUser(ctx context.Context, toID, fromID int, offset, limit int) ([]*model.PrivateMessage, error)
}

const (
	defaultOffset = 0
	defaultLimit  = 100
)

var (
	ErrInvalidOffset = errors.New("invalid offset provided")
	ErrInvalidLimit  = errors.New("invalid limit provided")
)

type MessageHandler struct {
	MessageService MessageService
	UserService    UserService
	AuthService    middleware.AuthService
	logger         *logrus.Logger
	validator      *validator.Validate
}

func NewMessageHandler(ms MessageService, us UserService, as middleware.AuthService, logger *logrus.Logger) *MessageHandler {
	return &MessageHandler{
		MessageService: ms,
		UserService:    us,
		AuthService:    as,
		logger:         logger,
		validator:      validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (mh *MessageHandler) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(mh.AuthService, mh.logger))

		r.Get("/public", mh.GetAllPublicMessages)
		r.Post("/public", mh.SendPublicMessage)

		r.Get("/private", mh.GetAllPrivateMessages)
		r.Post("/private", mh.SendPrivateMessage)

		r.Get("/private/user/{id}", mh.GetAllPrivateMessagesFromUser)
	})

	return router
}

// GetAllPublicMessages godoc
//
//	@Summary		Get all public messages
//	@Description	Get all public messages that were sent to chat
//	@Security		BasicAuth
//	@Tags			Message
//	@Produce		json
//	@Param			offset	query		int	true	"Offset"
//	@Param			limit	query		int	true	"Limit"
//	@Success		200		{object}	[]response.PublicMessageResponse
//	@Failure		401		{string}	Unauthorized
//	@Router			/messages/public [get]
func (mh *MessageHandler) GetAllPublicMessages(rw http.ResponseWriter, req *http.Request) {
	offset, limit := handlerutils.GetOffsetAndLimitFromQuery(req, defaultOffset, defaultLimit)
	if offset < 0 {
		respMsg := fmt.Sprintf("%s", ErrInvalidOffset)

		handlerutils.WriteErrResponseAndLog(rw, mh.logger, http.StatusBadRequest, "", respMsg)

		return
	}

	if limit < 0 {
		respMsg := fmt.Sprintf("%s", ErrInvalidLimit)

		handlerutils.WriteErrResponseAndLog(rw, mh.logger, http.StatusBadRequest, "", respMsg)

		return
	}

	messages := mh.MessageService.GetAllPublicMessages(req.Context(), offset, limit)

	resp := sliceutils.Map(messages, messagemapper.MapPublicMessageToPublicMsgResp)

	render.JSON(rw, req, resp)

	rw.WriteHeader(http.StatusOK)
}

// SendPublicMessage godoc
//
//	@Summary		Send public message to chat
//	@Description	Send public message to chat
//	@Security		BasicAuth
//	@Tags			Message
//	@Accept			json
//	@Produce		json
//	@Param			input	body		request.SendPublicMessageRequest	true	"public message schema"
//	@Success		200		{object}	[]response.PublicMessageResponse
//	@Failure		401		{string}	Unauthorized
//	@Failure		500		{string}	internal	error
//	@Router			/messages/public [post]
func (mh *MessageHandler) SendPublicMessage(rw http.ResponseWriter, req *http.Request) {
	id, err := handlerutils.GetIntHeaderByKey(req, "id")
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	var pubMsgReq request.SendPublicMessageRequest

	err = render.DecodeJSON(req.Body, &pubMsgReq)
	if err != nil {
		logMsg := fmt.Sprintf("error occurred validating PublicMessageRequest struct: %s", err)
		respMsg := fmt.Sprintf("invalid message provided: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, mh.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	if err = mh.validator.Struct(pubMsgReq); err != nil {
		logMsg := fmt.Sprintf("error occurred validating PublicMessageRequest struct: %s", err)
		respMsg := fmt.Sprintf("invalid message provided: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, mh.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	msg := dto.PublicMessageDTO{
		FromID:  id,
		Content: pubMsgReq.Content,
	}

	message, err := mh.MessageService.SendPublicMessage(req.Context(), msg)
	if err != nil {
		logMsg := fmt.Sprintf("error occurred saving public message: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, mh.logger, http.StatusInternalServerError, logMsg, "")

		return
	}

	resp := messagemapper.MapPublicMessageToPublicMsgResp(message)

	render.JSON(rw, req, resp)
	rw.WriteHeader(http.StatusCreated)
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
func (mh *MessageHandler) SendPrivateMessage(rw http.ResponseWriter, req *http.Request) {
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

	if err = mh.validator.Struct(privMsgReq); err != nil {
		logMsg := fmt.Sprintf("error occurred validating PrivateMessageRequest struct: %s", err)
		respMsg := fmt.Sprintf("invalid message provided: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, mh.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	msg := dto.PrivateMessageDTO{
		FromID:  id,
		ToID:    privMsgReq.ToID,
		Content: privMsgReq.Content,
	}

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

	resp := messagemapper.MapPrivateMessageToPrivateMsgResp(message)

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
//	@Success		200	{object}	[]response.PrivateMessageResponse
//	@Failure		401	{string}	Unauthorized
//	@Router			/messages/private [get]
func (mh *MessageHandler) GetAllPrivateMessages(rw http.ResponseWriter, req *http.Request) {
	id, err := handlerutils.GetIntHeaderByKey(req, "id")
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	offset, limit := handlerutils.GetOffsetAndLimitFromQuery(req, defaultOffset, defaultLimit)
	if offset < 0 {
		respMsg := fmt.Sprintf("%s", ErrInvalidOffset)

		handlerutils.WriteErrResponseAndLog(rw, mh.logger, http.StatusBadRequest, "", respMsg)

		return
	}

	if limit < 0 {
		respMsg := fmt.Sprintf("%s", ErrInvalidLimit)

		handlerutils.WriteErrResponseAndLog(rw, mh.logger, http.StatusBadRequest, "", respMsg)

		return
	}

	messages := mh.MessageService.GetAllPrivateMessages(req.Context(), id, offset, limit)

	resp := sliceutils.Map(messages, messagemapper.MapPrivateMessageToPrivateMsgResp)

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
//	@Param			offset	query		int	true	"Offset"
//	@Param			limit	query		int	true	"Limit"
//	@Param			user_id	path	int	true	"User ID"
//	@Para			page query int true "page"
//	@Success		200	{object}	[]response.PrivateMessageResponse
//	@Failure		401	{string}	Unauthorized
//	@Router			/messages/private/user/{user_id} [get]
func (mh *MessageHandler) GetAllPrivateMessagesFromUser(rw http.ResponseWriter, req *http.Request) {
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

	offset, limit := handlerutils.GetOffsetAndLimitFromQuery(req, defaultOffset, defaultLimit)
	if offset < 0 {
		respMsg := fmt.Sprintf("%s", ErrInvalidOffset)

		handlerutils.WriteErrResponseAndLog(rw, mh.logger, http.StatusBadRequest, "", respMsg)

		return
	}

	if limit < 0 {
		respMsg := fmt.Sprintf("%s", ErrInvalidLimit)

		handlerutils.WriteErrResponseAndLog(rw, mh.logger, http.StatusBadRequest, "", respMsg)

		return
	}

	messages, err := mh.MessageService.GetAllPrivateMessagesFromUser(ctx, id, fromID, offset, limit)
	if err != nil {
		logMsg := fmt.Sprintf("error occurred getting private messages from user: %s", err)
		respMsg := fmt.Sprintf("error occurred getting private messages from user: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, mh.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	resp := sliceutils.Map(messages, messagemapper.MapPrivateMessageToPrivateMsgResp)

	render.JSON(rw, req, resp)
	rw.WriteHeader(http.StatusOK)
}
