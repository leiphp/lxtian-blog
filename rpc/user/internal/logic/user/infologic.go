package userlogic

import (
	"context"
	"errors"
	"lxtian-blog/common/pkg/initcache"
	"lxtian-blog/rpc/user/internal/svc"
	"lxtian-blog/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoLogic {
	return &InfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InfoLogic) Info(in *user.InfoReq) (*user.InfoResp, error) {
	cache, err := initcache.InitCache(2, "user")
	if err != nil {
		return nil, err
	}
	cache.Set("userInfo", map[string]interface{}{
		"id":   in.Id,
		"name": "雷小天",
		"age":  30,
		"sex":  1,
	})
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
	// use value
	return &user.InfoResp{
		Data: value,
	}, nil
}
