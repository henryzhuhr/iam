package health

import (
	"net/http"

	"iam/internal/service/health"
	"iam/internal/svc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func HealthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 日志注入 user_id 字段
		ctx := logx.ContextWithFields(r.Context(), logx.Field("user_id", uuid.New().String()))
		l := health.NewHealthService(ctx, svcCtx)
		resp, err := l.Health()
		l.Logger.Infof("resp: %+v", resp)
		if err != nil {
			httpx.ErrorCtx(ctx, w, err)
		} else {
			httpx.OkJsonCtx(ctx, w, resp)
		}
	}
}
