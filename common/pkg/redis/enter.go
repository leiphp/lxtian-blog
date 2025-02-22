package redis

import "fmt"

const KeyPrefix = "blog:"

const (
	API_CACHE_HASH_CATEGORY = 1 //全量分类
	API_CACHE_HASH_TAG      = 2 //全量标签
	API_CACHE_HASH_USER     = 3 //全量用户
	UserTokenString         = 4
)

var apiCacheKeys = map[int]string{
	API_CACHE_HASH_CATEGORY: "BlogDb:Category", //全量分类
	API_CACHE_HASH_TAG:      "BlogDb:Tag",      //全量标签
	API_CACHE_HASH_USER:     "BlogDb:User",     //全量用户
	UserTokenString:         "user:token",      //用户token
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
