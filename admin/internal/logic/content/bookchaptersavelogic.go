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
	// req.Id 一定存在，因为 txy_chapter_data 是对 txy_chapter 的补充
	if req.Id == 0 {
		return nil, fmt.Errorf("章节ID不能为空")
	}

	now := time.Now()
	data := mysql.TxyChapterData{
		Id:      uint64(req.Id), // 必须设置 ID，因为表结构要求
		Title:   req.Title,
		Author:  req.Author,
		Content: req.Content,
		UpdatedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
	}

	// 判断是新增还是更新：cate == "insert" 表示新增，否则表示更新
	isInsert := req.Cate == "insert"

	if isInsert {
		// 新增：使用 req.Id 作为 id
		data.CreatedAt = sql.NullTime{
			Time:  now,
			Valid: true,
		}
		if err = l.svcCtx.DB.Create(&data).Error; err != nil {
			l.Errorf("创建章节失败: id=%d, err=%v", req.Id, err)
			return nil, err
		}
	} else {
		// 更新：使用 req.Id 作为 id
		result := l.svcCtx.DB.Model(&mysql.TxyChapterData{}).Where("id = ?", req.Id).Updates(data)
		if result.Error != nil {
			l.Errorf("更新章节失败: id=%d, err=%v", req.Id, result.Error)
			return nil, result.Error
		}
		// 检查是否有记录被更新，如果没有记录则尝试插入
		if result.RowsAffected == 0 {
			// 记录不存在，转为插入
			data.CreatedAt = sql.NullTime{
				Time:  now,
				Valid: true,
			}
			if err = l.svcCtx.DB.Create(&data).Error; err != nil {
				l.Errorf("插入章节失败: id=%d, err=%v", req.Id, err)
				return nil, err
			}
			l.Infof("章节记录不存在，已自动插入: id=%d", req.Id)
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
