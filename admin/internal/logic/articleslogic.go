package logic

import (
	"context"
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
	baseDB := l.svcCtx.DB.Table("txy_article as a").
		Joins("left join txy_category as c on c.id = a.cid").
		Joins("left join txy_tag as t on t.id = a.tid")
	order := "a.id desc"
	// 填充WHERE条件
	if req.Cid != 0 {
		baseDB = baseDB.Where("a.`cid` = ?", req.Cid)
	}
	if req.Keywords != "" {
		baseDB = baseDB.Where("a.title like ?", "%"+req.Keywords+"%")
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
	err = baseDB.Select("a.id,a.title,a.keywords,a.author,a.path,a.status,a.click,a.view_count,a.is_top,a.is_rec,a.is_hot,a.is_original,a.created_at,c.name cname,t.name tname").
		Limit(req.PageSize).
		Offset(offset).
		Order(order).
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	// 转换 []byte -> string（特别是中文字段）
	utils.ConvertByteFieldsToString(results)
	utils.FormatTimeFields(results, "created_at")
	utils.FormatBoolFields(results, "status", "is_hot", "is_rec", "is_top", "is_original")

	resp = new(types.ArticlesResp)
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	resp.Total = total
	resp.List = results
	return
}
