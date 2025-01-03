package userlogic

import (
	"context"
	"errors"
	"lxtian-blog/rpc/user/internal/svc"
	"lxtian-blog/rpc/user/user"
	"strings"

	"github.com/leiphp/unit-go-sdk/pkg/gconv"
	"github.com/zeromicro/go-zero/core/logx"
)

type QrStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQrStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QrStatusLogic {
	return &QrStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QrStatusLogic) QrStatus(in *user.QrStatusReq) (*user.QrStatusResp, error) {
	//检查uuid是否有效
	res, err := l.svcCtx.Rds.Get(in.Uuid)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(res) == "" {
		return nil, errors.New("小程序码已经失效")
	}
	switch in.Status {
	case 2:
		//已扫码
		l.svcCtx.Rds.Setex(in.Uuid, `{"code":2}`, 5*60)
	case 3:
		//正在登录
		l.svcCtx.Rds.Setex(in.Uuid, `{"code":3}`, 5*60)
	case 4:
		//取消登录
		l.svcCtx.Rds.Setex(in.Uuid, `{"code":4}`, 5*60)
	default:
		return nil, errors.New("登录状态码错误")
	}
	// 发生websocket msg todo
	return &user.QrStatusResp{
		Data: gconv.String(in.Status),
	}, nil
}