package logic

import (
	"context"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/pkg/utils"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticlesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文章管理
func NewArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticlesLogic {
	return &ArticlesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticlesLogic) Articles(req *types.ArticlesReq) (resp *types.ArticlesResp, err error) {
	// 基础查询构建（包含JOIN和公共WHERE条件）
	baseDB := l.svcCtx.DB.Model(&mysql.TxyArticle{}).
		Joins("left join txy_category as c on c.id = txy_article.cid").
		Joins("left join txy_tag as t on t.id = txy_article.tid")
	order := "txy_article.id desc"
	// 填充WHERE条件
	if req.Cid != 0 {
		baseDB = baseDB.Where("txy.`cid` = ?", req.Cid)
	}
	if req.Keywords != "" {
		baseDB = baseDB.Where("txy_user.username like ?", "%"+req.Keywords+"%")
	}
	// 计算总数（使用基础查询，无分页/排序）
	var total int64
	if err = baseDB.Count(&total).Error; err != nil {
		return nil, err
	}

	// 处理分页参数
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	offset := (req.Page - 1) * req.PageSize
	var results []map[string]interface{}
	err = baseDB.Select("txy_article.id,txy_article.title,txy_article.author,txy_article.path,txy_article.status,txy_article.created_at,txy_article.view_count,txy_article.keywords,c.name cname, t.name tname").
		Limit(req.PageSize).
		Offset(offset).
		Order(order).
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	// 转换 []byte -> string（特别是中文字段）
	utils.ConvertByteFieldsToString(results)

	resp = new(types.ArticlesResp)
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	resp.Total = total
	resp.List = results
	return
}
