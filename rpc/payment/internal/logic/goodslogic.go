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

type GoodsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	goodsService payment_repo.LxtPaymentGoodsRepo
}

func NewGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GoodsLogic {
	return &GoodsLogic{
		ctx:          ctx,
		svcCtx:       svcCtx,
		Logger:       logx.WithContext(ctx),
		goodsService: payment_repo.NewLxtPaymentGoodsRepo(svcCtx.DB),
	}
}

// 商品详情
func (l *GoodsLogic) Goods(in *payment.GoodsReq) (*payment.GoodsResp, error) {
	// 构建查询条件
	condition := make(map[string]interface{})

	// 状态：只查询已发布的商品
	condition["id"] = in.Id

	// 使用基础仓储的 GetList 方法
	goods, err := l.goodsService.GetByID(l.ctx, in.Id)
	if err != nil {
		l.Errorf("Failed to get goods: %v", err)
		return nil, fmt.Errorf("failed to get goods: %w", err)
	}

	// 构建响应数据
	goodItem := l.buildGoodsItem(goods)

	// 转换为JSON字符串
	jsonData, err := json.Marshal(goodItem)
	if err != nil {
		l.Errorf("Failed to marshal goods: %v", err)
		return nil, fmt.Errorf("failed to marshal goods: %w", err)
	}

	return &payment.GoodsResp{
		Data: string(jsonData),
	}, nil
}

// 构建商品项
func (l *GoodsLogic) buildGoodsItem(good *model.LxtPaymentGood) map[string]interface{} {
	item := map[string]interface{}{
		"id":             good.ID,
		"name":           good.Name,
		"desc":           good.Desc,
		"classify_id":    good.ClassifyID,
		"price":          good.Price,
		"original_price": good.OriginalPrice,
		"rating":         good.Rating,
		"sales":          good.Sales,
		"downloads":      good.Downloads,
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
