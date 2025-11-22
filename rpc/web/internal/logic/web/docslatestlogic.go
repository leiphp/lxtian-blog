package weblogic

import (
	"context"
	"encoding/json"

	"lxtian-blog/rpc/web/internal/svc"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logx"
)

type DocsLatestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDocsLatestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocsLatestLogic {
	return &DocsLatestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DocsLatestLogic) DocsLatest(in *web.DocsLatestReq) (*web.DocsLatestResp, error) {
	// 处理limit参数
	limit := int(in.Limit)
	if limit <= 0 {
		limit = 5 // 默认值
	}

	// 查询热门文档：按view倒序取limit条
	var results []map[string]interface{}
	err := l.svcCtx.DB.
		Table("txy_docs as d").
		Select("d.id,d.category_id,d.title,d.description,d.level,d.cover,d.created_at,d.view,d.like,d.comment").
		Where("d.deleted_at IS NULL").
		Order("d.created_at DESC").
		Limit(limit).
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	// 转换为JSON字符串
	jsonData, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	// 统计总数（用于返回）
	var total int64
	err = l.svcCtx.DB.
		Table("txy_docs").
		Where("deleted_at IS NULL").
		Count(&total).Error
	if err != nil {
		return nil, err
	}

	return &web.DocsLatestResp{
		List: string(jsonData),
	}, nil
}
