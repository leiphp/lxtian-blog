package userlogic

import (
	"context"
	"errors"
	"lxtian-blog/common/pkg/initcache"

	"lxtian-blog/rpc/user/internal/svc"
	"lxtian-blog/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	cache, err := initcache.InitCache(2, "user")
	if err != nil {
		return nil, err
	}
	v, exist := cache.Get("userInfo")
	if !exist {
		// deal with not exist
		return nil, errors.New("deal with not exist:数据不存在")
	}
	value, ok := v.(string)
	if !ok {
		// deal with type error
		return nil, errors.New("deal with type error:数据类型错误")
	}

	return &user.RegisterResp{
		Data: value,
	}, nil
}
