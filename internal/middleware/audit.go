package middleware

import (
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// AuditMiddleware logs request information for audit purposes.
// Will be connected to Kafka for async log publishing later.
func AuditMiddleware() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next(w, r)

			duration := time.Since(start)
			tenantID := GetTenantID(r.Context())

			logx.WithDuration(duration).Infof(
				"[audit] method=%s path=%s tenant_id=%d ip=%s duration=%s",
				r.Method, r.URL.Path, tenantID, r.RemoteAddr, duration,
			)

			// TODO: Send audit log via Kafka producer asynchronously.
		}
	}
}
