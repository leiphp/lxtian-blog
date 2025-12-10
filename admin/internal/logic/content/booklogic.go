package content

import (
	"context"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/pkg/utils"
	"strings"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BooKLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 书单管理
func NewBooKLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BooKLogic {
	return &BooKLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BooKLogic) BooK() (resp *types.BookResp, err error) {
	resp = new(types.BookResp)
	var results []map[string]interface{}
	err = l.svcCtx.DB.
		Model(&mysql.TxyBook{}).
		Select("*").
		Find(&results).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("书单不存在")
		}
		return nil, err // 其他数据库错误
	}
	for k, book := range results {
		// tags 可能是 json 字符串
		switch v := book["tags"].(type) {
		case string:
			var arr []string
			if json.Unmarshal([]byte(v), &arr) == nil {
				book["tags"] = arr
			} else {
				book["tags"] = []string{}
			}
		case []byte: // 数据库可能返回 []byte
			var arr []string
			if json.Unmarshal(v, &arr) == nil {
				book["tags"] = arr
			} else {
				book["tags"] = []string{}
			}
		default:
			// 不是字符串，直接返回空数组
			book["tags"] = []string{}
		}
		if !strings.HasPrefix(book["cover"].(string), "http://") && !strings.HasPrefix(book["cover"].(string), "https://") {
			results[k]["cover"] = l.svcCtx.QiniuClient.PrivateURL(book["cover"].(string), 3600)
		}
	}

	utils.FormatTimeFields(results, "created_at", "updated_at")
	utils.FormatBoolFields(results, "status")
	resp.Data = results
	return
}
