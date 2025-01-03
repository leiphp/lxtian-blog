package user

import (
	"context"
	"github.com/leiphp/unit-go-sdk/pkg/gconv"
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/rpc/user/user"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QrStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新扫码状态
func NewQrStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QrStatusLogic {
	return &QrStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QrStatusLogic) QrStatus(req *types.QrStatusReq) (resp *types.QrStatusResp, err error) {
	res, err := l.svcCtx.UserRpc.QrStatus(l.ctx, &user.QrStatusReq{
		Uuid:   req.Uuid,
		Status: req.Status,
	})
	if err != nil {
		logc.Errorf(l.ctx, "QrStatus error message: %s", err)
		return nil, err
	}
	return &types.QrStatusResp{
		Data: gconv.Uint32(res.Data),
	}, nil
}
