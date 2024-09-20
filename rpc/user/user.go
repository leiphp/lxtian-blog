package main

import (
	"flag"
	"fmt"
	"lxtian-blog/common/pkg/utils"
	"lxtian-blog/rpc/user/internal/utils/configcenter"
	"os"

	"lxtian-blog/rpc/user/internal/config"
	"lxtian-blog/rpc/user/internal/server/user"
	"lxtian-blog/rpc/user/internal/svc"
	"lxtian-blog/rpc/user/user"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	// 使用通用方法解析Etcd主机列表字符串
	c.Etcd.Hosts = utils.ParseHosts(os.Getenv("ETCD_HOSTS"))
	// 配置中心加载数据
	configcenter.LoadConfigFromEtcd(&c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
