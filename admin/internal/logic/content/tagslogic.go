package content

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
	"lxtian-blog/common/pkg/model/mysql"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 标签列表
func NewTagsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagsLogic {
	return &TagsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TagsLogic) Tags() (resp *types.TagsResp, err error) {
	resp = new(types.TagsResp)
	var results []map[string]interface{}
	err = l.svcCtx.DB.
		Model(&mysql.TxyTag{}).
		Select("*").
		Find(&results).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("数据不存在")
		}
		return nil, err // 其他数据库错误
	}
	resp.Data = results

	return
}
