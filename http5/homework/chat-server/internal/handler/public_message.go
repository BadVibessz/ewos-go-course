// nolint
package handler

import (
	"context"
	"fmt"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/mapper"
	handlerinternalutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/utils/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/request"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/middleware"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"

	handlerutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/handler"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
)

type PublicMessageService interface {
	SendPublicMessage(ctx context.Context, createModel entity.PublicMessage) (*model.PublicMessage, error)
	GetPublicMessage(ctx context.Context, id int) (*model.PublicMessage, error)
	GetAllPublicMessages(ctx context.Context, paginationOpts request.PaginationOptions) []*model.PublicMessage
}

type PublicMessageHandler struct {
	MessageService PublicMessageService
	UserService    UserService
	AuthService    middleware.AuthService
	logger         *logrus.Logger
}

func NewPublicMessageHandler(ms PublicMessageService, us UserService, as middleware.AuthService, logger *logrus.Logger) *PublicMessageHandler {
	return &PublicMessageHandler{
		MessageService: ms,
		UserService:    us,
		AuthService:    as,
		logger:         logger,
	}
}

func (mh *PublicMessageHandler) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(mh.AuthService, mh.logger))

		r.Get("/", mh.GetAllPublicMessages)
		r.Post("/", mh.SendPublicMessage)
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
func (mh *PublicMessageHandler) GetAllPublicMessages(rw http.ResponseWriter, req *http.Request) {
	paginationOpts := handlerinternalutils.GetPaginationOptsFromQuery(req, defaultOffset, defaultLimit)

	err := paginationOpts.Validate()
	if err != nil {
		respMsg := fmt.Sprintf("%v", err)

		handlerutils.WriteErrResponseAndLog(rw, mh.logger, http.StatusBadRequest, "", respMsg)

		return
	}

	messages := mh.MessageService.GetAllPublicMessages(req.Context(), paginationOpts)

	resp := sliceutils.Map(messages, mapper.MapPublicMessageToResponse)

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
func (mh *PublicMessageHandler) SendPublicMessage(rw http.ResponseWriter, req *http.Request) {
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

	pubMsgReq.FromID = id

	if err = pubMsgReq.Validate(); err != nil {
		logMsg := fmt.Sprintf("error occurred validating PublicMessageRequest struct: %s", err)
		respMsg := fmt.Sprintf("invalid message provided: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, mh.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	msg := mapper.MapPublicMessageRequestToEntity(&pubMsgReq)

	message, err := mh.MessageService.SendPublicMessage(req.Context(), msg)
	if err != nil {
		logMsg := fmt.Sprintf("error occurred saving public message: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, mh.logger, http.StatusInternalServerError, logMsg, "")

		return
	}

	resp := mapper.MapPublicMessageToResponse(message)

	render.JSON(rw, req, resp)
	rw.WriteHeader(http.StatusCreated)
}
