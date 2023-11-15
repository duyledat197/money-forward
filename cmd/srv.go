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
	"user-management/internal/repositories"
	"user-management/internal/services"
	"user-management/pkg/cache"
	"user-management/pkg/crypto_utils"
	"user-management/pkg/database"
	"user-management/pkg/http_server"
	"user-management/pkg/http_server/xcontext"
	"user-management/pkg/id_utils"
	log "user-management/pkg/logger"
	"user-management/pkg/lru"
	"user-management/pkg/postgres_client"
	"user-management/pkg/processor"
	"user-management/pkg/token_utils"

	"github.com/lmittmann/tint"
)

var (
	cfgs           *configs.Config
	logger         log.Logger
	httpServer     *http_server.HttpServer
	postgresClient *postgres_client.PostgresClient

	userByUserNameCache cache.Cache[string, *entities.User]
	userCache           cache.Cache[int64, *entities.UserWithAccounts]
	accountCache        cache.Cache[int64, *entities.Account]

	idGenerator    id_utils.IDGenerator
	tokenGenerator token_utils.Authenticator[*xcontext.UserInfo]

	userService    services.UserService
	authService    services.AuthService
	accountService services.AccountService

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

func loadGenerators() {
	var err error
	idGenerator = id_utils.NewSnowFlake(rand.Int63n(10))
	tokenGenerator, err = token_utils.NewPasetoAuthenticator[*xcontext.UserInfo](cfgs.SymetricKey)
	if err != nil {
		l.Fatalf("unable to create new token generator: %v", err)
	}
}

func loadHttpServer() {
	httpServer = http_server.NewHttpServer(
		cfgs.HTTP,
		logger,

		// middlewares will be handle by passing order.
		http_server.WithCors(), // using default allow access origin
		http_server.WithAuthenticate(tokenGenerator, []string{
			"POST /auth/login",
			"GET /users/{id}",
			"GET /users/{id}/accounts",
		}),
		http_server.WithRBAC(map[string][]entities.User_Role{
			"POST /users":   {entities.SuperAdminRole, entities.AdminRole},
			"PUT /users":    {entities.SuperAdminRole, entities.AdminRole, entities.UserRole},
			"DELETE /users": {entities.SuperAdminRole, entities.AdminRole},

			"POST /users/{id}/accounts": {entities.SuperAdminRole, entities.AdminRole, entities.UserRole},
			"PUT /accounts/{id}":        {entities.SuperAdminRole, entities.AdminRole, entities.UserRole},
			"DELETE /accounts/{id}":     {entities.SuperAdminRole, entities.AdminRole, entities.UserRole},
		}),
		http_server.WithRecovery(logger),
	)
}

func loadCaches() {
	userCache = lru.NewLRU[int64, *entities.UserWithAccounts](128, 24*time.Hour)
	accountCache = lru.NewLRU[int64, *entities.Account](128, 24*time.Hour)
	userByUserNameCache = lru.NewLRU[string, *entities.User](128, 24*time.Hour)
}

func loadPostgresClient() {
	postgresClient = postgres_client.NewPostgresClient(cfgs.PostgresDB.Address())
}

func loadServices() {
	userService = services.NewUserService(
		postgresClient,
		idGenerator,
		userCache,
		userByUserNameCache,
	)

	accountService = services.NewAccountService(postgresClient, accountCache)

	authService = services.NewAuthService(
		postgresClient,
		idGenerator,
		tokenGenerator,
		userByUserNameCache,
	)
}

func registerHandlers() {
	deliveries.RegisterUserDelivery(httpServer, userService)
	deliveries.RegisterAuthDelivery(httpServer, authService)
	deliveries.RegisterAccountDelivery(httpServer, accountService)
}

func registerFactories() {
	factories = append(factories, postgresClient)
}

func registerProcessors() {
	processors = append(processors, httpServer)
}

func migrateAdmin(ctx context.Context) {
	id := idGenerator.Int64()
	pwd, _ := crypto_utils.HashPassword(cfgs.SuperAdminPassword)
	userRepo := repositories.NewUserRepository()
	if err := userRepo.Upsert(ctx, postgresClient, &entities.User{
		ID:        id,
		Name:      database.NullString("admin"),
		UserName:  cfgs.SuperAdminUsername,
		Password:  pwd,
		Role:      entities.SuperAdminRole,
		CreatedBy: id,
	}); err != nil {
		l.Fatalf("unable to migrate super admin : %v", err)
	}
}

func loadDefault() {
	// loader
	loadConfigs()
	loadLogger()
	loadGenerators()
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
	for _, f := range factories {
		if err := f.Connect(ctx); err != nil {
			errChan <- err
		}
	}

	for _, p := range processors {
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
