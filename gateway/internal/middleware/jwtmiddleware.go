package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/common/pkg/jwts"
	"lxtian-blog/common/restful/response"
	"net/http"
	"strings"
	"time"
)

type JwtMiddleware struct {
	accessSecret string
	accessExpire int64
}

func NewJwtMiddleware(accessSecret string, accessExpire int64) *JwtMiddleware {
	return &JwtMiddleware{
		accessSecret: accessSecret,
		accessExpire: accessExpire,
	}
}

func (m *JwtMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// JWTAuthMiddleware implementation
		w.Header().Set("Content-Type", "application/json")
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			logc.Errorf(r.Context(), "JwtMiddleware error: %s", "请求头中token为空")
			response.Response(r, w, nil, errors.New("请求头中token为空"))
			return
		}
		parts := strings.Split(authorization, " ")
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			logc.Errorf(r.Context(), "JwtMiddleware error: %s", "请求头中token格式有误")
			response.Response(r, w, nil, errors.New("请求头中token格式有误"))
			return
		}
		claims, err := jwts.ParseToken(parts[1], m.accessSecret, m.accessExpire)
		if err != nil {
			logc.Errorf(r.Context(), "JwtMiddleware error: %s", err)
			response.Response(r, w, nil, response.ErrTokenInvalid)
			return
		}
		isExpire := claims.ExpiresAt.Before(time.Now())
		if isExpire {
			token, _ := jwts.GenToken(jwts.JwtPayLoad{
				UserID:   claims.UserID,
				Username: claims.Username,
				Role:     1,
			}, m.accessSecret, m.accessExpire)
			w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
		}
		r = r.WithContext(context.WithValue(r.Context(), "user_id", claims.UserID))
		r = r.WithContext(context.WithValue(r.Context(), "username", claims.Username))
		next(w, r)
	}
}
