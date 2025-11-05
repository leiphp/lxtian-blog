package content

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
	"lxtian-blog/common/pkg/model/mysql"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChapterDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 书单详情
func NewChapterDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChapterDetailLogic {
	return &ChapterDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChapterDetailLogic) ChapterDetail(req *types.ChapterDetailReq) (resp *types.ChapterDetailResp, err error) {
	resp = new(types.ChapterDetailResp)
	result := make(map[string]interface{})
	err = l.svcCtx.DB.
		Model(&mysql.TxyChapterData{}).
		Select("*").
		Where("id = ?", req.Id).
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("数据不存在")
		}
		return nil, err // 其他数据库错误
	}
	resp.Data = result

	return
}
