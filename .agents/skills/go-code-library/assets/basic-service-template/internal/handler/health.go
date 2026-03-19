package handler

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
)

func Health(logger *slog.Logger, pg *pgxpool.Pool, redis *goredis.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := pg.Ping(ctx); err != nil {
			logger.Error("pg health check failed", slog.Any("err", err))
			http.Error(w, "postgres unavailable", http.StatusServiceUnavailable)
			return
		}

		if err := redis.Ping(ctx).Err(); err != nil {
			logger.Error("redis health check failed", slog.Any("err", err))
			http.Error(w, "redis unavailable", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
}
