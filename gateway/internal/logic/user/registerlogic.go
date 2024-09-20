package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/rpc/user/user"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	res, err := l.svcCtx.UserRpc.Register(l.ctx, &user.RegisterReq{
		Username: req.Username,
	})
	fmt.Println("Res:", res)
	fmt.Println("err222:", err)
	if err != nil {
		logc.Errorf(l.ctx, "Register error message: %s", err)
		return nil, err
	}
	resp = new(types.RegisterResp)
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(res.Data), &result); err != nil {
		return nil, err
	}
	resp.Data = result
	return
}
