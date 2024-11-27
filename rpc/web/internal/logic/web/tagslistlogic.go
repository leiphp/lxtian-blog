package weblogic

import (
	"context"

	"lxtian-blog/rpc/web/internal/svc"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagsListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTagsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagsListLogic {
	return &TagsListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TagsListLogic) TagsList(in *web.TagsListReq) (*web.TagsListResp, error) {
	// todo: add your logic here and delete this line

	return &web.TagsListResp{}, nil
}
