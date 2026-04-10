// Package svc provides service context and dependency injection for the application.
package svc

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"fmt"
	"iam/infra/cache"
	"iam/infra/database"
	"iam/infra/queue"
	"iam/internal/config"
	"iam/internal/repository"
	tenantsvc "iam/internal/service/tenant"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// ServiceContext holds all application-level dependencies.
type ServiceContext struct {
	Config config.Config
	Logger logx.Logger

	// Infra
	DB            *sql.DB
	Redis         *redis.Redis
	KafkaProducer *queue.KafkaProducer

	// Repositories
	TenantRepo *repository.TenantRepository

	// Services
	TenantService *tenantsvc.Service

	// JWT (stub, to be implemented later)
	JWTKey    *rsa.PrivateKey
	JWTPubKey *rsa.PublicKey
}

// NewServiceContext creates and initializes all dependencies.
func NewServiceContext(c config.Config) (*ServiceContext, error) {
	ctx := context.Background()
	logger := logx.WithContext(ctx)

	// Initialize MySQL
	db, err := database.NewMySQL(c.DB)
	if err != nil {
		return nil, fmt.Errorf("init mysql: %w", err)
	}
	logger.Info("mysql connected")

	// Initialize Redis
	redisClient, err := cache.NewRedis(c.Redis)
	if err != nil {
		return nil, fmt.Errorf("init redis: %w", err)
	}
	logger.Info("redis connected")

	// Initialize Kafka (stub)
	kafkaProducer, err := queue.NewKafkaProducer(c.Kafka)
	if err != nil {
		return nil, fmt.Errorf("init kafka: %w", err)
	}
	logger.Info("kafka producer initialized (stub)")

	// Initialize Repositories
	tenantRepo := repository.NewTenantRepository(db)

	// Initialize Services
	tenantSvc := tenantsvc.NewService(tenantRepo)

	return &ServiceContext{
		Config:        c,
		Logger:        logger,
		DB:            db,
		Redis:         redisClient,
		KafkaProducer: kafkaProducer,
		TenantRepo:    tenantRepo,
		TenantService: tenantSvc,
	}, nil
}

// Close releases all resources.
func (sc *ServiceContext) Close() error {
	if sc.DB != nil {
		sc.DB.Close()
	}
	if sc.KafkaProducer != nil {
		sc.KafkaProducer.Close()
	}
	return nil
}
