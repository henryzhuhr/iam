package health

import (
	"context"
	"iam/internal/dto/health"
	"iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type HealthService struct {
	Logger logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHealthService(ctx context.Context, svcCtx *svc.ServiceContext) *HealthService {
	return &HealthService{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HealthService) Health() (resp *health.Response, err error) {
	// todo: add your logic here and delete this line
	l.Logger.Infof("health: logic 调用成功")

	resp = &health.Response{
		Message: "ok",
	}

	return
}
