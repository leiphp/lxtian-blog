package content

import (
	"context"
	"database/sql"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/repository/web_repo"
	"time"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 标签保存
func NewTagSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagSaveLogic {
	return &TagSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TagSaveLogic) TagSave(req *types.TagSaveReq) (resp *types.TagSaveResp, err error) {
	repo := web_repo.NewTxyTagRepository(l.svcCtx.DB)
	// 准备数据
	data := mysql.TxyTag{
		Id:          uint64(req.Id),
		Name:        req.Name,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		Keywords:    sql.NullString{String: req.Keywords, Valid: req.Keywords != ""},
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	// 判断是新增还是更新
	if req.Id == 0 {
		// 新增
		data.CreatedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		if err = repo.Create(l.ctx, &data); err != nil {
			return nil, err
		}
	} else {
		// 更新
		data.Id = uint64(req.Id)
		if err = repo.Update(l.ctx, &data); err != nil {
			return nil, err
		}
	}

	resp = &types.TagSaveResp{
		Data: true,
	}
	return
}
