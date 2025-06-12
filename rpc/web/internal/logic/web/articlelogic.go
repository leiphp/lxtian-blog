package weblogic

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"go.mongodb.org/mongo-driver/mongo"
	model "lxtian-blog/common/pkg/model/mongo"
	"lxtian-blog/common/pkg/model/mysql"
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
		Table("txy_article as a").
		Select("a.id,a.title,a.author,a.description,a.keywords,a.content,a.cid,a.tid,a.mid,a.view_count, DATE_FORMAT(a.created_at, '%Y-%m-%d %H:%i:%s') AS created_at, DATE_FORMAT(a.updated_at, '%Y-%m-%d %H:%i:%s') AS updated_at,c.name category_name").
		Joins("left join txy_category c on c.id = a.cid").
		Where(where).
		Debug().
		Find(&article).Error
	if err != nil {
		return nil, err
	}
	//判断数据是否初始化
	if len(article) > 0 {
		// mongodb获取文章内容
		conn := model.NewArticleModel(l.svcCtx.MongoUri, l.svcCtx.Config.MongoDB.DATABASE, "txy_article")
		contentId, ok := article["mid"].(string)
		if !ok {
			// 处理类型断言失败的情况
			return nil, errors.New("content_id is not a string")
		}
		if string(contentId) != "" {
			res, err := conn.FindOne(l.ctx, string(contentId))
			if err != nil {
				if err != mongo.ErrNoDocuments && err.Error() != "invalid objectId" {
					logc.Errorf(l.ctx, "Document not found or invalid ObjectId: %s", err)
					// 对于其他类型的错误，仍然返回
					return nil, err
				}
			} else {
				article["content"] = res.Content
			}
		}
	}
	// 查询上一篇文章
	var previousArticle mysql.TxyArticle
	err = l.svcCtx.DB.
		Model(&mysql.TxyArticle{}).
		Select("id,title").
		Where("id < ?", in.Id).
		Order("id DESC").
		Limit(1).
		Scan(&previousArticle).Error
	if err != nil {
		return nil, err
	}
	article["prev"] = map[string]interface{}{
		"id":    previousArticle.Id,
		"title": previousArticle.Title,
	}

	// 查询下一篇文章的 ID
	var nextArticle mysql.TxyArticle
	err = l.svcCtx.DB.
		Model(&mysql.TxyArticle{}).
		Select("id,title").
		Where("id > ?", in.Id).
		Order("id ASC").
		Limit(1).
		Scan(&nextArticle).Error
	if err != nil {
		return nil, err
	}
	article["next"] = map[string]interface{}{
		"id":    nextArticle.Id,
		"title": nextArticle.Title,
	}

	// 转换 _id 字段的类型
	if idBytes, ok := article["mid"].([]byte); ok {
		article["mid"] = string(idBytes)
	}
	jsonData, err := json.Marshal(article)
	if err != nil {
		return nil, err
	}
	return &web.ArticleResp{
		Data: string(jsonData),
	}, nil
}
