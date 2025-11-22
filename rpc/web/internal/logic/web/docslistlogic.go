package weblogic

import (
	"context"
	"encoding/json"
	"strings"

	"lxtian-blog/rpc/web/internal/svc"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logx"
)

type DocsListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDocsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocsListLogic {
	return &DocsListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DocsListLogic) DocsList(in *web.DocsListReq) (*web.DocsListResp, error) {
	// 基础查询构建（包含JOIN分类表）
	baseDB := l.svcCtx.DB.
		Table("txy_docs as d").
		Joins("left join txy_docs_categories as c on d.category_id = c.id")

	// 填充WHERE条件
	if in.CategoryId > 0 {
		baseDB = baseDB.Where("d.category_id = ?", in.CategoryId)
	}
	if in.Level != "" {
		baseDB = baseDB.Where("d.level = ?", in.Level)
	}

	// 关键词搜索
	if in.Keywords != "" {
		baseDB = baseDB.Where("d.title like ? OR d.description like ?", "%"+in.Keywords+"%", "%"+in.Keywords+"%")
	}

	// tag搜索：在JSON数组字段中搜索包含指定tag的文档
	if in.Tag != "" {
		// 使用 JSON_SEARCH 在 tags JSON 数组中搜索指定的 tag
		// JSON_SEARCH 返回匹配的路径，如果找到则不为 NULL
		baseDB = baseDB.Where("JSON_SEARCH(d.tags, 'one', ?, NULL, '$[*]') IS NOT NULL", in.Tag)
	}

	// 计算总数（使用基础查询，无分页/排序）
	var total int64
	if err := baseDB.Count(&total).Error; err != nil {
		return nil, err
	}

	// 处理分页参数
	if in.Page == 0 {
		in.Page = 1
	}
	if in.PageSize == 0 {
		in.PageSize = 10
	}
	offset := (in.Page - 1) * in.PageSize

	// 处理排序
	sortBy := in.SortBy
	if sortBy == "" {
		sortBy = "id"
	}
	// 验证排序字段，防止SQL注入
	allowedSortFields := map[string]bool{
		"id":         true,
		"created_at": true,
		"view":       true,
		"like":       true,
		"comment":    true,
	}
	if !allowedSortFields[strings.ToLower(sortBy)] {
		sortBy = "id"
	}

	sortOrder := strings.ToUpper(in.SortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}
	orderClause := "d." + sortBy + " " + sortOrder

	// 查询数据
	var results []map[string]interface{}
	err := baseDB.Select("d.id,d.category_id,c.name as category_name,d.title,d.description,d.level,d.tags,d.cover,d.created_at,d.view,d.like,d.comment,d.reading_time").
		Limit(int(in.PageSize)).
		Offset(int(offset)).
		Order(orderClause).
		Debug().
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	// 将 tags 字段从 JSON 字符串转换为数组
	for i := range results {
		if tagsStr, ok := results[i]["tags"].(string); ok && tagsStr != "" {
			var tagsArray []string
			if err := json.Unmarshal([]byte(tagsStr), &tagsArray); err == nil {
				results[i]["tags"] = tagsArray
			} else {
				// 如果解析失败，设置为空数组
				results[i]["tags"] = []string{}
			}
		} else {
			// 如果 tags 为空或不是字符串，设置为空数组
			results[i]["tags"] = []string{}
		}
	}

	// 转换为JSON字符串
	jsonData, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	return &web.DocsListResp{
		Page:     in.Page,
		PageSize: in.PageSize,
		Total:    total,
		List:     string(jsonData),
	}, nil
}
