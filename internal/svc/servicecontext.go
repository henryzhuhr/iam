// Package svc provides service context and dependency injection for the application.
package svc

import (
	"context"
	"iam/internal/config"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	// 全局配置
	Config config.Config
	// 全局日志
	Logger logx.Logger

	// Infra 基础设施配置
	Infra Infra

	// Repository
	Repository Repository
}

// Repository 结构体，包含所有仓库接口
type Repository struct {
}

// Infra 结构体，包含所有基础设施连接
type Infra struct {
}

// NewServiceContext 创建全局服务上下文实例。
// 返回错误时，调用方应处理该错误（如记录日志并退出程序）
func NewServiceContext(c config.Config) (*ServiceContext, error) {
	ctx := context.Background()
	// 初始化日志
	logger := logx.WithContext(ctx)

	return &ServiceContext{
		Config: c,
		Logger: logger,
	}, nil
}

// Close 关闭所有资源连接
func (sc *ServiceContext) Close() error {

	return nil
}
