package content

import (
	"context"
	"database/sql"
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

	// 判断是新增还是更新：req.Id == 0 表示新增，否则表示更新
	if req.Id == 0 {
		// 新增
		data.CreatedAt = sql.NullTime{
			Time:  now,
			Valid: true,
		}
		if err = l.svcCtx.DB.Create(&data).Error; err != nil {
			l.Errorf("创建章节失败: %v", err)
			return nil, err
		}
	} else {
		// 更新
		data.Id = uint64(req.Id)
		if err = l.svcCtx.DB.Model(&mysql.TxyChapterData{}).Where("id = ?", req.Id).Updates(data).Error; err != nil {
			l.Errorf("更新章节失败: id=%d, err=%v", req.Id, err)
			return nil, err
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
