package weblogic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	model "lxtian-blog/common/pkg/model/mongo"
	"lxtian-blog/common/pkg/model/mysql"
	redisutil "lxtian-blog/common/pkg/redis"
	"lxtian-blog/common/pkg/utils"
	"lxtian-blog/rpc/web/internal/svc"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logc"
	"go.mongodb.org/mongo-driver/mongo"

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
	// 记录浏览次数（如果有IP参数）
	if in.ClientIp != "" {
		go func() {
			viewCountUtil := utils.NewViewCountUtil(l.svcCtx.DB, l.svcCtx.Rds)
			if err := viewCountUtil.IncrementArticleView(l.ctx, in.Id, in.ClientIp); err != nil {
				logc.Errorf(l.ctx, "记录文章浏览次数失败: %s", err)
			}
		}()
	}
	articleID := uint64(in.Id)
	// 1. 尝试从缓存获取
	cachedArticle, err := l.getArticleFromCache(l.ctx, articleID)
	if err == nil && cachedArticle != "" {
		logx.Infof("从缓存获取文章详情: %d", articleID)
		// 将缓存数据转换为JSON字符串
		return &web.ArticleResp{
			Data: cachedArticle,
		}, nil
	}

	where := map[string]interface{}{}
	where["id"] = in.Id
	var article map[string]interface{}
	err = l.svcCtx.DB.
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
	// 异步写入缓存
	err = l.setArticleToCache(l.ctx, articleID, string(jsonData))
	if err != nil {
		logx.Errorf("写入文章缓存失败: %v", err)
		return nil, err
	}

	return &web.ArticleResp{
		Data: string(jsonData),
	}, nil
}

// getArticleCacheKey 获取文章缓存Key
func (l *ArticleLogic) getArticleCacheKey(articleID uint64) string {
	return fmt.Sprintf("%sarticle:detail:%d", redisutil.KeyPrefix, articleID)
}

// getArticleFromCache 从缓存获取文章
func (l *ArticleLogic) getArticleFromCache(ctx context.Context, articleID uint64) (string, error) {
	key := l.getArticleCacheKey(articleID)

	data, err := l.svcCtx.Rds.GetCtx(ctx, key)
	if err != nil {
		if err == redis.Nil {
			return "", nil // 缓存不存在
		}
		logx.Errorf("获取文章缓存失败: %v", err)
		return "", err
	}
	return data, nil
}

// setArticleToCache 设置文章缓存
func (l *ArticleLogic) setArticleToCache(ctx context.Context, articleID uint64, articleStr string) error {
	key := l.getArticleCacheKey(articleID)
	// 设置1小时过期
	err := l.svcCtx.Rds.SetexCtx(ctx, key, articleStr, 3600)
	if err != nil {
		logx.Errorf("设置文章缓存失败: %v", err)
		return err
	}

	logx.Infof("文章 %d 缓存设置成功", articleID)
	return nil
}

// deleteArticleCache 删除文章缓存
func (l *ArticleLogic) deleteArticleCache(ctx context.Context, articleID uint64) error {
	key := l.getArticleCacheKey(articleID)

	_, err := l.svcCtx.Rds.DelCtx(ctx, key)
	if err != nil {
		logx.Errorf("删除文章缓存失败: %v", err)
		return err
	}

	logx.Infof("文章 %d 缓存删除成功", articleID)
	return nil
}
