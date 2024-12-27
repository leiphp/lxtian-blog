package user

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/rpc/user/user"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetqrLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取二维码
func NewGetqrLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetqrLogic {
	return &GetqrLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetqrLogic) Getqr() (resp *types.GetqrResp, err error) {
	res, err := l.svcCtx.UserRpc.Getqr(l.ctx, &user.GetqrReq{})
	if err != nil {
		logc.Errorf(l.ctx, "Getqr error message: %s", err)
		return nil, err
	}
	resp = new(types.GetqrResp)
	resp.Uuid = res.Uuid
	resp.QrImg = res.QrImg
	return
}
