// Package cache provides Redis connection wrapper for the IAM application.
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

// RedisConfig holds Redis connection configuration.
type RedisConfig struct {
	Host     string `json:"Host"`
	Port     int    `json:"Port"`
	Password string `json:"Password,optional"`
	DB       int    `json:"DB,optional"`
}

// NewRedis creates a Redis client and verifies connectivity.
func NewRedis(cfg RedisConfig) (*redis.Redis, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	r := redis.New(addr, redis.WithPass(cfg.Password))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if !r.PingCtx(ctx) {
		return nil, fmt.Errorf("ping redis failed")
	}

	return r, nil
}
