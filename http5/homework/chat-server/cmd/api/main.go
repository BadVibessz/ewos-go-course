package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/config"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/domain/entity"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/fixtures"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/router"

	middlewares "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/middleware"

	authhandler "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/auth"
	privatemessagehandler "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/message/private"
	publicmessagehandler "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/message/public"
	userhandler "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler/user"

	authservice "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service/auth"
	privatemessageservice "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service/message/private"
	publicmessageservice "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service/message/public"
	userservice "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service/user"

	inmemoryrepository "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository/in-memory"
	postgresrepo "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository/postgres"

	inmemory "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/db/in-memory"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/docs"
	_ "github.com/jackc/pgx/v5/stdlib"
)

//	@title			Chat API
//	@version		1.0
//	@description	API Server for Web Chat

//	@BasePath	/chat

//	@securityDefinitions.basic	BasicAuth
//	@securityDefinitions.apikey	JWT
//	@in							header
//	@name						Authorization

const ( // todo: config file
	dbSavePath = "http5/homework/chat-server/internal/db/db_state.json"
	configPath = "http5/homework/chat-server/config"

	port         = 5000
	loadFixtures = true
)

type UserRepo interface {
	AddUser(ctx context.Context, user entity.User) (*entity.User, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	GetAllUsers(ctx context.Context, offset, limit int) []*entity.User
	DeleteUser(ctx context.Context, id int) (*entity.User, error)
	UpdateUser(ctx context.Context, id int, updateModel entity.User) (*entity.User, error)
	CheckUniqueConstraints(ctx context.Context, email, username string) error
}

type PublicMessageRepo interface {
	AddPublicMessage(ctx context.Context, msg entity.PublicMessage) (*entity.PublicMessage, error)
	GetAllPublicMessages(ctx context.Context, offset, limit int) []*entity.PublicMessage
	GetPublicMessage(ctx context.Context, id int) (*entity.PublicMessage, error)
}

type PrivateMessageRepo interface {
	AddPrivateMessage(ctx context.Context, msg entity.PrivateMessage) (*entity.PrivateMessage, error)
	GetAllPrivateMessages(ctx context.Context, offset, limit int) []*entity.PrivateMessage
	GetPrivateMessage(ctx context.Context, id int) (*entity.PrivateMessage, error)
}

func initDB(ctx context.Context) (*inmemory.InMemDB, <-chan any) {
	var inMemDB *inmemory.InMemDB

	var savedChan <-chan any

	dbStateRestored := true

	jsonDb, err := os.ReadFile(dbSavePath)
	if err != nil {
		dbStateRestored = false
	} else {
		inMemDB, savedChan, err = inmemory.NewInMemDBFromJSON(ctx, string(jsonDb), dbSavePath)
		if err != nil {
			dbStateRestored = false
		}
	}

	if !dbStateRestored {
		inMemDB, savedChan = inmemory.NewInMemDB(ctx, dbSavePath)
	}

	return inMemDB, savedChan
}

func initInMemRepos(ctx context.Context, conf *config.Config, savedChan *<-chan any) (*inmemoryrepository.UserRepo, *inmemoryrepository.PublicMessageRepo, *inmemoryrepository.PrivateMessageRepo) {
	db, ch := initDB(ctx)

	savedChan = &ch

	if conf.InMemoryDB.LoadFixtures {
		fixtures.LoadFixtures(db)
	}

	return inmemoryrepository.NewUserRepo(db), inmemoryrepository.NewPublicMessageRepo(db), inmemoryrepository.NewPrivateMessageRepo(db)
}

func initPostgresRepos(conf *config.Config, logger *logrus.Logger) (*postgresrepo.UserRepo, *postgresrepo.PublicMessageRepo, *postgresrepo.PrivateMessageRepo) {
	connStr := conf.Postgres.ConnectionURL()

	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		logger.Fatalf("cannot open database connection with connection string: %v, err: %v", connStr, err)
	}

	db := sqlx.NewDb(conn, "postgres")

	return postgresrepo.NewUserRepo(db), postgresrepo.NewPublicMessageRepo(db), postgresrepo.NewPrivateMessageRepo(db)
}

func initServices(
	usrRepo UserRepo,
	pubMsgRepo PublicMessageRepo,
	privMsgRepo PrivateMessageRepo) (
	*userservice.Service,
	*publicmessageservice.Service,
	*privatemessageservice.Service,
	*authservice.Service) {
	return userservice.New(usrRepo), publicmessageservice.New(pubMsgRepo, usrRepo), privatemessageservice.New(privMsgRepo, usrRepo), authservice.New(usrRepo)
}

func initConfig() (*config.Config, error) { // todo: to internals utils?
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var conf config.Config
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}

	// env variables
	if err := godotenv.Load(configPath + "/.env"); err != nil {
		return nil, err
	}

	viper.SetEnvPrefix("chat")
	viper.AutomaticEnv()

	// validate todo: VALIDATOR!

	conf.Jwt.Secret = viper.GetString("JWT_SECRET")
	if conf.Jwt.Secret == "" {
		return nil, errors.New("CHAT_JWT_SECRET env variable not set")
	}

	if !slices.Contains([]string{"jwt", "basic"}, strings.ToLower(conf.Auth)) {
		return nil, errors.New("invalid server.auth provided")
	}

	return &conf, nil
}

func initAuthMiddleware(typ string, secret string, authService authhandler.AuthService, logger *logrus.Logger, valid *validator.Validate) middlewares.Handler {
	switch typ {
	case "jwt":
		return middlewares.JWTAuthMiddleware(secret, logger)

	case "basic":
		return middlewares.BasicAuthMiddleware(authService, logger, valid)

	default: // jwt auth by default
		return middlewares.JWTAuthMiddleware(secret, logger)
	}
}

func main() {
	logger := logrus.New()
	ctx, cancel := context.WithCancel(context.Background())

	conf, err := initConfig()
	if err != nil {
		logger.Fatalf("init conf error: %v", err)
	}

	logger.Infof("CONFIG: %+v", conf)

	var (
		userRepo           UserRepo
		publicMessageRepo  PublicMessageRepo
		privateMessageRepo PrivateMessageRepo
	)

	var savedChan <-chan any

	switch conf.DB {
	case "postgres":
		userRepo, publicMessageRepo, privateMessageRepo = initPostgresRepos(conf, logger)

	case "inmem":
		userRepo, publicMessageRepo, privateMessageRepo = initInMemRepos(ctx, conf, &savedChan)

	default:
		userRepo, publicMessageRepo, privateMessageRepo = initPostgresRepos(conf, logger)
	}

	userService, publicMessageService, privateMessageService, authService := initServices(userRepo, publicMessageRepo, privateMessageRepo)

	valid := validator.New(validator.WithRequiredStructEnabled())

	authMiddleware := initAuthMiddleware(conf.Server.Auth, conf.Jwt.Secret, authService, logger, valid)
	loggingMiddleware := middlewares.LoggingMiddleware(logger, logrus.InfoLevel)
	recoveryMiddleware := middlewares.RecoveryMiddleware()

	authHandler := authhandler.New(userService, authService, conf.Jwt, logger, valid)
	userHandler := userhandler.New(userService, privateMessageService, logger, valid, authMiddleware)
	publicMessageHandler := publicmessagehandler.New(publicMessageService, userService, logger, valid, authMiddleware)
	privateMessageHandler := privatemessagehandler.New(privateMessageService, userService, logger, valid, authMiddleware)

	routers := make(map[string]chi.Router)

	routers["/auth"] = authHandler.Routes()
	routers["/users"] = userHandler.Routes()
	routers["/messages/public"] = publicMessageHandler.Routes()
	routers["/messages/private"] = privateMessageHandler.Routes()

	middlewars := []router.Middleware{
		recoveryMiddleware,
		loggingMiddleware,
	}

	r := router.MakeRoutes("/chat/api/v1", routers, middlewars)

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: r,
	}

	// add swagger middleware
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%v/swagger/doc.json", port)), // The url pointing to API definition
	))

	logger.Infof("server started at port %v", server.Addr)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Fatalf("server can't listen requests")
		}
	}()

	logger.Infof("documentation available on: http://localhost:%v/swagger/index.html", port)

	interrupt := make(chan os.Signal, 1)

	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(interrupt, syscall.SIGINT)

	go func() {
		<-interrupt

		logger.Info("interrupt signal caught")
		logger.Info("chat api server shutting down")

		if err := server.Shutdown(ctx); err != nil {
			logger.WithError(err).Fatalf("can't close server listening on '%s'", server.Addr)
		}

		cancel()
	}()

	// wait for inmem db being saved
	<-savedChan
}
