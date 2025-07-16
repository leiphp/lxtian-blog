package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/pkg/utils"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ColumnListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 专栏列表
func NewColumnListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ColumnListLogic {
	return &ColumnListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ColumnListLogic) ColumnList() (resp *types.ColumnListResp, err error) {
	resp = new(types.ColumnListResp)
	var results []map[string]interface{}
	err = l.svcCtx.DB.
		Model(&mysql.TxyColumn{}).
		Select("*").
		Find(&results).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("专栏不存在")
		}
		return nil, err // 其他数据库错误
	}
	utils.FormatTimeFields(results, "created_at", "updated_at")
	utils.FormatBoolFields(results, "status")
	resp.Data = results

	return
}
