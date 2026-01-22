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

	// 判断是新增还是更新：id == 0 表示新增，id > 0 表示更新
	isInsert := req.Id == 0

	data := mysql.TxyChapter{
		BookId: uint64(req.BookId),
		Title:  req.Title,
		UpdatedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
	}

	// 设置可选字段
	if req.IsGroup != nil {
		data.IsGroup = int64(*req.IsGroup)
	}

	if isInsert {
		// 新增章节
		data.CreatedAt = sql.NullTime{
			Time:  now,
			Valid: true,
		}
		if err = l.svcCtx.DB.Create(&data).Error; err != nil {
			l.Errorf("创建章节失败: err=%v", err)
			return nil, err
		}
		l.Infof("创建章节成功: id=%d", data.Id)
	} else {
		// 更新章节
		data.Id = uint64(req.Id)
		result := l.svcCtx.DB.Model(&mysql.TxyChapter{}).Where("id = ?", req.Id).Updates(data)
		if result.Error != nil {
			l.Errorf("更新章节失败: id=%d, err=%v", req.Id, result.Error)
			return nil, result.Error
		}
		// 检查是否有记录被更新
		if result.RowsAffected == 0 {
			l.Errorf("更新章节失败: id=%d 的记录不存在", req.Id)
			return nil, fmt.Errorf("章节不存在: id=%d", req.Id)
		}
		l.Infof("更新章节成功: id=%d", req.Id)
	}

	// 删除电子书缓存（如果 Redis 可用）
	if l.svcCtx.Rds != nil && req.BookId > 0 {
		cacheUtil := utils.NewCacheUtil(l.svcCtx.Rds)
		if err = cacheUtil.DeleteBookCache(l.ctx, uint64(req.BookId)); err != nil {
			logx.Errorf("删除电子书缓存失败: bookId=%d, err=%v", req.BookId, err)
			// 缓存删除失败不影响主流程，只记录日志
		}
	}

	resp = new(types.BookChapterSaveResp)
	resp.Data = true
	return
}
