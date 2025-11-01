package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"lxtian-blog/common/model"
	"lxtian-blog/common/repository/payment_repo"
	"lxtian-blog/rpc/payment/internal/svc"
	"lxtian-blog/rpc/payment/pb/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

type GoodsListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	goodsService payment_repo.LxtPaymentGoodsRepo
}

func NewGoodsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GoodsListLogic {
	return &GoodsListLogic{
		ctx:          ctx,
		svcCtx:       svcCtx,
		Logger:       logx.WithContext(ctx),
		goodsService: payment_repo.NewLxtPaymentGoodsRepo(svcCtx.DB),
	}
}

// 商品列表查询
func (l *GoodsListLogic) GoodsList(in *payment.GoodsListReq) (*payment.GoodsListResp, error) {
	// 参数验证
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}
	if in.PageSize > 100 {
		in.PageSize = 100 // 限制最大每页数量
	}

	// 构建查询条件
	condition := make(map[string]interface{})

	// 状态：只查询已发布的商品
	condition["status"] = 1

	// 分类ID过滤
	if in.ClassifyId > 0 {
		condition["classify_id"] = in.ClassifyId
	}

	// 价格范围过滤（需要在查询后处理，因为 GetList 不支持 BETWEEN）
	// 或者使用 GetDB 进行自定义查询

	// 排序字段
	orderBy := in.OrderBy
	if orderBy == "" {
		orderBy = "id desc"
	}

	// 使用基础仓储的 GetList 方法
	goodsList, total, err := l.goodsService.GetList(l.ctx, condition, int(in.Page), int(in.PageSize), orderBy, in.Keywords, "name", "desc", "detail")
	if err != nil {
		l.Errorf("Failed to get goods list: %v", err)
		return nil, fmt.Errorf("failed to get goods list: %w", err)
	}

	// 价格范围过滤
	var filteredGoods []*model.LxtPaymentGood
	for _, good := range goodsList {
		// 价格过滤
		if in.PriceMin > 0 && good.Price < float64(in.PriceMin) {
			continue
		}
		if in.PriceMax > 0 && good.Price > float64(in.PriceMax) {
			continue
		}
		filteredGoods = append(filteredGoods, good)
	}

	// 如果进行了价格过滤，需要重新计算总数
	if in.PriceMin > 0 || in.PriceMax > 0 {
		// 重新获取所有符合条件的记录以计算准确的总数
		allGoods, _, err := l.goodsService.GetList(l.ctx, condition, 0, 0, "", in.Keywords, "name", "desc", "detail")
		if err == nil {
			var count int64
			for _, good := range allGoods {
				if in.PriceMin > 0 && good.Price < float64(in.PriceMin) {
					continue
				}
				if in.PriceMax > 0 && good.Price > float64(in.PriceMax) {
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

	// 转换为JSON字符串
	listJson, err := json.Marshal(listData)
	if err != nil {
		l.Errorf("Failed to marshal goods list: %v", err)
		return nil, fmt.Errorf("failed to marshal goods list: %w", err)
	}

	return &payment.GoodsListResp{
		Page:     in.Page,
		PageSize: in.PageSize,
		Total:    uint64(total),
		List:     string(listJson),
	}, nil
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
