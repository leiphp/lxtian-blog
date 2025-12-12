package content

import (
	"context"
	"encoding/json"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
	"lxtian-blog/common/pkg/utils"
	"lxtian-blog/common/repository/web_repo"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type DocsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文档列表
func NewDocsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocsLogic {
	return &DocsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DocsLogic) Docs(req *types.DocsReq) (resp *types.DocsResp, err error) {
	// 分页兜底
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	repo := web_repo.NewTxyDocsRepository(l.svcCtx.DB)
	// 构建查询条件
	condition := make(map[string]interface{})
	if req.Status != nil {
		condition["status"] = *req.Status
	}
	if req.Cid > 0 {
		condition["category_id"] = req.Cid
	}
	result, total, err := repo.GetList(l.ctx, condition, req.Page, req.PageSize, "", req.Keywords, "title", "desc")
	if err != nil {
		return nil, err
	}

	// 转换为 map 便于前端消费
	list, err := utils.StructSliceToMapSliceUsingJSON(result)
	if err != nil {
		return nil, err
	}
	for k, item := range list {
		// tags 可能是 json 字符串
		switch v := item["tags"].(type) {
		case string:
			var arr []string
			if json.Unmarshal([]byte(v), &arr) == nil {
				item["tags"] = arr
			} else {
				item["tags"] = []string{}
			}
		case []byte: // 数据库可能返回 []byte
			var arr []string
			if json.Unmarshal(v, &arr) == nil {
				item["tags"] = arr
			} else {
				item["tags"] = []string{}
			}
		default:
			// 不是字符串，直接返回空数组
			item["tags"] = []string{}
		}
		if !strings.HasPrefix(item["cover"].(string), "http://") && !strings.HasPrefix(item["cover"].(string), "https://") {
			list[k]["cover"] = l.svcCtx.QiniuClient.PrivateURL(item["cover"].(string), 3600)
		}
	}
	utils.FormatTimeFields(list, "created_at", "updated_at")
	resp = &types.DocsResp{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		List:     list,
	}

	return
}
