package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"
	"lxtian-blog/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改用户信息
func NewUpdateInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateInfoLogic {
	return &UpdateInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateInfoLogic) UpdateInfo(req *types.UpdateInfoReq) (resp *types.UpdateInfoResp, err error) {
	//从中间件获取用户信息
	userId, ok := l.ctx.Value("user_id").(uint)
	if !ok {
		return nil, errors.New("user_id not found in context")
	}
	username, ok := l.ctx.Value("username").(string)
	if !ok {
		return nil, errors.New("username not found in context")
	}
	fmt.Println("userId:", userId)
	fmt.Println("userName:", username)

	if req.HeadImg == "" && req.Nickname == "" {
		return nil, errors.New("head_img, nickname is empty")
	}
	res, err := l.svcCtx.UserRpc.UpdateInfo(l.ctx, &user.UpdateInfoReq{
		Nickname: req.Nickname,
		HeadImg:  req.HeadImg,
		Id:       uint32(userId),
	})
	if err != nil {
		logc.Errorf(l.ctx, "UpdateInfo error: %s", err)
		return nil, err
	}

	resp = new(types.UpdateInfoResp)
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(res.Data), &result); err != nil {
		return nil, err
	}
	resp.Data = result
	return
}
