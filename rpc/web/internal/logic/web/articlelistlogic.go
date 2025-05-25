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
	// 基础查询构建（包含JOIN和公共WHERE条件）
	baseDB := l.svcCtx.DB.Model(&mysql.TxyArticle{}).
		Joins("left join txy_category as c on txy_article.cid = c.id").
		Joins("left join txy_tag as t on txy_article.tid = t.id")

	order := "id desc"
	// 填充WHERE条件
	if in.Cid > 0 {
		where["cid"] = in.Cid
	}
	if in.Tid > 0 {
		where["tid"] = in.Tid
	}
	if in.Types > 0 {
		switch in.Types {
		case consts.ArticleTypesRecommend:
			where["is_rec"] = 1
		case consts.ArticleTypesRank:
			order = "view_count desc"
		default:

		}
	}
	if len(where) > 0 {
		baseDB = baseDB.Where(where)
	}
	if in.Keywords != "" {
		baseDB = baseDB.Where("txy_article.title like ?", "%"+in.Keywords+"%")
	}

	// 计算总数（使用基础查询，无分页/排序）
	var total int64
	if err := baseDB.Count(&total).Error; err != nil {
		return nil, err
	}

	// 处理分页参数
	if in.Page == 0 {
		in.Page = 1
	}
	if in.PageSize == 0 {
		in.PageSize = 10
	}
	offset := (in.Page - 1) * in.PageSize
	var articles []map[string]interface{}
	err := baseDB.Select("txy_article.id,txy_article.path,txy_article.title,txy_article.author,txy_article.description,txy_article.keywords,txy_article.cid,txy_article.tid,txy_article.created_at,txy_article.view_count,c.name cname,t.name tname").
		Limit(int(in.PageSize)).
		Offset(int(offset)).
		Order(order).
		Find(&articles).Error
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(articles)
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
