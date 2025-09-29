package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	lxtPaymentOrdersFieldNames          = builder.RawFieldNames(&LxtPaymentOrders{})
	lxtPaymentOrdersRows                = strings.Join(lxtPaymentOrdersFieldNames, ",")
	lxtPaymentOrdersRowsExpectAutoSet   = strings.Join(stringx.Remove(lxtPaymentOrdersFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	lxtPaymentOrdersRowsWithPlaceHolder = strings.Join(stringx.Remove(lxtPaymentOrdersFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	lxtPaymentOrdersModel interface {
		Insert(ctx context.Context, data *LxtPaymentOrders) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*LxtPaymentOrders, error)
		FindOneByOutTradeNo(ctx context.Context, outTradeNo string) (*LxtPaymentOrders, error)
		FindOneByPaymentId(ctx context.Context, paymentId string) (*LxtPaymentOrders, error)
		Update(ctx context.Context, data *LxtPaymentOrders) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultLxtPaymentOrdersModel struct {
		conn  sqlx.SqlConn
		table string
	}

	LxtPaymentOrders struct {
		Id            uint64       `db:"id"`             // 主键ID
		PaymentId     string       `db:"payment_id"`     // 支付ID
		OrderId       string       `db:"order_id"`       // 订单ID
		OutTradeNo    string       `db:"out_trade_no"`   // 商户订单号
		UserId        int64        `db:"user_id"`        // 用户ID
		Amount        float64      `db:"amount"`         // 支付金额
		Subject       string       `db:"subject"`        // 订单标题
		Body          string       `db:"body"`           // 订单描述
		Status        string       `db:"status"`         // 支付状态
		TradeNo       string       `db:"trade_no"`       // 支付宝交易号
		TradeStatus   string       `db:"trade_status"`   // 支付宝交易状态
		BuyerUserId   string       `db:"buyer_user_id"`  // 买家支付宝用户ID
		BuyerLogonId  string       `db:"buyer_logon_id"` // 买家支付宝账号
		ReceiptAmount string       `db:"receipt_amount"` // 实收金额
		ProductCode   string       `db:"product_code"`   // 产品码
		ReturnUrl     string       `db:"return_url"`     // 支付成功跳转地址
		NotifyUrl     string       `db:"notify_url"`     // 支付结果异步通知地址
		Timeout       string       `db:"timeout"`        // 订单超时时间
		ClientIp      string       `db:"client_ip"`      // 客户端IP
		GmtPayment    sql.NullTime `db:"gmt_payment"`    // 支付时间
		GmtClose      sql.NullTime `db:"gmt_close"`      // 交易关闭时间
		CreatedAt     time.Time    `db:"created_at"`     // 创建时间
		UpdatedAt     time.Time    `db:"updated_at"`     // 更新时间
		DeletedAt     sql.NullTime `db:"deleted_at"`     // 删除时间
	}
)

func newLxtPaymentOrdersModel(conn sqlx.SqlConn) *defaultLxtPaymentOrdersModel {
	return &defaultLxtPaymentOrdersModel{
		conn:  conn,
		table: "`lxt_payment_orders`",
	}
}

func (m *defaultLxtPaymentOrdersModel) Delete(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultLxtPaymentOrdersModel) FindOne(ctx context.Context, id uint64) (*LxtPaymentOrders, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", lxtPaymentOrdersRows, m.table)
	var resp LxtPaymentOrders
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

func (m *defaultLxtPaymentOrdersModel) FindOneByOutTradeNo(ctx context.Context, outTradeNo string) (*LxtPaymentOrders, error) {
	var resp LxtPaymentOrders
	query := fmt.Sprintf("select %s from %s where `out_trade_no` = ? limit 1", lxtPaymentOrdersRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, outTradeNo)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultLxtPaymentOrdersModel) FindOneByPaymentId(ctx context.Context, paymentId string) (*LxtPaymentOrders, error) {
	var resp LxtPaymentOrders
	query := fmt.Sprintf("select %s from %s where `payment_id` = ? limit 1", lxtPaymentOrdersRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, paymentId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultLxtPaymentOrdersModel) Insert(ctx context.Context, data *LxtPaymentOrders) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, lxtPaymentOrdersRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.PaymentId, data.OrderId, data.OutTradeNo, data.UserId, data.Amount, data.Subject, data.Body, data.Status, data.TradeNo, data.TradeStatus, data.BuyerUserId, data.BuyerLogonId, data.ReceiptAmount, data.ProductCode, data.ReturnUrl, data.NotifyUrl, data.Timeout, data.ClientIp, data.GmtPayment, data.GmtClose, data.DeletedAt)
	return ret, err
}

func (m *defaultLxtPaymentOrdersModel) Update(ctx context.Context, newData *LxtPaymentOrders) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, lxtPaymentOrdersRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.PaymentId, newData.OrderId, newData.OutTradeNo, newData.UserId, newData.Amount, newData.Subject, newData.Body, newData.Status, newData.TradeNo, newData.TradeStatus, newData.BuyerUserId, newData.BuyerLogonId, newData.ReceiptAmount, newData.ProductCode, newData.ReturnUrl, newData.NotifyUrl, newData.Timeout, newData.ClientIp, newData.GmtPayment, newData.GmtClose, newData.DeletedAt, newData.Id)
	return err
}

func (m *defaultLxtPaymentOrdersModel) tableName() string {
	return m.table
}
