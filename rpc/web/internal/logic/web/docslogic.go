package weblogic

import (
	"context"
	"encoding/json"
	"lxtian-blog/common/repository/web_repo"
	"lxtian-blog/rpc/web/internal/svc"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DocsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDocsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocsLogic {
	return &DocsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DocsLogic) Docs(in *web.DocsReq) (*web.DocsResp, error) {
	// 创建文档repository
	docRepo := web_repo.NewTxyDocsRepository(l.svcCtx.DB)

	// 记录浏览次数（如果有IP参数）
	if in.ClientIp != "" {
		go func() {
			// 创建新的context，避免使用可能被取消的context
			ctx := context.Background()
			// 异步更新文档view
			if err := docRepo.IncrementDocView(ctx, int32(in.Id), in.ClientIp, l.svcCtx.Rds); err != nil {
				logc.Errorf(ctx, "记录文档浏览次数失败: %s", err)
			}
		}()
	}
	docID := int32(in.Id)
	// 1. 尝试从缓存获取文档详情
	doc, err := docRepo.GetDocDetail(l.ctx, docID, l.svcCtx.Rds)
	if err != nil {
		return nil, err
	}
	if doc != nil {
		logx.Infof("获取文档详情: %d", docID)
		// 将文档转换为 map，以便修改 tags 字段类型
		var docMap map[string]interface{}
		docBytes, err := json.Marshal(doc)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(docBytes, &docMap); err != nil {
			return nil, err
		}
		// 将 tags 字段从 JSON 字符串转换为数组
		if tagsStr, ok := docMap["tags"].(string); ok && tagsStr != "" {
			var tagsArray []string
			if err := json.Unmarshal([]byte(tagsStr), &tagsArray); err == nil {
				docMap["tags"] = tagsArray
			} else {
				// 如果解析失败，设置为空数组
				docMap["tags"] = []string{}
				logx.Errorf("解析 tags 失败: %v", err)
			}
		} else {
			// 如果 tags 为空或不存在，设置为空数组
			docMap["tags"] = []string{}
		}
		// 将文档数据转换为JSON字符串
		jsonData, err := json.Marshal(docMap)
		if err != nil {
			return nil, err
		}
		return &web.DocsResp{
			Data: string(jsonData),
		}, nil
	}

	return &web.DocsResp{
		Data: "",
	}, nil
}
