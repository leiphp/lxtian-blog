package middleware

import (
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"lxtian-blog/common/pkg/security"
	"lxtian-blog/common/pkg/utils"
	"lxtian-blog/common/restful/response"
	"net/http"
)

type RateLimitMiddleware struct {
	rateLimiter *security.RateLimiter
	config      security.RateLimitConfig
}

func NewRateLimitMiddleware(rds *redis.Redis) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		rateLimiter: security.NewRateLimiter(rds),
		config:      security.GetDefaultRateLimit(),
	}
}

func (m *RateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := utils.GetClientIp(r)
		endpoint := r.URL.Path

		// 检查限流
		allowed, err := m.rateLimiter.IsAllowed(r.Context(), clientIP, endpoint, m.config)
		if err != nil {
			logc.Errorf(r.Context(), "限流检查失败: %s", err)
			response.Response(r, w, nil, response.ErrServerError)
			return
		}

		if !allowed {
			response.Response(r, w, nil, &response.HttpError{
				Message:    "请求过于频繁，请稍后再试",
				StatusCode: http.StatusTooManyRequests,
			})
			return
		}

		// 添加限流信息到响应头
		remaining, _ := m.rateLimiter.GetRemainingRequests(r.Context(), clientIP, endpoint, m.config)
		w.Header().Set("X-RateLimit-Remaining", string(rune(remaining)))
		w.Header().Set("X-RateLimit-Limit", string(rune(m.config.MaxRequests)))

		next(w, r)
	}
}
