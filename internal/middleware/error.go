package middleware

import (
	"encoding/json"
	"net/http"

	"iam/internal/constant"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// ErrorMiddleware recovers from panics and returns a uniform error response.
func ErrorMiddleware() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					httpx.WriteJsonCtx(r.Context(), w, http.StatusInternalServerError, map[string]any{
						"code":    constant.CodeInternalError,
						"message": "internal server error",
						"data":    nil,
					})
				}
			}()

			next(w, r)
		}
	}
}

// WriteSuccess returns a uniform success response.
func WriteSuccess(w http.ResponseWriter, data any) {
	httpx.WriteJson(w, http.StatusOK, map[string]any{
		"code":    constant.CodeOK,
		"message": "success",
		"data":    data,
	})
}

// WriteError returns a uniform error response.
func WriteError(w http.ResponseWriter, httpStatus int, code int, message string) {
	httpx.WriteJson(w, httpStatus, map[string]any{
		"code":    code,
		"message": message,
		"data":    nil,
	})
}

// ToResponse converts an entity to a generic map for JSON response.
func ToResponse(entity any) map[string]any {
	if entity == nil {
		return nil
	}
	data, _ := json.Marshal(entity)
	var result map[string]any
	_ = json.Unmarshal(data, &result)
	return result
}
