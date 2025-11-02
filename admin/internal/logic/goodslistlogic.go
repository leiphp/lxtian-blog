package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
	"lxtian-blog/common/model"
	"lxtian-blog/common/repository/payment_repo"

	"github.com/zeromicro/go-zero/core/logx"
)

type GoodsListLogic struct {
	logx.Logger
	ctx                 context.Context
	svcCtx              *svc.ServiceContext
	lxtPaymentGoodsRepo payment_repo.LxtPaymentGoodsRepo
}

// 商品管理
func NewGoodsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GoodsListLogic {
	return &GoodsListLogic{
		Logger:              logx.WithContext(ctx),
		ctx:                 ctx,
		svcCtx:              svcCtx,
		lxtPaymentGoodsRepo: payment_repo.NewLxtPaymentGoodsRepo(svcCtx.DB),
	}
}

func (l *GoodsListLogic) GoodsList(req *types.GoodsListReq) (resp *types.GoodsListResp, err error) {
	// 参数验证
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100 // 限制最大每页数量
	}

	// 构建查询条件
	condition := make(map[string]interface{})

	// 状态：只查询已发布的商品
	condition["status"] = 1

	// 分类ID过滤
	if req.ClassifyId > 0 {
		condition["classify_id"] = req.ClassifyId
	}

	// 价格范围过滤（需要在查询后处理，因为 GetList 不支持 BETWEEN）
	// 或者使用 GetDB 进行自定义查询

	// 排序字段
	orderBy := req.OrderBy
	if orderBy == "" {
		orderBy = "id desc"
	}

	// 使用基础仓储的 GetList 方法
	goodsList, total, err := l.lxtPaymentGoodsRepo.GetList(l.ctx, condition, int(req.Page), int(req.PageSize), orderBy, req.Keywords, "name", "desc", "detail")
	if err != nil {
		l.Errorf("Failed to get goods list: %v", err)
		return nil, fmt.Errorf("failed to get goods list: %w", err)
	}

	// 价格范围过滤
	var filteredGoods []*model.LxtPaymentGood
	for _, good := range goodsList {
		// 价格过滤
		if req.PriceMin > 0 && good.Price < float64(req.PriceMin) {
			continue
		}
		if req.PriceMax > 0 && good.Price > float64(req.PriceMax) {
			continue
		}
		filteredGoods = append(filteredGoods, good)
	}

	// 如果进行了价格过滤，需要重新计算总数
	if req.PriceMin > 0 || req.PriceMax > 0 {
		// 重新获取所有符合条件的记录以计算准确的总数
		allGoods, _, err := l.lxtPaymentGoodsRepo.GetList(l.ctx, condition, 0, 0, "", req.Keywords, "name", "desc", "detail")
		if err == nil {
			var count int64
			for _, good := range allGoods {
				if req.PriceMin > 0 && good.Price < float64(req.PriceMin) {
					continue
				}
				if req.PriceMax > 0 && good.Price > float64(req.PriceMax) {
					continue
				}
				count++
			}
			total = count
		}
	}

	// 构建响应数据
	var listData []map[string]interface{}
	for _, good := range filteredGoods {
		item := l.buildGoodsItem(good)
		listData = append(listData, item)
	}

	resp = &types.GoodsListResp{
		BasePageRes: types.BasePageRes{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
			List:     listData,
		},
	}
	return resp, nil
}

// 构建商品项
func (l *GoodsListLogic) buildGoodsItem(good *model.LxtPaymentGood) map[string]interface{} {
	item := map[string]interface{}{
		"id":             good.ID,
		"name":           good.Name,
		"desc":           good.Desc,
		"classify_id":    good.ClassifyID,
		"price":          good.Price,
		"original_price": good.OriginalPrice,
		"rating":         good.Rating,
		"sales":          good.Sales,
		"download":       good.Download,
		"size":           good.Size,
		"status":         good.Status,
		"created_at":     good.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at":     good.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if good.Detail != nil {
		item["detail"] = *good.Detail
	}
	if good.ProductCode != nil {
		item["product_code"] = *good.ProductCode
	}
	if good.PicURL != nil {
		item["pic_url"] = *good.PicURL
	}
	if good.Tags != nil {
		// 将 tags JSON 字符串解析为数组
		var tagsArray []string
		if err := json.Unmarshal([]byte(*good.Tags), &tagsArray); err == nil {
			item["tags"] = tagsArray
		} else {
			// 如果解析失败，设置为空数组
			item["tags"] = []string{}
		}
	} else {
		// 如果 tags 为 nil，设置为空数组
		item["tags"] = []string{}
	}

	return item
}
