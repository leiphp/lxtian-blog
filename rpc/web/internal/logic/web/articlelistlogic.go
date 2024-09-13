package weblogic

import (
	"context"
	"lxtian-blog/rpc/web/model/mysql"

	"lxtian-blog/rpc/web/internal/svc"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleListLogic {
	return &ArticleListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleListLogic) ArticleList(in *web.ArticleListReq) (*web.ArticleListResp, error) {
	where := map[string]interface{}{}
	if in.Cid > 0 {
		where["cid"] = in.Cid
	}
	if in.Page == 0 {
		in.Page = 1
	}
	if in.PageSize == 0 {
		in.PageSize = 10
	}
	offset := (in.Page - 1) * in.PageSize
	var lxtArticles []mysql.TxyArticle
	err := l.svcCtx.DB.Where(where).Limit(int(in.PageSize)).Offset(int(offset)).Order("id desc").Debug().Find(&lxtArticles).Error
	if err != nil {
		return nil, err
	}

	// 创建一个空的切片来存放转换后的结果
	var articleList []*web.Article
	for _, article := range lxtArticles {
		convertedArticle := &web.Article{
			Id:          uint32(article.Id),
			Title:       article.Title,
			Author:      article.Author,
			Keywords:    article.Keywords,
			Description: article.Description,
			Content:     article.Content,
		}
		// 将转换后的 Article 指针添加到切片中
		articleList = append(articleList, convertedArticle)
	}

	//计算当前type的总数，给分页算总页
	var total int64
	err = l.svcCtx.DB.Model(&mysql.TxyArticle{}).Where(where).Count(&total).Error
	if err != nil {
		return nil, err
	}

	return &web.ArticleListResp{
		Page:     in.Page,
		PageSize: in.PageSize,
		Total:    uint32(total),
		List:     articleList,
	}, nil
}
