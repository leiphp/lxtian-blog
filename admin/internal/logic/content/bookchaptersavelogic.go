package content

import (
	"context"
	"database/sql"
	"fmt"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/pkg/utils"
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
	now := time.Now()
	data := mysql.TxyChapterData{
		Title:   req.Title,
		Author:  req.Author,
		Content: req.Content,
		UpdatedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
	}

	// 判断是新增还是更新：
	// 1. 如果 cate == "insert"，强制新增（忽略 id）
	// 2. 否则，req.Id == 0 表示新增，否则表示更新
	isInsert := req.Cate == "insert" || req.Id == 0

	if isInsert {
		// 新增
		data.CreatedAt = sql.NullTime{
			Time:  now,
			Valid: true,
		}
		if err = l.svcCtx.DB.Create(&data).Error; err != nil {
			l.Errorf("创建章节失败: %v", err)
			return nil, err
		}
		// Create 后会自动填充 data.Id（自增ID）
	} else {
		// 更新
		data.Id = uint64(req.Id)
		result := l.svcCtx.DB.Model(&mysql.TxyChapterData{}).Where("id = ?", req.Id).Updates(data)
		if result.Error != nil {
			l.Errorf("更新章节失败: id=%d, err=%v", req.Id, result.Error)
			return nil, result.Error
		}
		// 检查是否有记录被更新
		if result.RowsAffected == 0 {
			l.Errorf("更新章节失败: id=%d 的记录不存在", req.Id)
			return nil, fmt.Errorf("章节不存在: id=%d", req.Id)
		}
	}
	// 删除缓存（如果 Redis 可用）
	if l.svcCtx.Rds != nil {
		cacheUtil := utils.NewCacheUtil(l.svcCtx.Rds)
		if err = cacheUtil.DeleteChapterCache(l.ctx, data.Id); err != nil {
			logx.Errorf("删除章节缓存失败: %v", err)
			// 缓存删除失败不影响主流程，只记录日志
		}
	}

	resp = new(types.BookChapterSaveResp)
	resp.Data = true
	return
}
