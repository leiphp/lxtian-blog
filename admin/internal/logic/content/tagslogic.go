package content

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

func (l *TagsLogic) Tags(req *types.TagsReq) (resp *types.TagsResp, err error) {
	resp = new(types.TagsResp)
	var results []map[string]interface{}

	// 默认分页参数兜底
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	db := l.svcCtx.DB.Model(&mysql.TxyTag{})
	if req.Keywords != "" {
		like := "%" + req.Keywords + "%"
		db = db.Where("name LIKE ? OR keywords LIKE ?", like, like)
	}

	// 统计总数
	var total int64
	if err = db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err = db.Order("id desc").
		Limit(req.PageSize).
		Offset(offset).
		Find(&results).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	utils.FormatTimeFields(results, "created_at", "updated_at")
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	resp.Total = total
	resp.List = results
	return
}
