package middleware

import (
	"net/http"
	"strings"

	"lxtian-blog/common/pkg/security"
	"lxtian-blog/common/restful/response"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

// SecurityMiddleware Gateway层安全中间件
type SecurityMiddleware struct {
	rateLimiter *security.RateLimiter
	antiSpam    *security.AntiSpam
}

// NewSecurityMiddleware 创建Gateway层安全中间件
func NewSecurityMiddleware(rds *redis.Redis) *SecurityMiddleware {
	return &SecurityMiddleware{
		rateLimiter: security.NewRateLimiter(rds),
		antiSpam:    security.NewAntiSpam(rds),
	}
}

// GetRateLimiter 获取限流器
func (sm *SecurityMiddleware) GetRateLimiter() *security.RateLimiter {
	return sm.rateLimiter
}

// GetAntiSpam 获取反刷检测器
func (sm *SecurityMiddleware) GetAntiSpam() *security.AntiSpam {
	return sm.antiSpam
}

// RateLimitMiddleware 限流中间件
func (sm *SecurityMiddleware) RateLimitMiddleware(config security.RateLimitConfig) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)
			endpoint := r.URL.Path

			// 检查限流
			allowed, err := sm.rateLimiter.IsAllowed(r.Context(), clientIP, endpoint, config)
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
			remaining, _ := sm.rateLimiter.GetRemainingRequests(r.Context(), clientIP, endpoint, config)
			w.Header().Set("X-RateLimit-Remaining", string(rune(remaining)))
			w.Header().Set("X-RateLimit-Limit", string(rune(config.MaxRequests)))

			next(w, r)
		}
	}
}

// AntiSpamMiddleware 反刷中间件
func (sm *SecurityMiddleware) AntiSpamMiddleware() rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)
			userAgent := r.Header.Get("User-Agent")
			endpoint := r.URL.Path

			// 检查是否为恶意请求
			isSpam, err := sm.antiSpam.CheckSpam(r.Context(), clientIP, userAgent, endpoint)
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
}

// IPWhitelistMiddleware IP白名单中间件
func (sm *SecurityMiddleware) IPWhitelistMiddleware(whitelist []string) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)

			// 检查IP是否在白名单中
			if !isIPInWhitelist(clientIP, whitelist) {
				logc.Errorf(r.Context(), "IP %s 不在白名单中，拒绝访问", clientIP)
				response.Response(r, w, nil, &response.HttpError{
					Message:    "访问被拒绝",
					StatusCode: http.StatusForbidden,
				})
				return
			}

			next(w, r)
		}
	}
}

// RefererCheckMiddleware Referer检查中间件
func (sm *SecurityMiddleware) RefererCheckMiddleware(allowedReferers []string) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			referer := r.Header.Get("Referer")

			// 如果配置了允许的Referer，则检查
			if len(allowedReferers) > 0 {
				if !isRefererAllowed(referer, allowedReferers) {
					logc.Errorf(r.Context(), "Referer %s 不在允许列表中", referer)
					response.Response(r, w, nil, &response.HttpError{
						Message:    "访问被拒绝",
						StatusCode: http.StatusForbidden,
					})
					return
				}
			}

			next(w, r)
		}
	}
}

// getClientIP 获取客户端真实IP
func getClientIP(r *http.Request) string {
	// 检查 X-Forwarded-For 头
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 检查 X-Real-IP 头
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// 使用 RemoteAddr
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}

	return ip
}

// isIPInWhitelist 检查IP是否在白名单中
func isIPInWhitelist(clientIP string, whitelist []string) bool {
	if len(whitelist) == 0 {
		return true // 如果没有配置白名单，则允许所有IP
	}

	for _, allowedIP := range whitelist {
		if clientIP == allowedIP {
			return true
		}
	}

	return false
}

// isRefererAllowed 检查Referer是否被允许
func isRefererAllowed(referer string, allowedReferers []string) bool {
	if referer == "" {
		return true // 如果没有Referer，则允许
	}

	for _, allowed := range allowedReferers {
		if strings.HasPrefix(referer, allowed) {
			return true
		}
	}

	return false
}
