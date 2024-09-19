package configcenter

import (
	"context"
	"github.com/zeromicro/go-zero/core/configcenter"
	"github.com/zeromicro/go-zero/core/configcenter/subscriber"
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/common/pkg/define"
	"lxtian-blog/gateway/internal/config"
)

// LoadConfigFromEtcd 从 etcd配置中心获取配置并更新
func LoadConfigFromEtcd(c *config.Config) {
	// 创建 etcd subscriber
	ss := subscriber.MustNewEtcdSubscriber(subscriber.EtcdConf{
		Hosts: c.WebRpc.Etcd.Hosts, // etcd 地址
		Key:   "/gateway",          // 配置key
	})

	// 创建 configurator
	cc := configurator.MustNewConfigCenter[define.GatewayOverrides](configurator.Config{
		Type: "yaml", // 配置值类型：json,yaml,toml
	}, ss)

	// 获取配置
	// 注意: 配置如果发生变更，调用的结果永远获取到最新的配置
	v, err := cc.GetConfig()
	if err != nil {
		logc.Errorf(context.Background(), "cc.GetConfig error: %s", err)
	}
	// 如果想监听配置变化，可以添加 listener
	cc.AddListener(func() {
		v, err := cc.GetConfig()
		if err != nil {
			logc.Errorf(context.Background(), "cc.AddListener c.GetConfig( error: %s", err)
		}
		//虽然可以监听到配置文件的变化，但是无法重新赋值到配置文件
		logc.Infof(context.Background(), "update config: %s", v)
	})
	c.Telemetry.Name = v.Telemetry.Name
	c.Telemetry.Endpoint = v.Telemetry.Endpoint
	c.Telemetry.Batcher = v.Telemetry.Batcher
	c.Telemetry.Sampler = v.Telemetry.Sampler
}
