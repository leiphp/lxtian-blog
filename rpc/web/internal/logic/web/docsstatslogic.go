package weblogic

import (
	"context"

	"lxtian-blog/rpc/web/internal/svc"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logx"
)

type DocsStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDocsStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocsStatsLogic {
	return &DocsStatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DocsStatsLogic) DocsStats(in *web.DocsStatsReq) (*web.DocsStatsResp, error) {
	// 在一个查询中统计文档数量、分类数量和总阅读量
	type StatsResult struct {
		TotalDocs       int64 `gorm:"column:total_docs"`
		TotalCategories int64 `gorm:"column:total_categories"`
		TotalViews      int64 `gorm:"column:total_views"`
		TotalLikes      int64 `gorm:"column:total_likes"`
	}

	var result StatsResult
	err := l.svcCtx.DB.
		Table("txy_docs").
		Select("COUNT(1) as total_docs, COUNT(DISTINCT category_id) as total_categories, COALESCE(SUM(`view`), 0) as total_views, COALESCE(SUM(`like`), 0) as total_likes").
		Where("deleted_at IS NULL").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return &web.DocsStatsResp{
		TotalDocs:       result.TotalDocs,
		TotalCategories: result.TotalCategories,
		TotalViews:      result.TotalViews,
		TotalLikes:      result.TotalLikes,
	}, nil
}
