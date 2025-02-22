package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"lxtian-blog/common/pkg/initdb"
	"lxtian-blog/gateway/internal/config"
	"lxtian-blog/gateway/internal/middleware"
	"lxtian-blog/rpc/user/client/user"
	"lxtian-blog/rpc/web/client/web"
)

type ServiceContext struct {
	Config        config.Config
	Rds           *redis.Redis
	WebRpc        web.Web
	UserRpc       user.User
	JwtMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds := initdb.InitRedis(c.RedisConfig.Host, c.RedisConfig.Type, c.RedisConfig.Pass, c.RedisConfig.Tls)
	return &ServiceContext{
		Config:        c,
		Rds:           rds,
		WebRpc:        web.NewWeb(zrpc.MustNewClient(c.WebRpc)),
		UserRpc:       user.NewUser(zrpc.MustNewClient(c.UserRpc)),
		JwtMiddleware: middleware.NewJwtMiddleware(c.Auth.AccessSecret, c.Auth.AccessExpire).Handle,
	}
}
