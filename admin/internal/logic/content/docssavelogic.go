package content

import (
	"context"
	"encoding/json"
	"lxtian-blog/common/model"
	"lxtian-blog/common/repository/web_repo"
	"time"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DocsSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文档保存
func NewDocsSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocsSaveLogic {
	return &DocsSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DocsSaveLogic) DocsSave(req *types.DocsSaveReq) (resp *types.DocsSaveResp, err error) {
	repo := web_repo.NewTxyDocsRepository(l.svcCtx.DB)

	// 序列化 tags 为 JSON 字符串
	tagsJSON, err := json.Marshal(req.Tags)
	if err != nil {
		l.Errorf("marshal tags failed, err:%v", err)
		return nil, err
	}

	// 准备数据
	now := time.Now()
	data := model.TxyDoc{
		CategoryID:  int32(req.CategoryId),
		Title:       req.Title,
		Description: req.Description,
		Content:     &req.Content,
		Level:       req.Level,
		Tags:        string(tagsJSON),
		Status:      req.Status,
		Cover:       req.Cover,
		View:        int32(req.View),
		UpdatedAt:   &now,
	}

	// 判断是新增还是更新
	if req.Id == 0 {
		// 新增
		data.CreatedAt = &now
		if err = repo.Create(l.ctx, &data); err != nil {
			return nil, err
		}
	} else {
		// 更新 - 使用 UpdateByCondition，只更新非零值字段
		condition := map[string]interface{}{
			"id": req.Id,
		}

		// 构建更新字段 map，排除零值字段
		updates := make(map[string]interface{})
		if data.CategoryID != 0 {
			updates["category_id"] = data.CategoryID
		}
		if data.Title != "" {
			updates["title"] = data.Title
		}
		if data.Description != "" {
			updates["description"] = data.Description
		}
		if data.Content != nil {
			updates["content"] = data.Content
		}
		if data.Level != "" {
			updates["level"] = data.Level
		}
		if data.Tags != "" {
			updates["tags"] = data.Tags
		}
		if data.Cover != "" {
			updates["cover"] = data.Cover
		}
		if data.Author != "" {
			updates["author"] = data.Author
		}
		// Status 是 bool 类型，需要显式设置
		updates["status"] = data.Status
		if data.View != 0 {
			updates["view"] = data.View
		}
		// UpdatedAt 总是更新
		updates["updated_at"] = data.UpdatedAt

		if err = repo.UpdateByCondition(l.ctx, condition, updates); err != nil {
			return nil, err
		}
	}

	resp = &types.DocsSaveResp{
		Data: true,
	}
	return
}
