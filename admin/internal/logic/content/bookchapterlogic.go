package content

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/pkg/utils"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BookChapterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 书单章节
func NewBookChapterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BookChapterLogic {
	return &BookChapterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BookChapterLogic) BookChapter(req *types.BookChapterReq) (resp *types.BookChapterResp, err error) {
	resp = new(types.BookChapterResp)
	var results []map[string]interface{}
	err = l.svcCtx.DB.
		Model(&mysql.TxyChapter{}).
		Select("id,title,parent_id,is_group").
		Where("book_id = ?", req.Id).
		Find(&results).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("章节不存在")
		}
		return nil, err // 其他数据库错误
	}
	tree := utils.BuildTreeMap(results, 0)
	resp.Data = tree
	return
}
