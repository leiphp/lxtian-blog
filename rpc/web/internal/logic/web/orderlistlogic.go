package weblogic

import (
	"context"
	"encoding/json"
	"lxtian-blog/rpc/web/internal/consts"
	"lxtian-blog/rpc/web/internal/svc"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderListLogic {
	return &OrderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderListLogic) OrderList(in *web.OrderListReq) (*web.OrderListResp, error) {
	where := map[string]interface{}{}
	where["status"] = consts.PaymentStatusPaid
	if in.Page == 0 {
		in.Page = 1
	}
	if in.PageSize == 0 {
		in.PageSize = 10
	}
	offset := (in.Page - 1) * in.PageSize
	var results []map[string]interface{}
	err := l.svcCtx.DB.
		Table("txy_order as o").
		Select("o.id,o.amount,o.pay_type,o.user_id,o.status,o.created_at,o.remark,u.nickname,u.head_img").
		Joins("left join txy_user u on u.id = o.user_id").
		Where(where).
		Limit(int(in.PageSize)).
		Offset(int(offset)).
		Order("o.id desc").
		Debug().
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}
	//计算当前type的总数，给分页算总页
	var total int64
	err = l.svcCtx.DB.Table("txy_order as o").Where(where).Count(&total).Error
	if err != nil {
		return nil, err
	}

	return &web.OrderListResp{
		Page:     in.Page,
		PageSize: in.PageSize,
		Total:    uint32(total),
		List:     string(jsonData),
	}, nil
}
