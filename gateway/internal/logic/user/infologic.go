package user

import (
	"context"
	"errors"
	"lxtian-blog/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logc"

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

	// 检查返回结果
	if res == nil {
		logc.Errorf(l.ctx, "Info response is nil")
		return nil, errors.New("info response is nil")
	}
	if res.User == nil {
		logc.Errorf(l.ctx, "Info response User is nil")
		return nil, errors.New("info response user is nil")
	}

	// 映射用户信息
	infoResp := &types.InfoResp{
		Id:       int64(res.User.Id),
		Uid:      int64(res.User.Uid),
		Username: res.User.Username,
		Nickname: res.User.Nickname,
		Email:    res.User.Email,
		HeadImg:  res.User.HeadImg,
		Gold:     int(res.User.Gold),
		Score:    int(res.User.Score),
		Type:     int(res.User.Type),
		Status:   int(res.User.Status),
	}

	// 映射会员信息
	if res.Membership != nil {
		infoResp.Vip = &types.MemberShip{
			Is_valid:   res.Membership.IsValid,
			Is_active:  int(res.Membership.IsActive),
			Levle:      int(res.Membership.Level),
			Start_time: res.Membership.StartTime,
			End_time:   res.Membership.EndTime,
			Type_id:    int(res.Membership.TypeId),
		}
	}

	return infoResp, nil
}
