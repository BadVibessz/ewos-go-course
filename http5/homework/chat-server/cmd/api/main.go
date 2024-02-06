package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/fixtures"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handler"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/router"

	inmemory "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/db/in-memory"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/docs"
)

//	@title			Chat API
//	@version		1.0
//	@description	API Server for Web Chat

//	@BasePath	/chat/api/v1

//	@securityDefinitions.basic	BasicAuth
//	@in							header
//	@name						Authorization

const ( // todo: config file
	dbSavePath = "http5/homework/chat-server/internal/db/db_state.json"
	port       = 5000
)

func initDB(ctx context.Context) (*inmemory.InMemDB, <-chan any) {
	var inMemDB *inmemory.InMemDB

	var savedChan <-chan any

	dbStateRestored := true

	jsonDb, err := os.ReadFile(dbSavePath)
	if err == nil {
		inMemDB, savedChan, err = inmemory.NewInMemDBFromJSON(ctx, string(jsonDb), dbSavePath)
		if err != nil {
			dbStateRestored = false
		}
	} else {
		dbStateRestored = false
	}

	if !dbStateRestored {
		inMemDB, savedChan = inmemory.NewInMemDB(ctx, dbSavePath)
	}

	return inMemDB, savedChan
}

func initInMemServices(db inmemory.InMemoryDB) (*service.UserService, *service.MessageService, *service.AuthBasicService) {
	userRepo := repository.NewInMemUserRepo(db)
	privateMsgRepo := repository.NewInMemPrivateMessageRepo(db)
	publicMsgRepo := repository.NewInMemPublicMessageRepo(db)

	userService := service.NewUserService(userRepo)
	messageService := service.NewMessageService(privateMsgRepo, publicMsgRepo, userRepo)
	authService := service.NewBasicAuthService(userRepo)

	return userService, messageService, authService
}

func main() {
	logger := logrus.New()

	ctx, cancel := context.WithCancel(context.Background())

	inMemDB, savedChan := initDB(ctx)
	userService, messageService, authService := initInMemServices(inMemDB)

	fixtures.LoadFixtures(ctx, userService, messageService) // todo: maybe load fixtures via db layer, not service?

	userHandler := handler.NewUserHandler(userService, messageService, authService, logger)
	messageHandler := handler.NewMessageHandler(messageService, userService, authService, logger)

	routers := make(map[string]chi.Router)

	routers["/users"] = userHandler.Routes()
	routers["/messages"] = messageHandler.Routes()

	middlewares := []router.Middleware{
		middleware.Recoverer,
	}

	r := router.MakeRoutes("/chat/api/v1", routers, middlewares)

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
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
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
