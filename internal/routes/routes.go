// Package routes provides HTTP route registration for the application.
package routes

import (
	"iam/internal/routes/health"
	"iam/internal/routes/tenant"
	"iam/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

// RegisterHandlers registers all HTTP routes.
func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// Register health check routes
	healthRouter := health.NewHealthRouter(server, serverCtx)
	healthRouter.Register()

	// Register tenant management routes
	tenantRouter := tenant.NewRouter(server, serverCtx)
	tenantRouter.Register()
}
