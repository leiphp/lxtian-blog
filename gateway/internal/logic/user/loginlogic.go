package user

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/rpc/user/user"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	res, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		logc.Errorf(l.ctx, "Login error message: %s", err)
		return nil, err
	}
	fmt.Println("res:", res)
	resp = new(types.LoginResp)
	resp.Data = map[string]interface{}{
		"name": "leixiaotian",
	}
	return
}
