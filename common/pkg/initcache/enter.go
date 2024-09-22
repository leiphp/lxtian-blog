package initcache

import (
	"context"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/logc"
	"time"
)

// InitCache 初始化缓存
func InitCache(num int64, name string) (*collection.Cache, error) {
	if num == 0 {
		num = 10
	}
	expireTime := time.Minute * time.Duration(num)
	cache, err := collection.NewCache(expireTime, collection.WithName(name))
	if err != nil {
		logc.Errorf(context.Background(), "InitCache error message: %s", err)
	}
	return cache, err
}
