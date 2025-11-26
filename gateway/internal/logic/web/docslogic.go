package web

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/common/pkg/utils"
	"lxtian-blog/rpc/web/web"
	"net/http"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DocsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文档详情
func NewDocsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocsLogic {
	return &DocsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DocsLogic) Docs(req *types.DocsReq, r *http.Request) (resp *types.DocsResp, err error) {
	// 获取客户端IP
	clientIP := ""
	clientIP = utils.GetClientIP(r)
	logc.Infof(l.ctx, "文档 %d 被IP %s 访问", req.Id, clientIP)

	res, err := l.svcCtx.WebRpc.Docs(l.ctx, &web.DocsReq{
		Id:       req.Id,
		ClientIp: clientIP,
	})
	if err != nil {
		logc.Errorf(l.ctx, "Docs error: %s", err)
		return nil, err
	}
	resp = new(types.DocsResp)
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(res.Data), &result); err != nil {
		return nil, err
	}
	resp.Data = result
	return
}
