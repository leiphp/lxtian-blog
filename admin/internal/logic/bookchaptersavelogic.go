package logic

import (
	"context"
	"database/sql"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/pkg/redis"
	"time"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BookChapterSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 章节保存
func NewBookChapterSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BookChapterSaveLogic {
	return &BookChapterSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BookChapterSaveLogic) BookChapterSave(req *types.BookChapterSaveReq) (resp *types.BookChapterSaveResp, err error) {
	// 1.判断是插入还是更新
	data := mysql.TxyChapterData{
		Id:      uint64(req.Id),
		Title:   req.Title,
		Author:  req.Author,
		Content: req.Content,
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	if req.Cate == "insert" {
		data.CreatedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		if err = l.svcCtx.DB.Create(&data).Error; err != nil {
			return nil, err
		}
	} else {
		if err = l.svcCtx.DB.Model(&data).Updates(data).Error; err != nil {
			return nil, err
		}
	}
	// 删除缓存
	cacheKey := redis.ReturnRedisKey(redis.ApiWebStringBookChapter, req.Id)
	if _, err = l.svcCtx.Rds.Del(cacheKey); err != nil {
		return nil, err
	}

	resp = new(types.BookChapterSaveResp)
	resp.Data = true
	return
}
