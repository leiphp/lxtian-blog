package logic

import (
	"context"
	"database/sql"
	"lxtian-blog/common/pkg/model/mysql"
	"time"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文章保存
func NewArticleSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleSaveLogic {
	return &ArticleSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleSaveLogic) ArticleSave(req *types.ArticleSaveReq) (resp *types.ArticleSaveResp, err error) {
	db := l.svcCtx.DB

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1.判断是插入还是更新
	data := mysql.TxyArticle{
		Title:       req.Title,
		Author:      req.Author,
		Content:     req.Content,
		Keywords:    req.Keywords,
		Path:        req.Path,
		Description: req.Description,
		Status:      req.Status,
		IsHot:       req.IsHot,
		IsRec:       req.IsRec,
		IsTop:       uint64(req.IsTop),
		IsOriginal:  uint64(req.IsOriginal),
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	if req.Id == 0 {
		data.CreatedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		if err = tx.Create(&data).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

	} else {
		data.Id = uint64(req.Id)
		if err = tx.Model(&data).Updates(data).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	// 3. 提交事务
	if err = tx.Commit().Error; err != nil {
		return nil, err
	}

	resp = new(types.ArticleSaveResp)
	resp.Data = true
	return
}
