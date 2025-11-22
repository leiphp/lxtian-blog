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

type DocsCategoriesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 教程分类
func NewDocsCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocsCategoriesLogic {
	return &DocsCategoriesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DocsCategoriesLogic) DocsCategories() (resp *types.DocsCategoriesResp, err error) {
	res, err := l.svcCtx.WebRpc.DocsCategories(l.ctx, &web.DocsCategoriesReq{})
	if err != nil {
		logc.Errorf(l.ctx, "DocsCategories error message: %s", err)
		return nil, err
	}
	var result []*types.CategoryItem
	if err := json.Unmarshal([]byte(res.List), &result); err != nil {
		return nil, err
	}
	resp = new(types.DocsCategoriesResp)
	resp.List = result
	return
}
