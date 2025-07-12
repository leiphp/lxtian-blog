package weblogic

import (
	"context"
	"encoding/json"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/pkg/redis"
	"lxtian-blog/common/pkg/utils"
	"lxtian-blog/rpc/web/internal/svc"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logx"
)

type BookChapterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBookChapterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BookChapterLogic {
	return &BookChapterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BookChapterLogic) BookChapter(in *web.BookChapterReq) (*web.BookChapterResp, error) {
	// Redis key
	cacheKey := redis.ReturnRedisKey(redis.ApiWebStringBookChapter, in.Id)

	// 1. 查缓存
	cacheStr, err := l.svcCtx.Rds.Get(cacheKey)
	if err == nil && cacheStr != "" {
		// 缓存命中，直接返回
		return &web.BookChapterResp{
			Data: cacheStr,
		}, nil
	}

	var detail map[string]interface{}
	err = l.svcCtx.DB.
		Model(&mysql.TxyChapterData{}).
		Select("id,title,author,content,updated_at").
		Where("id =?", in.Id).
		Order("id desc").
		Debug().
		Find(&detail).Error
	if err != nil {
		return nil, err
	}
	utils.FormatTimeFieldsInMap(detail, "updated_at")
	jsonData, err := json.Marshal(detail)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.Rds.Set(redis.ReturnRedisKey(redis.ApiWebStringBookChapter, in.Id), string(jsonData))
	if err != nil {
		return nil, err
	}
	return &web.BookChapterResp{
		Data: string(jsonData),
	}, nil
}
