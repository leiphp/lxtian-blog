// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.2

package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	txyOrderFieldNames          = builder.RawFieldNames(&TxyOrder{})
	txyOrderRows                = strings.Join(txyOrderFieldNames, ",")
	txyOrderRowsExpectAutoSet   = strings.Join(stringx.Remove(txyOrderFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	txyOrderRowsWithPlaceHolder = strings.Join(stringx.Remove(txyOrderFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	txyOrderModel interface {
		Insert(ctx context.Context, data *TxyOrder) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*TxyOrder, error)
		Update(ctx context.Context, data *TxyOrder) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultTxyOrderModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TxyOrder struct {
		Id         uint64  `db:"id"`
		OutTradeNo string  `db:"out_trade_no"` // 商户订单号
		OrderSn    string  `db:"order_sn"`     // 订单号
		PayMoney   float64 `db:"pay_money"`    // 支付金额
		PayType    int64   `db:"pay_type"`     // 1:捐赠2:购买模板3:直接消费
		UserId     uint64  `db:"user_id"`
		Status     int64   `db:"status"`     // 支付状态：0待支付1已支付2已发货3交易完成
		GoodsName  string  `db:"goods_name"` // 商品名称
		Ctime      int64   `db:"ctime"`      // 创建时间
		Remark     string  `db:"remark"`
	}
)

func newTxyOrderModel(conn sqlx.SqlConn) *defaultTxyOrderModel {
	return &defaultTxyOrderModel{
		conn:  conn,
		table: "`txy_order`",
	}
}

func (m *defaultTxyOrderModel) Delete(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultTxyOrderModel) FindOne(ctx context.Context, id uint64) (*TxyOrder, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", txyOrderRows, m.table)
	var resp TxyOrder
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTxyOrderModel) Insert(ctx context.Context, data *TxyOrder) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, txyOrderRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.OutTradeNo, data.OrderSn, data.PayMoney, data.PayType, data.UserId, data.Status, data.GoodsName, data.Ctime, data.Remark)
	return ret, err
}

func (m *defaultTxyOrderModel) Update(ctx context.Context, data *TxyOrder) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, txyOrderRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.OutTradeNo, data.OrderSn, data.PayMoney, data.PayType, data.UserId, data.Status, data.GoodsName, data.Ctime, data.Remark, data.Id)
	return err
}

func (m *defaultTxyOrderModel) tableName() string {
	return m.table
}
