package alipay

import (
	"encoding/json"
	"fmt"
	"sync"
)

// AlipayConfig 支付宝配置
type AlipayConfig struct {
	AppId           string `json:"app_id"`            // 应用ID
	AppPrivateKey   string `json:"app_private_key"`   // 应用私钥
	AlipayPublicKey string `json:"alipay_public_key"` // 支付宝公钥
	GatewayUrl      string `json:"gateway_url"`       // 支付宝网关地址
	NotifyUrl       string `json:"notify_url"`        // 异步通知地址
	ReturnUrl       string `json:"return_url"`        // 同步跳转地址
	IsProd          bool   `json:"is_prod"`           // 是否生产环境
	SignType        string `json:"sign_type"`         // 签名类型，默认RSA2
	Charset         string `json:"charset"`           // 字符集，默认utf-8
	Format          string `json:"format"`            // 数据格式，默认JSON
	Version         string `json:"version"`           // 接口版本，默认1.0
	Timeout         string `json:"timeout"`           // 订单超时时间
}

// DefaultConfig 默认配置
func DefaultConfig() *AlipayConfig {
	return &AlipayConfig{
		SignType: "RSA2",
		Charset:  "utf-8",
		Format:   "JSON",
		Version:  "1.0",
		Timeout:  "30m",
	}
}

// SandboxConfig 沙箱环境配置
func SandboxConfig() *AlipayConfig {
	config := DefaultConfig()
	config.IsProd = false
	config.GatewayUrl = "https://openapi.alipaydev.com/gateway.do"
	return config
}

// ProductionConfig 生产环境配置
func ProductionConfig() *AlipayConfig {
	config := DefaultConfig()
	config.IsProd = true
	config.GatewayUrl = "https://openapi.alipay.com/gateway.do"
	return config
}

// ConfigManager 配置管理器
type ConfigManager struct {
	configs map[string]*AlipayConfig
	mu      sync.RWMutex
}

// NewConfigManager 创建配置管理器
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		configs: make(map[string]*AlipayConfig),
	}
}

// SetConfig 设置配置
func (cm *ConfigManager) SetConfig(key string, config *AlipayConfig) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.configs[key] = config
}

// GetConfig 获取配置
func (cm *ConfigManager) GetConfig(key string) (*AlipayConfig, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	config, exists := cm.configs[key]
	if !exists {
		return nil, fmt.Errorf("config not found for key: %s", key)
	}
	return config, nil
}

// GetDefaultConfig 获取默认配置
func (cm *ConfigManager) GetDefaultConfig() (*AlipayConfig, error) {
	return cm.GetConfig("default")
}

// LoadConfigFromJSON 从JSON字符串加载配置
func LoadConfigFromJSON(jsonStr string) (*AlipayConfig, error) {
	var config AlipayConfig
	if err := json.Unmarshal([]byte(jsonStr), &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 设置默认值
	if config.SignType == "" {
		config.SignType = "RSA2"
	}
	if config.Charset == "" {
		config.Charset = "utf-8"
	}
	if config.Format == "" {
		config.Format = "JSON"
	}
	if config.Version == "" {
		config.Version = "1.0"
	}
	if config.Timeout == "" {
		config.Timeout = "30m"
	}
	if config.GatewayUrl == "" {
		if config.IsProd {
			config.GatewayUrl = "https://openapi.alipay.com/gateway.do"
		} else {
			config.GatewayUrl = "https://openapi.alipaydev.com/gateway.do"
		}
	}

	return &config, nil
}

// Validate 验证配置
func (c *AlipayConfig) Validate() error {
	if c.AppId == "" {
		return fmt.Errorf("app_id is required")
	}
	if c.AppPrivateKey == "" {
		return fmt.Errorf("app_private_key is required")
	}
	if c.AlipayPublicKey == "" {
		return fmt.Errorf("alipay_public_key is required")
	}
	if c.GatewayUrl == "" {
		return fmt.Errorf("gateway_url is required")
	}
	return nil
}

// String 返回配置的字符串表示
func (c *AlipayConfig) String() string {
	return fmt.Sprintf("AlipayConfig{AppId:%s, IsProd:%v, GatewayUrl:%s}",
		c.AppId, c.IsProd, c.GatewayUrl)
}
