package svc

import (
	"lxtian-blog/common/pkg/initdb"
	"lxtian-blog/gateway/internal/config"
	"lxtian-blog/gateway/internal/middleware"
	"lxtian-blog/rpc/user/client/user"
	"lxtian-blog/rpc/web/client/web"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config             config.Config
	Rds                *redis.Redis
	WebRpc             web.Web
	UserRpc            user.User
	JwtMiddleware      rest.Middleware
	SecurityMiddleware *middleware.SecurityMiddleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds := initdb.InitRedis(c.RedisConfig.Host, c.RedisConfig.Type, c.RedisConfig.Pass, c.RedisConfig.Tls)
	securityMiddleware := middleware.NewSecurityMiddleware(rds)

	return &ServiceContext{
		Config:             c,
		Rds:                rds,
		WebRpc:             web.NewWeb(zrpc.MustNewClient(c.WebRpc)),
		UserRpc:            user.NewUser(zrpc.MustNewClient(c.UserRpc)),
		JwtMiddleware:      middleware.NewJwtMiddleware(c.Auth.AccessSecret, c.Auth.AccessExpire).Handle,
		SecurityMiddleware: securityMiddleware,
	}
}

var (
	// 自定义 QPS 计数器（示例：统计订单创建QPS）
	OrderCreateQPS = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_create_qps_total",
			Help: "Total number of order creation requests",
		},
		[]string{"method"}, // 标签维度（按方法名分组）
	)

	// 自定义 QPS 计数器（示例：统计文章浏览QPS）
	ArticleViewQPS = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "article_view_qps_total",
			Help: "Total number of article view requests",
		},
		[]string{"method"}, // 标签维度（按方法名分组）
	)
)

func init() {
	// 注册自定义指标到全局Prometheus注册表
	prometheus.MustRegister(OrderCreateQPS, ArticleViewQPS)
}
