package web

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/rpc/web/web"

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

func (l *TagsListLogic) TagsList() (resp *types.TagsListResp, err error) {
	res, err := l.svcCtx.WebRpc.TagsList(l.ctx, &web.TagsListReq{})
	if err != nil {
		logc.Errorf(l.ctx, "TagsList error message: %s", err)
		return nil, err
	}
	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(res.List), &result); err != nil {
		return nil, err
	}
	resp = new(types.TagsListResp)
	resp.Data = result
	return
}
