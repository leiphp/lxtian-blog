package redis

import (
	"fmt"
	"time"
)

const KeyPrefix = "blog:"

const (
	ApiWebStringCategory     = 1  //全量分类
	ApiWebStringTags         = 2  //全量tag
	ApiUserStringUser        = 3  //全量用户
	UserTokenString          = 4  //用户token
	UserScanString           = 5  //用户扫码
	WsUserIdString           = 6  //wx用户ID
	ApiWebStringColumn       = 7  //全量column
	ApiWebStringBook         = 8  //图书book
	ApiWebStringBookChapter  = 9  //图书章节
	ArticleViewString        = 10 //文章浏览次数记录
	OAuthStateString         = 11 //OAuth state
	UserMemberShipString     = 12 //用户会员
	DonatePendingOrderString = 13 //捐赠待支付订单
	DonatePendingOrderSet    = 14 //捐赠订单集合
	ApiUserInfoSet           = 15 //用户详情
)

var apiCacheKeys = map[int]string{
	ApiWebStringCategory:     "web:category",
	ApiWebStringTags:         "web:tags",
	ApiUserStringUser:        "user:list",
	UserTokenString:          "user:token",
	UserScanString:           "user:scan",
	WsUserIdString:           "user:ws",
	UserMemberShipString:     "user:membership",
	ApiWebStringColumn:       "web:column",
	ApiWebStringBook:         "web:book",
	ApiWebStringBookChapter:  "web:book:chapter",
	ArticleViewString:        "article:view",
	OAuthStateString:         "oauth:state",
	DonatePendingOrderString: "donate:order",
	DonatePendingOrderSet:    "donate:pending",
	ApiUserInfoSet:           "user:info:set",
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

/**
 * 获取文章浏览次数的Redis Key
 * 格式: blog:article:view:{article_id}:{ip}:{date}
 */
func GetArticleViewKey(articleID uint32, clientIP, date string) string {
	return fmt.Sprintf("%s%s:%d:%s:%s", KeyPrefix, apiCacheKeys[ArticleViewString], articleID, clientIP, date)
}

/**
 * 获取文章浏览次数的Redis Key（使用当前日期）
 * 格式: blog:article:view:{article_id}:{ip}:{today}
 */
func GetArticleViewKeyToday(articleID uint32, clientIP string) string {
	today := time.Now().Format("2006-01-02")
	return GetArticleViewKey(articleID, clientIP, today)
}
