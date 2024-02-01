package handler

import (
	"context"
	"fmt"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/dto"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/requset"
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

func (mh *MessageHandler) GetAllPublicMessages(w http.ResponseWriter, req *http.Request) {
	messages := mh.MessageService.GetAllPublicMessages(req.Context())

	page, _ := strconv.Atoi(req.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	if limit == 0 {
		limit = defaultLimit
	}

	leftBound := page*limit - limit
	rightBound := leftBound + limit - 1

	if rightBound >= len(messages) {
		rightBound = len(messages) - 1
	}

	messages = messages[leftBound:rightBound]

	resp := sliceutils.Map(messages[leftBound:rightBound+1], func(msg *model.PublicMessage) response.PublicMessageResponse {
		return response.PublicMsgRespFromMessage(*msg)
	})

	render.JSON(w, req, resp)

	w.WriteHeader(http.StatusOK)
}

func (mh *MessageHandler) SendPublicMessage(rw http.ResponseWriter, req *http.Request) {
	user, ok := req.Context().Value("user").(model.User)
	if !ok { // unauthorized
		rw.WriteHeader(http.StatusUnauthorized)
	}

	var pubMsgReq requset.SendPublicMessageRequest

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

func (mh *MessageHandler) SendPrivateMessage(rw http.ResponseWriter, req *http.Request) {
	user, ok := req.Context().Value("user").(model.User)
	if !ok { // unauthorized
		rw.WriteHeader(http.StatusUnauthorized)
	}

	var privMsgReq requset.SendPrivateMessageRequest

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

func (mh *MessageHandler) GetAllPrivateMessages(rw http.ResponseWriter, req *http.Request) {
	user, ok := req.Context().Value("user").(model.User)
	if !ok { // unauthorized
		rw.WriteHeader(http.StatusUnauthorized)
	}

	messages := mh.MessageService.GetAllPrivateMessages(req.Context(), &user)

	page, _ := strconv.Atoi(req.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	if limit == 0 {
		limit = defaultLimit
	}

	leftBound := page*limit - limit
	rightBound := leftBound + limit - 1

	if rightBound >= len(messages) {
		rightBound = len(messages) - 1
	}

	messages = messages[leftBound:rightBound]

	resp := sliceutils.Map(messages, func(msg *model.PrivateMessage) response.PrivateMessageResponse {
		return response.PrivateMsgRespFromMessage(*msg)
	})

	render.JSON(rw, req, resp)
	rw.WriteHeader(http.StatusOK)
}

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

	resp := sliceutils.Map(messages, func(msg *model.PrivateMessage) response.PrivateMessageResponse {
		return response.PrivateMsgRespFromMessage(*msg)

	})

	render.JSON(rw, req, resp)
	rw.WriteHeader(http.StatusOK)
}
