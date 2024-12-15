package weblogic

import (
	"context"
	"encoding/json"
	"lxtian-blog/rpc/web/internal/consts"
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
	db := l.svcCtx.DB.Model(&mysql.TxyArticle{})
	order := "id desc"
	if in.Cid > 0 {
		where["cid"] = in.Cid
	}
	if in.Tid > 0 {
		where["tid"] = in.Tid
	}
	if in.Types > 0 {
		switch in.Types {
		case consts.ArticleTypesRecommend:
			where["is_tuijian"] = 1
		case consts.ArticleTypesRank:
			order = "view_count desc"
		default:

		}
	}
	if in.Keywords != "" {
		db = db.Where("title like ?", "%"+in.Keywords+"%")
	}
	if in.Page == 0 {
		in.Page = 1
	}
	if in.PageSize == 0 {
		in.PageSize = 10
	}
	offset := (in.Page - 1) * in.PageSize
	var articles []map[string]interface{}
	err := db.Select("txy_article.id,txy_article.path,txy_article.title,txy_article.author,txy_article.description,txy_article.keywords,txy_article.cid,txy_article.tid,txy_article.created_at,txy_article.view_count,c.name cname,t.name tname").
		Joins("left join txy_category as c on txy_article.cid = c.id").
		Joins("left join txy_tag as t on txy_article.tid = t.id").
		Where(where).
		Limit(int(in.PageSize)).
		Offset(int(offset)).
		Order(order).
		Debug().
		Find(&articles).Error
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(articles)
	if err != nil {
		return nil, err
	}
	//计算当前type的总数，给分页算总页
	var total int64
	err = db.Where(where).Count(&total).Error
	if err != nil {
		return nil, err
	}

	return &web.ArticleListResp{
		Page:     in.Page,
		PageSize: in.PageSize,
		Total:    uint32(total),
		List:     string(jsonData),
	}, nil
}
