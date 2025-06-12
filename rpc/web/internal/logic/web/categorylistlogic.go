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

type CategoryListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryListLogic {
	return &CategoryListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CategoryListLogic) CategoryList(in *web.CategoryListReq) (*web.CategoryListResp, error) {
	// Redis key
	cacheKey := redis.ReturnRedisKey(redis.ApiWebStringCategory, nil)

	// 1. 查缓存
	cacheStr, err := l.svcCtx.Rds.Get(cacheKey)
	if err == nil && cacheStr != "" {
		// 缓存命中，直接返回
		return &web.CategoryListResp{
			List: cacheStr,
		}, nil
	}
	where := map[string]interface{}{}
	where["status"] = 1
	if in.Page == 0 {
		in.Page = 1
	}
	if in.PageSize == 0 {
		in.PageSize = 10
	}
	offset := (in.Page - 1) * in.PageSize
	var results []map[string]interface{}
	err = l.svcCtx.DB.
		Model(&mysql.TxyCategory{}).
		Select("id,name,seoname,description,keywords,sort,status").
		Where(where).
		Limit(int(in.PageSize)).
		Offset(int(offset)).
		Order("sort asc").
		Debug().
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}
	//计算当前type的总数，给分页算总页
	var total int64
	err = l.svcCtx.DB.Model(&mysql.TxyCategory{}).Where(where).Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 缓存分类
	err = l.svcCtx.Rds.Set(redis.ReturnRedisKey(redis.ApiWebStringCategory, nil), string(jsonData))
	if err != nil {
		return nil, err
	}
	return &web.CategoryListResp{
		Page:     in.Page,
		PageSize: in.PageSize,
		Total:    uint32(total),
		List:     string(jsonData),
	}, nil
}
