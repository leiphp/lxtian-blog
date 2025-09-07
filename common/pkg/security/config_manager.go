package security

import (
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
)

// ConfigManager 配置管理器
type ConfigManager struct {
	config *SecurityConfig
}

// SecurityConfig 安全配置结构
type SecurityConfig struct {
	RateLimit RateLimitConfigs `json:",optional"`
}

// RateLimitConfigs 限流配置集合
type RateLimitConfigs struct {
	Default  RateLimitConfig            `json:",optional"`
	Article  RateLimitConfig            `json:",optional"`
	Category RateLimitConfig            `json:",optional"`
	User     RateLimitConfig            `json:",optional"`
	Custom   map[string]RateLimitConfig `json:",optional"`
}

// NewConfigManager 创建配置管理器
func NewConfigManager(configFile string) (*ConfigManager, error) {
	cm := &ConfigManager{}

	// 设置默认配置
	cm.setDefaultConfig()

	// 如果配置文件存在，则加载
	if configFile != "" {
		err := conf.Load(configFile, &cm.config)
		if err != nil {
			return nil, fmt.Errorf("加载安全配置文件失败: %w", err)
		}
	}

	return cm, nil
}

// setDefaultConfig 设置默认配置
func (cm *ConfigManager) setDefaultConfig() {
	cm.config = &SecurityConfig{
		RateLimit: RateLimitConfigs{
			Default: RateLimitConfig{
				WindowSize:  time.Minute,
				MaxRequests: 60,
				KeyPrefix:   "rate_limit",
			},
			Article: RateLimitConfig{
				WindowSize:  time.Minute,
				MaxRequests: 30,
				KeyPrefix:   "article_rate",
			},
			Category: RateLimitConfig{
				WindowSize:  time.Minute,
				MaxRequests: 20,
				KeyPrefix:   "category_rate",
			},
			User: RateLimitConfig{
				WindowSize:  time.Minute,
				MaxRequests: 100,
				KeyPrefix:   "user_rate",
			},
			Custom: make(map[string]RateLimitConfig),
		},
	}
}

// GetRateLimitConfig 获取限流配置
func (cm *ConfigManager) GetRateLimitConfig(configType string) RateLimitConfig {
	switch configType {
	case "default":
		return cm.config.RateLimit.Default
	case "article":
		return cm.config.RateLimit.Article
	case "category":
		return cm.config.RateLimit.Category
	case "user":
		return cm.config.RateLimit.User
	default:
		// 检查自定义配置
		if custom, exists := cm.config.RateLimit.Custom[configType]; exists {
			return custom
		}
		// 返回默认配置
		return cm.config.RateLimit.Default
	}
}

// GetDefaultRateLimit 获取默认限流配置
func (cm *ConfigManager) GetDefaultRateLimit() RateLimitConfig {
	return cm.config.RateLimit.Default
}

// GetArticleRateLimit 获取文章限流配置
func (cm *ConfigManager) GetArticleRateLimit() RateLimitConfig {
	return cm.config.RateLimit.Article
}

// GetCategoryRateLimit 获取分类限流配置
func (cm *ConfigManager) GetCategoryRateLimit() RateLimitConfig {
	return cm.config.RateLimit.Category
}

// GetUserRateLimit 获取用户限流配置
func (cm *ConfigManager) GetUserRateLimit() RateLimitConfig {
	return cm.config.RateLimit.User
}

// AddCustomRateLimit 添加自定义限流配置
func (cm *ConfigManager) AddCustomRateLimit(name string, config RateLimitConfig) {
	if cm.config.RateLimit.Custom == nil {
		cm.config.RateLimit.Custom = make(map[string]RateLimitConfig)
	}
	cm.config.RateLimit.Custom[name] = config
}

// 全局配置管理器实例
var globalConfigManager *ConfigManager

// InitConfigManager 初始化全局配置管理器
func InitConfigManager(configFile string) error {
	var err error
	globalConfigManager, err = NewConfigManager(configFile)
	return err
}

// GetGlobalConfigManager 获取全局配置管理器
func GetGlobalConfigManager() *ConfigManager {
	if globalConfigManager == nil {
		// 如果没有初始化，使用默认配置
		globalConfigManager, _ = NewConfigManager("")
	}
	return globalConfigManager
}

// 为了向后兼容，提供全局函数
func GetDefaultRateLimit() RateLimitConfig {
	return GetGlobalConfigManager().GetDefaultRateLimit()
}

func GetArticleRateLimit() RateLimitConfig {
	return GetGlobalConfigManager().GetArticleRateLimit()
}

func GetCategoryRateLimit() RateLimitConfig {
	return GetGlobalConfigManager().GetCategoryRateLimit()
}

func GetUserRateLimit() RateLimitConfig {
	return GetGlobalConfigManager().GetUserRateLimit()
}
