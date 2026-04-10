package middleware

import (
	"net/http"
	"strings"

	"iam/internal/constant"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// AuthMiddleware validates Bearer token in Authorization header.
// Currently a stub — full JWT validation to be implemented later.
func AuthMiddleware() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httpx.WriteJsonCtx(r.Context(), w, http.StatusUnauthorized, map[string]any{
					"code":    constant.CodeTokenInvalid,
					"message": "missing Authorization header",
					"data":    nil,
				})
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == authHeader {
				httpx.WriteJsonCtx(r.Context(), w, http.StatusUnauthorized, map[string]any{
					"code":    constant.CodeTokenInvalid,
					"message": "invalid Authorization format, expected 'Bearer <token>'",
					"data":    nil,
				})
				return
			}

			// TODO: Full JWT validation (parse, verify signature, check expiration, extract claims).
			logx.Debugf("auth middleware: token format ok (stub validation)")
			next(w, r)
		}
	}
}

// SkipAuthPaths lists paths that do not require authentication.
var SkipAuthPaths = map[string]bool{
	"/api/v1/auth/login":          true,
	"/api/v1/auth/register":       true,
	"/api/v1/auth/password/reset": true,
	"/api/v1/auth/code/send":      true,
	"/api/v1/auth/code/login":     true,
	"/health":                     true,
	"/api/health":                 true,
	"/api/v1/clients/token":       true,
}
