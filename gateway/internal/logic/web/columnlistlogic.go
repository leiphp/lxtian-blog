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

type ColumnListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 专栏列表
func NewColumnListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ColumnListLogic {
	return &ColumnListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ColumnListLogic) ColumnList() (resp *types.ColumnListResp, err error) {
	res, err := l.svcCtx.WebRpc.ColumnList(l.ctx, &web.ColumnListReq{})
	if err != nil {
		logc.Errorf(l.ctx, "ColumnList error message: %s", err)
		return nil, err
	}
	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(res.List), &result); err != nil {
		return nil, err
	}
	resp = new(types.ColumnListResp)
	resp.Data = result
	return
}
