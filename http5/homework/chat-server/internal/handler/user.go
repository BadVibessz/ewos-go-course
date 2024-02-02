package handler

import (
	"context"
	"fmt"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/dto"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/request"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/response"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/middleware"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"
	handlerutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/utils/handler"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UserService interface {
	RegisterUser(ctx context.Context, user dto.CreateUserDTO) (*model.User, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetAllUsers(ctx context.Context) []*model.User
	UpdateUser(ctx context.Context, id int, updateModel dto.UpdateUserDTO) (*model.User, error)
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

// Routes todo: maybe here accept context and then use callback with context?
func (uh *UserHandler) Routes() chi.Router {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Post("/register", uh.Register)
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(uh.AuthService))
		r.Get("/all", uh.GetAll)
		r.Get("/messages", uh.GetAllUsersThatSentMessage)
	})

	return router
}

// Register godoc
// @Summary      Register new user
// @Description  to register new user
// @Tags         User
// @Accept       json
// @Produce      plain
// @Param input  body request.RegisterRequest true "registration info"
// @Success      200  {object}  response.UserResponse
// @Failure 	 400 {string}	invalid registration data provided
// @Router       /users/register [post]
func (uh *UserHandler) Register(rw http.ResponseWriter, req *http.Request) { // TODO: PANIC IF TRYING REGISTER 2nd USER
	registerReq := request.RegisterRequest{}
	err := render.DecodeJSON(req.Body, &registerReq)
	if err != nil {
		logMsg := fmt.Sprintf("error occured decoding request body to RegisterRequest struct: %s", err)
		respMsg := fmt.Sprintf("invalid registration data provided")

		handlerutils.WriteResponseAndLogError(rw, uh.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	err = uh.validator.Struct(registerReq)
	if err != nil {
		logMsg := fmt.Sprintf("error occurred validating RegisterRequest struct: %s", err)
		respMsg := fmt.Sprintf("invalid registration data provided")

		handlerutils.WriteResponseAndLogError(rw, uh.logger, http.StatusBadRequest, logMsg, respMsg)

		return
	}

	// todo: understand what is bcrypt.cost!
	hash, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)

	createModel := dto.CreateUserDTO{
		Email:          registerReq.Email,
		Username:       registerReq.Username,
		HashedPassword: string(hash),
	}

	// todo: understand how to manage ctx here
	ctx := context.Background()

	user, err := uh.UserService.RegisterUser(ctx, createModel)
	if err != nil {
		handlerutils.WriteResponseAndLogError(rw, uh.logger, http.StatusBadRequest, // todo: specify what user did wrong
			"", "invalid registration data provided")

		return
	}

	render.JSON(rw, req, response.UserToUserResponse(user))
	rw.WriteHeader(http.StatusCreated)
}

// GetAll godoc
// @Summary      Get all users
// @Description  Get all users
// @Security 	 BasicAuth
// @Tags         User
// @Produce      json
// @Success      200  {object}  []response.UserResponse
// @Router       /users/all [get]
func (uh *UserHandler) GetAll(rw http.ResponseWriter, req *http.Request) {
	users := uh.UserService.GetAllUsers(req.Context())
	resp := sliceutils.Map(users, func(user *model.User) response.UserResponse { return response.UserToUserResponse(user) })

	render.JSON(rw, req, resp)

	rw.WriteHeader(http.StatusOK)
}

func (uh *UserHandler) GetAllUsersThatSentMessage(rw http.ResponseWriter, req *http.Request) {
	user, ok := req.Context().Value("user").(model.User)
	if !ok { // unauthorized
		rw.WriteHeader(http.StatusUnauthorized)
	}

	messages := uh.MessageService.GetAllPrivateMessages(req.Context(), &user)

	usersDuplicates := sliceutils.Map(messages, func(msg *model.PrivateMessage) *model.User { return msg.From })
	users := sliceutils.Unique(usersDuplicates)

	resp := sliceutils.Map(users, func(user *model.User) response.UserResponse { return response.UserToUserResponse(user) })

	render.JSON(rw, req, resp)

	rw.WriteHeader(http.StatusOK)
}
