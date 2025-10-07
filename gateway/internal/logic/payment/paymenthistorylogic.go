package payment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"lxtian-blog/rpc/payment/pb/payment"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PaymentHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 支付记录查询
func NewPaymentHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentHistoryLogic {
	return &PaymentHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PaymentHistoryLogic) PaymentHistory(req *types.PaymentHistoryReq) (resp *types.PaymentHistoryResp, err error) {
	//从中间件获取用户信息
	userId, ok := l.ctx.Value("user_id").(uint)
	if !ok {
		return nil, errors.New("user_id not found in context")
	}
	res, err := l.svcCtx.PaymentRpc.PaymentHistory(l.ctx, &payment.PaymentHistoryReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		UserId:   uint64(userId),
	})
	if err != nil {
		return nil, err
	}

	// 解析JSON字符串为列表
	var list []map[string]interface{}
	if res.List != "" {
		if err := json.Unmarshal([]byte(res.List), &list); err != nil {
			return nil, fmt.Errorf("failed to unmarshal payment list: %w", err)
		}
	}

	return &types.PaymentHistoryResp{
		Page:     res.Page,
		PageSize: res.PageSize,
		Total:    res.Total,
		List:     list,
	}, nil
}
