package weblogic

import (
	"context"
	"encoding/json"

	"lxtian-blog/rpc/web/internal/svc"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logx"
)

type DocsCategoriesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDocsCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocsCategoriesLogic {
	return &DocsCategoriesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DocsCategoriesLogic) DocsCategories(in *web.DocsCategoriesReq) (*web.DocsCategoriesResp, error) {
	var results []map[string]interface{}

	// 查询分类并统计每个分类下的文档数量
	// LEFT JOIN 文档表，统计未删除且状态为1的文档数量
	err := l.svcCtx.DB.
		Table("txy_docs_categories as c").
		Select("c.id, c.name, c.icon, c.status, c.color, c.sort, COUNT(d.id) as count").
		Joins("left join txy_docs as d on d.category_id = c.id AND d.deleted_at IS NULL AND d.status = 1").
		Where("c.deleted_at IS NULL AND c.status = 1").
		Group("c.id").
		Order("c.id ASC").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	// 转换为JSON字符串
	jsonData, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	return &web.DocsCategoriesResp{
		List: string(jsonData),
	}, nil
}
