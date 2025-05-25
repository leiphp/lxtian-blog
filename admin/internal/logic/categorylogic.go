package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文章分类
func NewCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryLogic {
	return &CategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategoryLogic) Category() (resp *types.CategoryResp, err error) {
	resp = new(types.CategoryResp)
	var results []map[string]interface{}
	err = l.svcCtx.DB.
		Model(&mysql.TxyCategory{}).
		Select("*").
		Find(&results).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, err // 其他数据库错误
	}
	utils.FormatTimeFields(results, "created_at", "updated_at")
	utils.FormatBoolFields(results, "status")
	resp.Data = results
	return
}
