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
		// 更新
		data.ID = int32(req.Id)
		if err = repo.Update(l.ctx, &data); err != nil {
			return nil, err
		}
	}

	resp = &types.DocsSaveResp{
		Data: true,
	}
	return
}
