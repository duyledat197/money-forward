package cmd

import (
	"context"
	"log/slog"
	"math/rand"
	"os"
	"time"

	l "log"

	"user-management/configs"
	deliveries "user-management/internal/deliveries/http"
	"user-management/internal/entities"
	"user-management/internal/services"
	"user-management/pkg/cache"
	"user-management/pkg/http_server"
	"user-management/pkg/id_utils"
	log "user-management/pkg/logger"
	"user-management/pkg/lru"
	"user-management/pkg/postgres_client"
	"user-management/pkg/processor"

	"github.com/lmittmann/tint"
)

var (
	cfgs           *configs.Config
	logger         log.Logger
	httpServer     *http_server.HttpServer
	postgresClient *postgres_client.PostgresClient

	userCache cache.Cache[int64, *entities.User]

	idGenerator id_utils.IDGenerator

	userService services.UserService

	processors []processor.Processor
	factories  []processor.Factory
)

func loadConfigs() {
	var err error
	cfgs, err = configs.LoadConfig("developments", "dev")
	if err != nil {
		l.Fatalln(err.Error())
	}
	l.Println(cfgs.PostgresDB.Address())
}

func loadLogger() {
	logger = slog.New(tint.NewHandler(os.Stdout, nil))
}

func loadIDGenerator() {
	idGenerator = id_utils.NewSnowFlake(rand.Int63n(10))
}

func loadHttpServer() {
	httpServer = http_server.NewHttpServer(cfgs.HTTP, logger)
}

func loadCaches() {
	userCache = lru.NewLRU[int64, *entities.User](128, 24*time.Hour)
}

func loadPostgresClient() {
	postgresClient = postgres_client.NewPostgresClient(cfgs.PostgresDB.Address())
}

func loadServices() {
	userService = services.NewUserService(
		postgresClient,
		idGenerator,
		userCache,
	)
}

func registerHandlers() {
	deliveries.RegisterUserDelivery(httpServer, userService)
}

func registerFactories() {
	factories = append(factories, postgresClient)
}

func registerProcessors() {
	processors = append(processors, httpServer)
}

func loadDefault() {
	// loader
	loadConfigs()
	loadLogger()
	loadIDGenerator()
	loadPostgresClient()
	loadCaches()
	loadServices()
	loadHttpServer()

	// register
	registerHandlers()
	registerFactories()
	registerProcessors()
}

func start(ctx context.Context, errChan chan error) {
	logger.Info("processors", len(processors))
	logger.Info("factories", len(factories))
	for _, f := range factories {
		if err := f.Connect(ctx); err != nil {
			errChan <- err
		}
	}

	for _, p := range processors {
		logger.Info("start")
		go func(pr processor.Processor) {
			if err := pr.Start(ctx); err != nil {
				errChan <- err
			}
		}(p)
	}
}

func stop(ctx context.Context) {
	for _, f := range factories {
		if err := f.Close(ctx); err != nil {
			logger.Error("unable to stop factory", "err", err)
		}
	}

	for _, p := range processors {
		if err := p.Stop(ctx); err != nil {
			logger.Error("unable to stop processor", "err", err)
		}
	}
}
