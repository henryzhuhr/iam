package middleware

import (
	"context"
	"net/http"
	"strconv"

	"iam/internal/constant"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type tenantIDKey struct{}

// TenantMiddleware extracts tenant_id from request and injects it into context.
// Once AuthMiddleware is implemented, tenant_id should come from JWT claims.
func TenantMiddleware() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tenantIDStr := r.Header.Get("X-Tenant-ID")
			if tenantIDStr == "" {
				// AuthMiddleware not yet implemented, inject default value.
				// TODO: Extract tenant_id from JWT claims after AuthMiddleware is done.
				ctx := context.WithValue(r.Context(), tenantIDKey{}, int64(0))
				next(w, r.WithContext(ctx))
				return
			}

			tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
			if err != nil {
				logx.Errorf("invalid tenant_id: %s", tenantIDStr)
				httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, map[string]any{
					"code":    constant.CodeTenantNotFound,
					"message": "invalid tenant_id",
					"data":    nil,
				})
				return
			}

			ctx := context.WithValue(r.Context(), tenantIDKey{}, tenantID)
			next(w, r.WithContext(ctx))
		}
	}
}

// GetTenantID extracts tenant_id from context.
func GetTenantID(ctx context.Context) int64 {
	if v, ok := ctx.Value(tenantIDKey{}).(int64); ok {
		return v
	}
	return 0
}
