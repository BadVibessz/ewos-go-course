package handler

import (
	"context"
	"fmt"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/dto"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/request"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/response"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/middleware"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"strconv"

	handlerutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/utils/handler"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"

	"net/http"
)

type MessageService interface {
	SendPrivateMessage(ctx context.Context, createModel dto.CreatePrivateMessageDTO) (*model.PrivateMessage, error)
	SendPublicMessage(ctx context.Context, createModel dto.CreatePublicMessageDTO) (*model.PublicMessage, error)

	GetPrivateMessage(ctx context.Context, id int) (*model.PrivateMessage, error)
	GetPublicMessage(ctx context.Context, id int) (*model.PublicMessage, error)

	GetAllPrivateMessages(ctx context.Context, user *model.User) []*model.PrivateMessage
	GetAllPublicMessages(ctx context.Context) []*model.PublicMessage

	GetAllPrivateMessagesFromUser(ctx context.Context, user *model.User, id int) []*model.PrivateMessage

	UpdatePrivateMessage(ctx context.Context, id int, newContent string) (*model.PrivateMessage, error)
	UpdatePublicMessage(ctx context.Context, id int, newContent string) (*model.PublicMessage, error)

	DeletePrivateMessage(ctx context.Context, id int) (*model.PrivateMessage, error)
	DeletePublicMessage(ctx context.Context, id int) (*model.PublicMessage, error)
}

const defaultLimit = 10
const defaultPage = 1

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

func (mh *MessageHandler) Routes() chi.Router {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(mh.AuthService))

		r.Get("/public", mh.GetAllPublicMessages)
		r.Post("/public", mh.SendPublicMessage)

		r.Get("/private", mh.GetAllPrivateMessages)
		r.Post("/private", mh.SendPrivateMessage)

		r.Get("/private/user/{id}", mh.GetAllPrivateMessagesFromUser)
	})

	return router
}

// GetAllPublicMessages godoc
// @Summary      Get all public messages
// @Description  Get all public messages that were sent to chat
// @Security 	 BasicAuth
// @Tags         Message
// @Produce      json
// @Success      200  {object}  []response.PublicMessageResponse
// @Failure 	 401  {string}  Unauthorized
// @Router       /messages/public [get]
func (mh *MessageHandler) GetAllPublicMessages(w http.ResponseWriter, req *http.Request) {
	messages := mh.MessageService.GetAllPublicMessages(req.Context())
	messages = handlerutils.Paginate(req, defaultPage, defaultLimit, messages)

	// TODO: TEST!!
	resp := sliceutils.Map(messages, func(msg *model.PublicMessage) response.PublicMessageResponse {
		return response.PublicMsgRespFromMessage(*msg)
	})

	render.JSON(w, req, resp)

	w.WriteHeader(http.StatusOK)
}

// SendPublicMessage godoc
// @Summary      Send public message to chat
// @Description  Send public message to chat
// @Security 	 BasicAuth
// @Tags         Message
// @Accept	  	 json
// @Produce      json
// @Param input  body request.SendPublicMessageRequest true "public message schema"
// @Success      200  {object}  []response.PublicMessageResponse
// @Failure 	 401  {string}  Unauthorized
// @Router       /messages/public [post]
func (mh *MessageHandler) SendPublicMessage(rw http.ResponseWriter, req *http.Request) {
	user, ok := req.Context().Value("user").(model.User)
	if !ok { // unauthorized
		rw.WriteHeader(http.StatusUnauthorized)
	}

	var pubMsgReq request.SendPublicMessageRequest

	err := render.DecodeJSON(req.Body, &pubMsgReq)
	if err != nil {
		logMsg := fmt.Sprintf("error occurred validating PublicMessageRequest struct: %s", err)
		respMsg := fmt.Sprintf("invalid message provided")

		handlerutils.WriteResponseAndLogError(rw, mh.logger, http.StatusBadRequest, logMsg, respMsg)

		rw.WriteHeader(http.StatusBadRequest)

		return
	}

	if err = mh.validator.Struct(pubMsgReq); err != nil {
		logMsg := fmt.Sprintf("error occurred validating PublicMessageRequest struct: %s", err)
		respMsg := fmt.Sprintf("invalid message provided")

		handlerutils.WriteResponseAndLogError(rw, mh.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	msg := dto.CreatePublicMessageDTO{ // todo: maybe dto should use mode.User struct?
		FromID:  user.ID,
		Content: pubMsgReq.Content,
	}

	message, err := mh.MessageService.SendPublicMessage(req.Context(), msg)
	if err != nil {
		logMsg := fmt.Sprintf("error occurred saving public message: %s", err)

		handlerutils.WriteResponseAndLogError(rw, mh.logger, http.StatusInternalServerError, logMsg, "") // todo: correct status code?

		return
	}

	resp := response.PublicMsgRespFromMessage(*message)

	render.JSON(rw, req, resp)
	rw.WriteHeader(http.StatusCreated)
}

// SendPrivateMessage godoc
// @Summary      Send private message to user
// @Description  Send private message to user
// @Security 	 BasicAuth
// @Tags         Message
// @Accept	  	 json
// @Produce      json
// @Param input  body request.SendPrivateMessageRequest true "private message schema"
// @Success      200  {object}  []response.PrivateMessageResponse
// @Failure 	 401  {string}  Unauthorized
// @Failure 	 400  {string}  invalid message provided
// @Router       /messages/private [post]
func (mh *MessageHandler) SendPrivateMessage(rw http.ResponseWriter, req *http.Request) {
	user, ok := req.Context().Value("user").(model.User)
	if !ok { // unauthorized
		rw.WriteHeader(http.StatusUnauthorized)
	}

	var privMsgReq request.SendPrivateMessageRequest

	err := render.DecodeJSON(req.Body, &privMsgReq)
	if err != nil {
		logMsg := fmt.Sprintf("error occurred validating PrivateMessageRequest struct: %s", err)
		respMsg := fmt.Sprintf("invalid message provided")

		handlerutils.WriteResponseAndLogError(rw, mh.logger, http.StatusBadRequest, logMsg, respMsg)

		rw.WriteHeader(http.StatusBadRequest)

		return
	}

	if err = mh.validator.Struct(privMsgReq); err != nil {
		logMsg := fmt.Sprintf("error occurred validating PrivateMessageRequest struct: %s", err)
		respMsg := fmt.Sprintf("invalid message provided")

		handlerutils.WriteResponseAndLogError(rw, mh.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	msg := dto.CreatePrivateMessageDTO{ // todo: maybe dto should use mode.User struct?
		FromID:  user.ID,
		ToID:    privMsgReq.ToID,
		Content: privMsgReq.Content,
	}

	message, err := mh.MessageService.SendPrivateMessage(req.Context(), msg)
	if err != nil {
		logMsg := fmt.Sprintf("error occurred saving private message: %s", err)

		handlerutils.WriteResponseAndLogError(rw, mh.logger, http.StatusInternalServerError, logMsg, "") // todo: correct status code?

		return
	}

	resp := response.PrivateMsgRespFromMessage(*message)

	render.JSON(rw, req, resp)
	rw.WriteHeader(http.StatusCreated)
}

// GetAllPrivateMessages godoc
// @Summary      Get all private messages
// @Description  Get all private messages that were sent to chat
// @Security 	 BasicAuth
// @Tags         Message
// @Produce      json
// @Success      200  {object}  []response.PrivateMessageResponse
// @Failure 	 401  {string}  Unauthorized
// @Router       /messages/private [get]
func (mh *MessageHandler) GetAllPrivateMessages(rw http.ResponseWriter, req *http.Request) {
	user, ok := req.Context().Value("user").(model.User)
	if !ok { // unauthorized
		rw.WriteHeader(http.StatusUnauthorized)
	}

	messages := mh.MessageService.GetAllPrivateMessages(req.Context(), &user)
	messages = handlerutils.Paginate(req, defaultPage, defaultLimit, messages)

	resp := sliceutils.Map(messages, func(msg *model.PrivateMessage) response.PrivateMessageResponse {
		return response.PrivateMsgRespFromMessage(*msg)
	})

	render.JSON(rw, req, resp)
	rw.WriteHeader(http.StatusOK)
}

// GetAllPrivateMessagesFromUser godoc
// @Summary      Get all private messages from user
// @Description  Get all private messages from user
// @Security 	 BasicAuth
// @Tags         Message
// @Produce      json
// @Param 		 user_id  path int true "User ID"
// @Para		 page query int true "page"
// @Success      200  {object}  []response.PrivateMessageResponse
// @Failure 	 401  {string}  Unauthorized
// @Router       /messages/private/user/{user_id} [get]
func (mh *MessageHandler) GetAllPrivateMessagesFromUser(rw http.ResponseWriter, req *http.Request) {
	user, ok := req.Context().Value("user").(model.User)
	if !ok { // unauthorized
		rw.WriteHeader(http.StatusUnauthorized)
	}

	ctx := req.Context()

	fromID, err := strconv.Atoi(chi.URLParam(req, "id"))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
	}

	messages := mh.MessageService.GetAllPrivateMessagesFromUser(ctx, &user, fromID)
	messages = handlerutils.Paginate(req, defaultPage, defaultLimit, messages)

	resp := sliceutils.Map(messages, func(msg *model.PrivateMessage) response.PrivateMessageResponse {
		return response.PrivateMsgRespFromMessage(*msg)
	})

	render.JSON(rw, req, resp)
	rw.WriteHeader(http.StatusOK)
}
