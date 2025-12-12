package content

import (
	"context"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
	"lxtian-blog/common/repository/web_repo"

	"github.com/zeromicro/go-zero/core/logx"
)

type DocsDelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文档删除
func NewDocsDelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocsDelLogic {
	return &DocsDelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DocsDelLogic) DocsDel(req *types.DocsDelReq) (resp *types.DocsDelResp, err error) {
	repo := web_repo.NewTxyDocsRepository(l.svcCtx.DB)

	// 调用通用删除方法
	if err = repo.Delete(l.ctx, uint64(req.Id)); err != nil {
		l.Errorf("delete doc failed, id:%d, err:%v", req.Id, err)
		return nil, err
	}

	resp = &types.DocsDelResp{
		Data: map[string]interface{}{
			"status": "success",
		},
	}
	return
}
