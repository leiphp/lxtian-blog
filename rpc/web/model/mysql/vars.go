package mysql

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

var ErrNotFound = sqlx.ErrNotFound

type JSONTime time.Time

// 实现 JSON 序列化接口
func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := "\"" + time.Time(t).Format("2006-01-02 15:04:05") + "\""
	return []byte(formatted), nil
}
