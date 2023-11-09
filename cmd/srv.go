package cmd

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"user-management/configs"
	"user-management/internal/deliveries"
	"user-management/pkg/http_server"
	"user-management/pkg/postgres_client"
	"user-management/pkg/processor"
)

var (
	logger         *slog.Logger
	httpServer     *http_server.HttpServer
	postgresClient *postgres_client.PostgresClient

	userDelivery deliveries.UserDelivery

	processors []processor.Processor
	factories  []processor.Factory
)

func loadLogger() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func loadHttpServer() {
	httpServer = http_server.NewHttpServer(&configs.Endpoint{
		Port: "8080",
	}, logger)

}

func loadPostgresClient() {
	postgresClient = postgres_client.NewPostgresClient("")
}

func loadDeliveries() {
	userDelivery = deliveries.NewUserDelivery()
}

func registerHandlers() {
	http_server.Register(httpServer, http.MethodPost, "/users", userDelivery.CreateUser)
	http_server.Register(httpServer, http.MethodGet, "/users/{id}", userDelivery.GetUserByID)
}

func registerFactories() {
	factories = append(factories, postgresClient)
}

func registerProcessors() {
	processors = append(processors, httpServer)
}

func loadDefault() {
	loadLogger()
	loadPostgresClient()
	loadHttpServer()
	loadDeliveries()
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
