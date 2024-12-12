package weblogic

import (
	"context"
	"encoding/json"
	"lxtian-blog/rpc/web/model/mysql"

	"lxtian-blog/rpc/web/internal/svc"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleLikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleLikeLogic {
	return &ArticleLikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleLikeLogic) ArticleLike(in *web.ArticleLikeReq) (*web.ArticleLikeResp, error) {
	// 获取文章列表
	//var articles []mysql.TxyArticle
	var articles []map[string]interface{}
	err := l.svcCtx.DB.
		Model(&mysql.TxyArticle{}).
		Select("id,title,path").
		Where("id != ?", in.Id).
		Order("RAND()").
		Limit(4).
		Scan(&articles).Error
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(articles)
	if err != nil {
		return nil, err
	}
	return &web.ArticleLikeResp{
		Data: string(jsonData),
	}, nil
}
