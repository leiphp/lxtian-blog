package weblogic

import (
	"context"
	"encoding/json"
	"lxtian-blog/rpc/web/model/mysql"

	"lxtian-blog/rpc/web/internal/svc"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleLogic {
	return &ArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleLogic) Article(in *web.ArticleReq) (*web.ArticleResp, error) {
	where := map[string]interface{}{}
	where["id"] = in.Id
	var article map[string]interface{}
	err := l.svcCtx.DB.
		Model(&mysql.TxyArticle{}).
		Select("id,title,author,description,keywords,cid").
		Where(where).
		Debug().
		Find(&article).Error
	if err != nil {
		return nil, err
	}
	//判断数据是否初始化
	if len(article) > 0 {
		article["content"] = "这是测试文章内容"
	}
	jsonData, err := json.Marshal(article)
	if err != nil {
		return nil, err
	}
	return &web.ArticleResp{
		Data: string(jsonData),
	}, nil
}
