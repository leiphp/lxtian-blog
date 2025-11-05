package content

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
	"lxtian-blog/common/pkg/model/mysql"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文章详情
func NewArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleLogic {
	return &ArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleLogic) Article(req *types.ArticleReq) (resp *types.ArticleResp, err error) {
	resp = new(types.ArticleResp)
	var result map[string]interface{}
	err = l.svcCtx.DB.
		Model(&mysql.TxyArticle{}).
		Select("*").
		Where("id = ?", req.Id).
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("数据不存在")
		}
		return nil, err // 其他数据库错误
	}
	path := result["path"].(string)
	if !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") {
		path = l.svcCtx.QiniuClient.PrivateURL(path, 3600)
	}
	result["path"] = path
	resp.Data = result
	return
}
