package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sony/sonyflake"
)

var (
	sf      *sonyflake.Sonyflake
	chars   = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	numbers = []byte("0123456789")
)

func init() {
	sf = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: time.Date(2024, 5, 27, 0, 0, 0, 0, time.UTC),
	})
}

// GenerateOrderSN 生成订单号（优化版）
func GenerateOrderSN(prefix string) string {
	// 方案1: 使用雪花算法 + 时间戳
	timestamp := time.Now().Format("20060102150405") // YYYYMMDDHHMMSS
	snowflakeID := Snowflake()

	// 截取雪花算法ID的后8位
	shortID := snowflakeID % 100000000

	return fmt.Sprintf("%s%s%08d", prefix, timestamp, shortID)
}

// GenerateOrderSNWithLength 生成指定长度的订单号
func GenerateOrderSNWithLength(prefix string, length int) string {
	timestamp := time.Now().UnixNano() / 1e6
	snowflakeID := Snowflake()

	// 计算需要填充的长度
	baseLength := len(prefix) + len(strconv.FormatInt(timestamp, 10))
	remainingLength := length - baseLength

	if remainingLength <= 0 {
		remainingLength = 8 // 默认8位
	}

	// 截取雪花算法ID
	maxValue := int64(1) << (remainingLength * 4) // 16进制
	shortID := snowflakeID % uint64(maxValue)

	return fmt.Sprintf("%s%d%0*d", prefix, timestamp, remainingLength, shortID)
}

// Snowflake 雪花算法
func Snowflake() uint64 {
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}
	return id
}

// UUID uuid
func UUID() string {
	return uuid.New().String()
}

// GenerateMD5 生成MD5
func GenerateMD5(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	md5sum := hash.Sum(nil)
	return hex.EncodeToString(md5sum) // 将[]byte转为16进制字符串
}

// HashStringToRange 将字符串映射到 0-rang 范围内的数字
func HashStringToRange(s string, rang uint32) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	hashValue := h.Sum32()
	return int(hashValue % rang)
}

// GenerateRandomNumber 生成n位随机数
func GenerateRandomNumber(n int) string {
	buf := &strings.Builder{}
	l := len(chars)
	for i := 0; i < n; i++ {
		buf.WriteByte(chars[rand.Intn(l)])
	}
	return buf.String()
}

// GenerateRandomNumberOnly 生成n位随机数
func GenerateRandomNumberOnly(n int) string {
	buf := &strings.Builder{}
	length := len(numbers)
	for i := 0; i < n; i++ {
		buf.WriteByte(numbers[rand.Intn(length)])
	}
	return buf.String()
}

// GenerateRandomNumberNotFour 生成n位随机数 不包含4
func GenerateRandomNumberNotFour(n int) string {
	buf := &strings.Builder{}
	l := len(numbers)
	for i := 0; i < n; i++ {
		n := numbers[rand.Intn(l-1)]
		for n == '4' {
			n = numbers[rand.Intn(l-1)]
		}
		buf.WriteByte(n)
	}
	return buf.String()
}

// RandomNumber 得到范围内的随机数
func RandomNumber(min, max int) int {
	if min == max {
		return min
	}
	return rand.Intn(max-min) + min
}

// GenerateRandomString 随机生成16为字母数字符号大小写字符串
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"
	var result []byte
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charset))
		result = append(result, charset[randomIndex])
	}
	return string(result)
}
