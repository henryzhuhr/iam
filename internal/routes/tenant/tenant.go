// Package tenant provides route registration for tenant management.
package tenant

import (
	"net/http"

	tenanthandler "iam/internal/handler/tenant"
	"iam/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

// Router registers tenant-related routes.
type Router struct {
	server  *rest.Server
	svcCtx  *svc.ServiceContext
	handler *tenanthandler.Handler
}

// NewRouter creates a new tenant Router.
func NewRouter(server *rest.Server, svcCtx *svc.ServiceContext) *Router {
	return &Router{
		server:  server,
		svcCtx:  svcCtx,
		handler: tenanthandler.NewHandler(svcCtx.TenantService),
	}
}

// Register registers all tenant routes.
func (r *Router) Register() {
	r.server.AddRoutes(
		[]rest.Route{
			{Method: http.MethodGet, Path: "/tenants", Handler: r.handler.List},
			{Method: http.MethodGet, Path: "/tenants/:id", Handler: r.handler.Get},
			{Method: http.MethodPost, Path: "/tenants", Handler: r.handler.Create},
			{Method: http.MethodPut, Path: "/tenants/:id", Handler: r.handler.Update},
			{Method: http.MethodDelete, Path: "/tenants/:id", Handler: r.handler.Delete},
			{Method: http.MethodPut, Path: "/tenants/:id/status", Handler: r.handler.UpdateStatus},
		},
		rest.WithPrefix("/api/v1"),
	)
}
