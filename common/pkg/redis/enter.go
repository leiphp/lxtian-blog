package redis

import "fmt"

const KeyPrefix = "blog:"

const (
	ApiWebStringCategory = 1 //全量分类
	ApiWebStringTags     = 2 //全量tag
	ApiUserStringUser    = 3 //全量用户
	UserTokenString      = 4 //用户token
	UserScanString       = 5 //用户扫码
	WsUserIdString       = 6 //wx用户ID
)

var apiCacheKeys = map[int]string{
	ApiWebStringCategory: "web:category",
	ApiWebStringTags:     "web:tags",
	ApiUserStringUser:    "user:list",
	UserTokenString:      "user:token",
	UserScanString:       "user:scan",
	WsUserIdString:       "user:ws",
}

/**
 * 获取对应的key
 */
func ReturnRedisKey(keyType int, key interface{}) string {

	var suffix string
	if key != nil {
		suffix = ":" + fmt.Sprintf("%v", key)
	}
	redisKey := KeyPrefix + apiCacheKeys[keyType] + suffix
	return redisKey
}
