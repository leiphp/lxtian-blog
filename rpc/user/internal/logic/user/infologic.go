package userlogic

import (
	"context"
	"encoding/json"
	"fmt"
	"lxtian-blog/common/pkg/initcache"
	"lxtian-blog/common/pkg/utils"
	"lxtian-blog/rpc/user/internal/svc"
	"lxtian-blog/rpc/user/model"
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
	// 初始化缓存，缓存过期时间为 60 * 24 分钟
	cache, err := initcache.InitCache(60*24, "user")
	if err != nil {
		return nil, err
	}
	// 尝试从缓存获取用户信息
	v, exist := cache.Get(fmt.Sprintf("userInfo:%d", in.Id))
	if exist {
		// 缓存中存在数据，直接返回
		return l.returnData(v)
	}
	// 如果缓存中没有，查询数据库
	txyUser, err := l.getUserFromDB(in.Id)
	if err != nil {
		return nil, err
	}
	// 转换为JSON小写标签格式
	res, err := utils.ConvertToLowercaseJSONTags(txyUser)
	if err != nil {
		return nil, err
	}
	// 设置缓存
	cache.Set(fmt.Sprintf("userInfo:%d", in.Id), res)
	// 返回结果
	return l.returnData(res)
}

// 从数据库获取用户信息
func (l *InfoLogic) getUserFromDB(id uint32) (*model.TxyUser, error) {
	var txyUser model.TxyUser
	if err := l.svcCtx.DB.First(&txyUser, "id = ?", id).Debug().Error; err != nil {
		return nil, err
	}
	return &txyUser, nil
}

// 返回InfoResp
func (l *InfoLogic) returnData(data interface{}) (*user.InfoResp, error) {
	// 转换缓存中的数据为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return &user.InfoResp{
		Data: string(jsonData),
	}, nil
}
