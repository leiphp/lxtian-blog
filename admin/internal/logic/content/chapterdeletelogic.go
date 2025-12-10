package content

import (
	"context"
	"errors"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
	"lxtian-blog/common/pkg/model/mysql"

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

	result["status"] = "success"
	resp.Data = result
	return
}
