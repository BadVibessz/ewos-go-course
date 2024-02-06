package handler

import (
	"context"
	"fmt"
	"math"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/dto"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/request"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/middleware"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"

	usermapper "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/mapper/user"
	handlerutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/handler"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
)

type UserService interface {
	RegisterUser(ctx context.Context, user dto.UserDTO) (*model.User, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetAllUsers(ctx context.Context, offset, limit int) []*model.User
	UpdateUser(ctx context.Context, id int, updateModel dto.UserDTO) (*model.User, error)
	DeleteUser(ctx context.Context, id int) (*model.User, error)
}

type UserHandler struct {
	UserService    UserService
	MessageService MessageService
	AuthService    middleware.AuthService
	logger         *logrus.Logger
	validator      *validator.Validate
}

func NewUserHandler(us UserService, ms MessageService, as middleware.AuthService, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		UserService:    us,
		MessageService: ms,
		AuthService:    as,
		logger:         logger,
		validator:      validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (uh *UserHandler) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Post("/register", uh.Register)
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(uh.AuthService, uh.logger))
		r.Get("/all", uh.GetAll)
		r.Get("/messages", uh.GetAllUsersThatSentMessage)
	})

	return router
}

// Register godoc
//
//	@Summary		Register new user
//	@Description	to register new user
//	@Tags			User
//	@Accept			json
//	@Produce		plain
//	@Param			input	body		request.RegisterRequest	true	"registration info"
//	@Success		200		{object}	response.UserResponse
//	@Failure		400		{string}	invalid		registration	data	provided
//	@Failure		500		{string}	internal	error
//	@Router			/users/register [post]
func (uh *UserHandler) Register(rw http.ResponseWriter, req *http.Request) {
	registerReq := request.RegisterRequest{}

	err := render.DecodeJSON(req.Body, &registerReq)
	if err != nil {
		logMsg := fmt.Sprintf("error occurred decoding request body to RegisterRequest struct: %s", err)
		respMsg := fmt.Sprintf("invalid registration data provided: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, uh.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	err = uh.validator.Struct(registerReq)
	if err != nil {
		logMsg := fmt.Sprintf("error occurred validating RegisterRequest struct: %s", err)
		respMsg := fmt.Sprintf("invalid registration data provided: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, uh.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	// todo: understand what is bcrypt.cost!
	hash, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		logMsg := fmt.Sprintf("error occurred generating hash from password: %s", err)
		respMsg := fmt.Sprintf("error occurred generating hash from password: %s", err)

		handlerutils.WriteErrResponseAndLog(rw, uh.logger, http.StatusInternalServerError, logMsg, respMsg)
	}

	createModel := dto.UserDTO{
		Email:          registerReq.Email,
		Username:       registerReq.Username,
		HashedPassword: string(hash),
	}

	ctx := req.Context()

	user, err := uh.UserService.RegisterUser(ctx, createModel)
	if err != nil {
		handlerutils.WriteErrResponseAndLog(rw, uh.logger, http.StatusBadRequest,
			"", fmt.Sprintf("invalid registration data provided: %s", err))

		return
	}

	render.JSON(rw, req, usermapper.MapUserToUserResponse(user))
	rw.WriteHeader(http.StatusCreated)
}

// GetAll godoc
//
//	@Summary		Get all users
//	@Description	Get all users
//	@Security		BasicAuth
//	@Tags			User
//	@Produce		json
//	@Success		200	{object}	[]response.UserResponse
//	@Failure		401	{string}	Unauthorized
//	@Router			/users/all [get]
func (uh *UserHandler) GetAll(rw http.ResponseWriter, req *http.Request) {
	offset, limit := handlerutils.GetOffsetAndLimitFromQuery(req, defaultOffset, defaultLimit)

	users := uh.UserService.GetAllUsers(req.Context(), offset, limit)
	resp := sliceutils.Map(users, usermapper.MapUserToUserResponse)

	render.JSON(rw, req, resp)

	rw.WriteHeader(http.StatusOK)
}

// GetAllUsersThatSentMessage godoc
//
//	@Summary		Get all users that sent message to current user
//	@Description	Get all users that sent message to current user
//	@Security		BasicAuth
//	@Tags			User
//	@Produce		json
//	@Success		200	{object}	[]response.UserResponse
//	@Failure		401	{string}	Unauthorized
//	@Router			/users/messages [get]
func (uh *UserHandler) GetAllUsersThatSentMessage(rw http.ResponseWriter, req *http.Request) {
	id, err := handlerutils.GetIntHeaderByKey(req, "id")
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	messages := uh.MessageService.GetAllPrivateMessages(req.Context(), id, 0, math.MaxInt64)

	usersDuplicates := sliceutils.Map(messages, func(msg *model.PrivateMessage) *model.User { return msg.From })
	users := sliceutils.Unique(usersDuplicates)

	resp := sliceutils.Map(users, usermapper.MapUserToUserResponse)

	render.JSON(rw, req, resp)

	rw.WriteHeader(http.StatusOK)
}
