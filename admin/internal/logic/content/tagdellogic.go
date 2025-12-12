package content

import (
	"context"
	"lxtian-blog/common/repository/web_repo"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagDelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 标签删除
func NewTagDelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagDelLogic {
	return &TagDelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TagDelLogic) TagDel(req *types.TagDelReq) (resp *types.TagDelResp, err error) {
	repo := web_repo.NewTxyTagRepository(l.svcCtx.DB)

	// 调用通用删除方法
	if err = repo.Delete(l.ctx, uint64(req.Id)); err != nil {
		l.Errorf("delete tag failed, id:%d, err:%v", req.Id, err)
		return nil, err
	}

	resp = &types.TagDelResp{
		Data: map[string]interface{}{
			"status": "success",
		},
	}

	return
}
