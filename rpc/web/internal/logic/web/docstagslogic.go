package weblogic

import (
	"context"
	"encoding/json"

	"lxtian-blog/rpc/web/internal/svc"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logx"
)

type DocsTagsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDocsTagsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocsTagsLogic {
	return &DocsTagsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DocsTagsLogic) DocsTags(in *web.DocsTagsReq) (*web.DocsTagsResp, error) {
	// 查询所有文档的tags字段
	type DocTag struct {
		Tags string `gorm:"column:tags"`
	}
	var docs []DocTag
	err := l.svcCtx.DB.
		Table("txy_docs").
		Select("tags").
		Where("deleted_at IS NULL AND tags IS NOT NULL AND tags != '' AND tags != '[]'").
		Find(&docs).Error
	if err != nil {
		return nil, err
	}

	// 统计每个tag出现的次数
	tagCountMap := make(map[string]int64)
	for _, doc := range docs {
		var tags []string
		// 解析JSON格式的tags字段
		if err := json.Unmarshal([]byte(doc.Tags), &tags); err != nil {
			// 如果解析失败，跳过这条记录
			continue
		}
		// 遍历每个tag并计数
		for _, tag := range tags {
			if tag != "" {
				tagCountMap[tag]++
			}
		}
	}

	// 转换为包含tag和count的对象列表
	var results []map[string]interface{}
	for name, count := range tagCountMap {
		results = append(results, map[string]interface{}{
			"name":  name,
			"count": count,
		})
	}

	// 转换为JSON字符串
	jsonData, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	return &web.DocsTagsResp{
		List: string(jsonData),
	}, nil
}
