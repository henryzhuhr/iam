package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"example.com/basic-service/internal/config"
	"example.com/basic-service/internal/handler"
	"example.com/basic-service/internal/platform/httpserver"
	"example.com/basic-service/internal/platform/logging"
	"example.com/basic-service/internal/platform/postgres"
	redisclient "example.com/basic-service/internal/platform/redis"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := logging.New(cfg.AppEnv)

	pgPool, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		logger.Error("init postgres", slog.Any("err", err))
		os.Exit(1)
	}
	defer pgPool.Close()

	redisClient, err := redisclient.New(ctx, cfg.Redis)
	if err != nil {
		logger.Error("init redis", slog.Any("err", err))
		os.Exit(1)
	}
	defer redisClient.Close()

	mux := http.NewServeMux()
	mux.Handle("GET /healthz", handler.Health(logger, pgPool, redisClient))

	server := httpserver.New(cfg.HTTP, logger, mux)

	if err := server.Start(); err != nil {
		logger.Error("http server stopped", slog.Any("err", err))
		os.Exit(1)
	}

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown server", slog.Any("err", err))
	}
}
