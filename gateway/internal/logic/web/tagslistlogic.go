package web

import (
	"context"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagsListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 标签列表
func NewTagsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagsListLogic {
	return &TagsListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TagsListLogic) TagsList(req *types.TagsListReq) (resp *types.TagsListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
