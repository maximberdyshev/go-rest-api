package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go-rest-api/config"
	"go-rest-api/internal/composite"
	http_v1_route "go-rest-api/internal/transport/http/v1/route"
	http_server "go-rest-api/pkg/http-server"
	"go-rest-api/pkg/logger"
	"go-rest-api/pkg/postgres"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

func Run(ctx context.Context) {
	logger := logger.FromContext(ctx)
	cfg := config.FromContext(ctx)

	pgClient, err := postgres.New((*postgres.Config)(&cfg.Postgres))
	if err != nil {
		logger.Fatal("Error initialize DB", zap.Error(err))
	}

	composite := composite.New(ctx, pgClient)

	router := httprouter.New()
	http_v1_route.SwaggerRouteRegister(ctx, router)
	http_v1_route.MusicRouteRegister(ctx, router, composite)

	server := http_server.New(router, http_server.Port(cfg.HTTP.Port))
	logger.Info("HTTP-server started")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("Received interrupt signal", zap.String("signal", s.String()))
	case err := <-server.Notify():
		logger.Error("HTTP-server received error", zap.Error(err))
	}

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("HTTP-server shutdown error", zap.Error(err))
	} else {
		logger.Info("HTTP-server stopped")
	}
}
