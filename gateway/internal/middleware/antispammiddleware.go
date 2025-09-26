package middleware

import (
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"lxtian-blog/common/pkg/security"
	"lxtian-blog/common/pkg/utils"
	"lxtian-blog/common/restful/response"
	"net/http"
)

type AntiSpamMiddleware struct {
	antiSpam *security.AntiSpam
}

func NewAntiSpamMiddleware(rds *redis.Redis) *AntiSpamMiddleware {
	return &AntiSpamMiddleware{
		antiSpam: security.NewAntiSpam(rds),
	}
}

func (m *AntiSpamMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := utils.GetClientIp(r)
		userAgent := r.Header.Get("User-Agent")
		endpoint := r.URL.Path

		// 检查是否为恶意请求
		isSpam, err := m.antiSpam.CheckSpam(r.Context(), clientIP, userAgent, endpoint)
		if err != nil {
			logc.Errorf(r.Context(), "反刷检查失败: %s", err)
			response.Response(r, w, nil, response.ErrServerError)
			return
		}

		if isSpam {
			response.Response(r, w, nil, &response.HttpError{
				Message:    "访问被拒绝",
				StatusCode: http.StatusForbidden,
			})
			return
		}
		next(w, r)
	}
}
