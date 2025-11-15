package user

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/rpc/user/user"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoLogic {
	return &InfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoLogic) Info() (resp *types.InfoResp, err error) {
	//从中间件获取用户信息
	userId, ok := l.ctx.Value("user_id").(uint)
	if !ok {
		logx.Errorf("Info user_id not found in context")
		return nil, errors.New("user_id not found in context")
	}
	username, ok := l.ctx.Value("username").(string)
	if !ok {
		return nil, errors.New("username not found in context")
	}
	logx.Infof("获取中间件用户信息: userId=%d, userName=%s", userId, username)
	//从user.rpc服务获取用户信息
	res, err := l.svcCtx.UserRpc.Info(l.ctx, &user.InfoReq{
		Id: uint32(userId),
	})
	if err != nil {
		logc.Errorf(l.ctx, "Info error: %s", err)
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(res.Data), &result); err != nil {
		return nil, err
	}
	return &types.InfoResp{
		Data: result,
	}, nil
}
