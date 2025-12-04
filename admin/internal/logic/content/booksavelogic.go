package content

import (
	"context"
	"database/sql"
	"encoding/json"
	"lxtian-blog/common/pkg/model/mysql"
	"time"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BookSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 书单保存
func NewBookSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BookSaveLogic {
	return &BookSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BookSaveLogic) BookSave(req *types.BookSaveReq) (resp *types.BookSaveResp, err error) {
	var status int64 = 0
	if req.Status {
		status = 1
	}
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

	// 序列化 tags 为 JSON 字符串
	tagsJSON, err := json.Marshal(req.Tags)
	if err != nil {
		// 处理错误
		return nil, err
	}

	// 1.判断是插入还是更新
	data := mysql.TxyBook{
		ColumnId:    req.ColumnId,
		Title:       req.Title,
		Slug:        req.Slug,
		Description: req.Description,
		Author:      req.Author,
		Tags: sql.NullString{
			String: string(tagsJSON),
			Valid:  len(req.Tags) > 0, // 如果 tags 为空，Valid 为 false
		},
		Status: status,
		Cover:  req.Cover,
		Badge:  int64(req.Badge),
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
		if err = tx.Debug().Create(&data).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

	} else {
		data.Id = uint64(req.Id)
		if err = tx.Model(&data).Select("column_id", "title", "slug", "description", "status", "cover", "auth", "tags", "updated_at").Debug().Updates(data).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	// 3. 提交事务
	if err = tx.Commit().Error; err != nil {
		return nil, err
	}

	resp = new(types.BookSaveResp)
	resp.Data = true

	return
}
