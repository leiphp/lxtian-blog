package main

import (
	"flag"
	"fmt"
	"lxtian-blog/common/pkg/utils"
	"lxtian-blog/gateway/internal/utils/configcenter"
	"net/http"
	"os"

	"lxtian-blog/gateway/internal/config"
	"lxtian-blog/gateway/internal/handler"
	"lxtian-blog/gateway/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/gateway-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	// 使用通用方法解析Etcd主机列表字符串
	c.WebRpc.Etcd.Hosts = utils.ParseHosts(os.Getenv("ETCD_HOSTS"))
	// 配置中心加载数据
	configcenter.LoadConfigFromEtcd(&c)
	// 解决跨域
	//server := rest.MustNewServer(c.RestConf)
	server := rest.MustNewServer(c.RestConf, rest.WithCustomCors(nil, func(w http.ResponseWriter) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}, "*"))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
