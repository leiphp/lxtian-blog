package user

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/common/pkg/jwts"
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
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(res.Data), &result); err != nil {
		return nil, err
	}
	// 获取token
	auth := l.svcCtx.Config.Auth
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		UserID:   uint(result["id"].(float64)),
		Username: result["username"].(string),
	}, auth.AccessSecret, auth.AccessExpire)
	if err != nil {
		return nil, err
	}
	resp = new(types.LoginResp)
	resp.User = result
	resp.AccessToken = token
	resp.ExpiresIn = uint64(auth.AccessExpire) * 3600
	return
}
