package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"lxtian-blog/gateway/internal/config"
	"lxtian-blog/rpc/web/client/web"
)

type ServiceContext struct {
	Config config.Config
	WebRpc web.Web
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		WebRpc: web.NewWeb(zrpc.MustNewClient(c.WebRpc)),
	}
}
