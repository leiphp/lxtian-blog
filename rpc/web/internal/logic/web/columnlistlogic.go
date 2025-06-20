package weblogic

import (
	"context"
	"encoding/json"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/pkg/redis"

	"lxtian-blog/rpc/web/internal/svc"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logx"
)

type ColumnListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewColumnListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ColumnListLogic {
	return &ColumnListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ColumnListLogic) ColumnList(in *web.ColumnListReq) (*web.ColumnListResp, error) {
	// Redis key
	cacheKey := redis.ReturnRedisKey(redis.ApiWebStringColumn, nil)

	// 1. 查缓存
	cacheStr, err := l.svcCtx.Rds.Get(cacheKey)
	if err == nil && cacheStr != "" {
		// 缓存命中，直接返回
		return &web.ColumnListResp{
			List: cacheStr,
		}, nil
	}

	where := map[string]interface{}{}
	var results []map[string]interface{}
	err = l.svcCtx.DB.
		Model(&mysql.TxyColumn{}).
		Select("id,name,slug,cover").
		Where(where).
		Order("id desc").
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	// 缓存tag
	err = l.svcCtx.Rds.Set(redis.ReturnRedisKey(redis.ApiWebStringColumn, nil), string(jsonData))
	if err != nil {
		return nil, err
	}
	return &web.ColumnListResp{
		List: string(jsonData),
	}, nil
}
