// Package health provides HTTP route registration for health-related endpoints.
package health

import (
	"net/http"

	healthHandler "iam/internal/handler/health"
	"iam/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

type healthRouter struct {
	server    *rest.Server
	serverCtx *svc.ServiceContext
}

func NewHealthRouter(server *rest.Server, serverCtx *svc.ServiceContext) *healthRouter {
	return &healthRouter{
		server:    server,
		serverCtx: serverCtx,
	}
}

// Register 注册健康检查相关路由
func (r *healthRouter) Register() {
	r.server.AddRoutes(
		[]rest.Route{
			{
				// 健康检查
				Method:  http.MethodGet,
				Path:    "/health",
				Handler: healthHandler.HealthHandler(r.serverCtx),
			},
		},
		rest.WithPrefix("/api"),
	)
}
