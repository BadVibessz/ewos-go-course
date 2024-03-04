package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/config"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository/in-memory"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

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

	inmemory "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/db/in-memory"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/docs"
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

func initInMemServices(db inmemory.InMemoryDB) (*userservice.Service, *publicmessageservice.Service, *privatemessageservice.Service, *authservice.Service) {
	userRepo := in_memory.NewInMemUserRepo(db)
	privateMsgRepo := in_memory.NewInMemPrivateMessageRepo(db)
	publicMsgRepo := in_memory.NewInMemPublicMessageRepo(db)

	userService := userservice.New(userRepo)
	publicMessageService := publicmessageservice.New(publicMsgRepo, userRepo)
	privateMessageService := privatemessageservice.New(privateMsgRepo, userRepo)
	authService := authservice.New(userRepo)

	return userService, publicMessageService, privateMessageService, authService
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

	conf, err := initConfig()
	if err != nil {
		logger.Fatalf("init conf error: %v", err)
	}

	logger.Infof("CONFIG: %+v", conf)

	ctx, cancel := context.WithCancel(context.Background())

	inMemDB, savedChan := initDB(ctx)
	if loadFixtures { // todo: to config
		fixtures.LoadFixtures(inMemDB)
	}

	userService, publicMessageService, privateMessageService, authService := initInMemServices(inMemDB)

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

	<-savedChan
}
