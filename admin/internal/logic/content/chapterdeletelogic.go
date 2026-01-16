package content

import (
	"context"
	"errors"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ChapterDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除章节
func NewChapterDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChapterDeleteLogic {
	return &ChapterDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChapterDeleteLogic) ChapterDelete(req *types.ChapterDeleteReq) (resp *types.ChapterDeleteResp, err error) {
	resp = new(types.ChapterDeleteResp)
	result := make(map[string]interface{})

	// 先查询章节信息，获取 book_id，用于删除电子书缓存
	var chapter mysql.TxyChapter
	err = l.svcCtx.DB.Where("id = ?", req.Id).First(&chapter).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("章节不存在")
		}
		logx.Errorf("查询章节失败：%v", err)
		return nil, err
	}
	bookId := chapter.BookId

	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 先删除章节内容（可能不存在，若不存在不视为错误）
		if err := tx.Where("id = ?", req.Id).Delete(&mysql.TxyChapterData{}).Error; err != nil {
			return err
		}

		// 删除章节（一定存在，若不存在则返回错误）
		chapterTx := tx.Where("id = ?", req.Id).Delete(&mysql.TxyChapter{})
		if chapterTx.Error != nil {
			return chapterTx.Error
		}
		if chapterTx.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("章节不存在")
		}
		logx.Errorf("删除章节失败：%v", err)
		return nil, err
	}

	// 删除电子书缓存（如果 Redis 可用）
	if l.svcCtx.Rds != nil && bookId > 0 {
		cacheUtil := utils.NewCacheUtil(l.svcCtx.Rds)
		if err = cacheUtil.DeleteBookCache(l.ctx, bookId); err != nil {
			logx.Errorf("删除电子书缓存失败: bookId=%d, err=%v", bookId, err)
			// 缓存删除失败不影响主流程，只记录日志
		}
	}

	result["status"] = "success"
	resp.Data = result
	return
}
