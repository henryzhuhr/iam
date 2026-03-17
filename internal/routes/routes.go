// routes 路由注册包
package routes

import (
	"iam/internal/routes/health"
	"iam/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {

	// 注册健康检查相关路由
	healthRouter := health.NewHealthRouter(server, serverCtx)
	healthRouter.Register()
}
