package model

import (
	"context"
	"database/sql"
	"fmt"

	"lxtian-blog/common/pkg/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PaymentModel = (*customPaymentModel)(nil)

type (
	// PaymentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPaymentModel.
	PaymentModel interface {
		// 支付订单相关方法
		InsertPaymentOrder(ctx context.Context, data *model.PaymentOrder) (sql.Result, error)
		FindPaymentOrderById(ctx context.Context, id uint64) (*model.PaymentOrder, error)
		FindPaymentOrderByPaymentId(ctx context.Context, paymentId string) (*model.PaymentOrder, error)
		FindPaymentOrderByOrderId(ctx context.Context, orderId string) (*model.PaymentOrder, error)
		FindPaymentOrderByOutTradeNo(ctx context.Context, outTradeNo string) (*model.PaymentOrder, error)
		UpdatePaymentOrder(ctx context.Context, data *model.PaymentOrder) error
		UpdatePaymentOrderStatus(ctx context.Context, paymentId string, status string) error
		UpdatePaymentOrderTradeInfo(ctx context.Context, paymentId string, tradeNo string, tradeStatus string, buyerUserId string, buyerLogonId string, receiptAmount float64, gmtPayment interface{}) error
		DeletePaymentOrder(ctx context.Context, id uint64) error
		FindPaymentOrdersByUserId(ctx context.Context, userId uint64, offset, limit int) ([]*model.PaymentOrder, error)
		FindPaymentOrdersByStatus(ctx context.Context, status string, offset, limit int) ([]*model.PaymentOrder, error)
		CountPaymentOrdersByUserId(ctx context.Context, userId uint64) (int64, error)
		CountPaymentOrdersByStatus(ctx context.Context, status string) (int64, error)

		// 支付退款相关方法
		InsertPaymentRefund(ctx context.Context, data *model.PaymentRefund) (sql.Result, error)
		FindPaymentRefundById(ctx context.Context, id uint64) (*model.PaymentRefund, error)
		FindPaymentRefundByRefundId(ctx context.Context, refundId string) (*model.PaymentRefund, error)
		FindPaymentRefundByOutRequestNo(ctx context.Context, outRequestNo string) (*model.PaymentRefund, error)
		UpdatePaymentRefund(ctx context.Context, data *model.PaymentRefund) error
		UpdatePaymentRefundStatus(ctx context.Context, refundId string, status string) error
		DeletePaymentRefund(ctx context.Context, id uint64) error
		FindPaymentRefundsByPaymentId(ctx context.Context, paymentId string, offset, limit int) ([]*model.PaymentRefund, error)
		FindPaymentRefundsByUserId(ctx context.Context, userId uint64, offset, limit int) ([]*model.PaymentRefund, error)
		CountPaymentRefundsByPaymentId(ctx context.Context, paymentId string) (int64, error)
		CountPaymentRefundsByUserId(ctx context.Context, userId uint64) (int64, error)

		// 支付通知相关方法
		InsertPaymentNotify(ctx context.Context, data *model.PaymentNotify) (sql.Result, error)
		FindPaymentNotifyById(ctx context.Context, id uint64) (*model.PaymentNotify, error)
		FindPaymentNotifyByNotifyId(ctx context.Context, notifyId string) (*model.PaymentNotify, error)
		UpdatePaymentNotify(ctx context.Context, data *model.PaymentNotify) error
		UpdatePaymentNotifyVerifyStatus(ctx context.Context, notifyId string, verifyStatus string) error
		UpdatePaymentNotifyProcessStatus(ctx context.Context, notifyId string, processStatus string, errorMessage string) error
		DeletePaymentNotify(ctx context.Context, id uint64) error
		FindPaymentNotifiesByPaymentId(ctx context.Context, paymentId string, offset, limit int) ([]*model.PaymentNotify, error)
		FindPaymentNotifiesByType(ctx context.Context, notifyType string, offset, limit int) ([]*model.PaymentNotify, error)
		CountPaymentNotifiesByPaymentId(ctx context.Context, paymentId string) (int64, error)
		CountPaymentNotifiesByType(ctx context.Context, notifyType string) (int64, error)

		// 支付配置相关方法
		InsertPaymentConfig(ctx context.Context, data *model.PaymentConfig) (sql.Result, error)
		FindPaymentConfigById(ctx context.Context, id uint64) (*model.PaymentConfig, error)
		FindPaymentConfigByAppId(ctx context.Context, appId string) (*model.PaymentConfig, error)
		UpdatePaymentConfig(ctx context.Context, data *model.PaymentConfig) error
		DeletePaymentConfig(ctx context.Context, id uint64) error
		FindAllPaymentConfigs(ctx context.Context, offset, limit int) ([]*model.PaymentConfig, error)
		CountPaymentConfigs(ctx context.Context) (int64, error)

		withSession(session sqlx.Session) PaymentModel
	}

	customPaymentModel struct {
		*defaultPaymentModel
	}

	defaultPaymentModel struct {
		conn                 sqlx.SqlConn
		paymentOrdersTable   string
		paymentRefundsTable  string
		paymentNotifiesTable string
		paymentConfigsTable  string
	}
)

// NewPaymentModel returns a model for the database table.
func NewPaymentModel(conn sqlx.SqlConn) PaymentModel {
	return &customPaymentModel{
		defaultPaymentModel: newPaymentModel(conn),
	}
}

func (m *customPaymentModel) withSession(session sqlx.Session) PaymentModel {
	return NewPaymentModel(sqlx.NewSqlConnFromSession(session))
}

func newPaymentModel(conn sqlx.SqlConn) *defaultPaymentModel {
	return &defaultPaymentModel{
		conn:                 conn,
		paymentOrdersTable:   "lxt_payment_orders",
		paymentRefundsTable:  "lxt_payment_refunds",
		paymentNotifiesTable: "lxt_payment_notifies",
		paymentConfigsTable:  "lxt_payment_configs",
	}
}

// 支付订单相关方法实现

func (m *defaultPaymentModel) InsertPaymentOrder(ctx context.Context, data *model.PaymentOrder) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (payment_id, order_id, out_trade_no, user_id, amount, subject, body, status, product_code, return_url, notify_url, timeout, client_ip) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.paymentOrdersTable)
	return m.conn.ExecCtx(ctx, query, data.PaymentId, data.OrderId, data.OutTradeNo, data.UserId, data.Amount, data.Subject, data.Body, data.Status, data.ProductCode, data.ReturnUrl, data.NotifyUrl, data.Timeout, data.ClientIP)
}

func (m *defaultPaymentModel) FindPaymentOrderById(ctx context.Context, id uint64) (*model.PaymentOrder, error) {
	var resp model.PaymentOrder
	query := fmt.Sprintf("select %s from %s where id = ? limit 1", paymentOrderRows, m.paymentOrdersTable)
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultPaymentModel) FindPaymentOrderByPaymentId(ctx context.Context, paymentId string) (*model.PaymentOrder, error) {
	var resp model.PaymentOrder
	query := fmt.Sprintf("select %s from %s where payment_id = ? limit 1", paymentOrderRows, m.paymentOrdersTable)
	err := m.conn.QueryRowCtx(ctx, &resp, query, paymentId)
	switch err {
	case nil:
		return &resp, nil
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultPaymentModel) FindPaymentOrderByOrderId(ctx context.Context, orderId string) (*model.PaymentOrder, error) {
	var resp model.PaymentOrder
	query := fmt.Sprintf("select %s from %s where order_id = ? limit 1", paymentOrderRows, m.paymentOrdersTable)
	err := m.conn.QueryRowCtx(ctx, &resp, query, orderId)
	switch err {
	case nil:
		return &resp, nil
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultPaymentModel) FindPaymentOrderByOutTradeNo(ctx context.Context, outTradeNo string) (*model.PaymentOrder, error) {
	var resp model.PaymentOrder
	query := fmt.Sprintf("select %s from %s where out_trade_no = ? limit 1", paymentOrderRows, m.paymentOrdersTable)
	err := m.conn.QueryRowCtx(ctx, &resp, query, outTradeNo)
	switch err {
	case nil:
		return &resp, nil
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultPaymentModel) UpdatePaymentOrder(ctx context.Context, data *model.PaymentOrder) error {
	query := fmt.Sprintf("update %s set %s where payment_id = ?", m.paymentOrdersTable, paymentOrderRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.OrderId, data.OutTradeNo, data.UserId, data.Amount, data.Subject, data.Body, data.Status, data.TradeNo, data.TradeStatus, data.BuyerUserId, data.BuyerLogonId, data.ReceiptAmount, data.ProductCode, data.ReturnUrl, data.NotifyUrl, data.Timeout, data.ClientIP, data.GmtPayment, data.GmtClose, data.PaymentId)
	return err
}

func (m *defaultPaymentModel) UpdatePaymentOrderStatus(ctx context.Context, paymentId string, status string) error {
	query := fmt.Sprintf("update %s set status = ? where payment_id = ?", m.paymentOrdersTable)
	_, err := m.conn.ExecCtx(ctx, query, status, paymentId)
	return err
}

func (m *defaultPaymentModel) UpdatePaymentOrderTradeInfo(ctx context.Context, paymentId string, tradeNo string, tradeStatus string, buyerUserId string, buyerLogonId string, receiptAmount float64, gmtPayment interface{}) error {
	query := fmt.Sprintf("update %s set trade_no = ?, trade_status = ?, buyer_user_id = ?, buyer_logon_id = ?, receipt_amount = ?, gmt_payment = ?, status = ? where payment_id = ?", m.paymentOrdersTable)
	_, err := m.conn.ExecCtx(ctx, query, tradeNo, tradeStatus, buyerUserId, buyerLogonId, receiptAmount, gmtPayment, model.PaymentStatusPaid, paymentId)
	return err
}

func (m *defaultPaymentModel) DeletePaymentOrder(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where id = ?", m.paymentOrdersTable)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultPaymentModel) FindPaymentOrdersByUserId(ctx context.Context, userId uint64, offset, limit int) ([]*model.PaymentOrder, error) {
	var resp []*model.PaymentOrder
	query := fmt.Sprintf("select %s from %s where user_id = ? order by created_at desc limit ? offset ?", paymentOrderRows, m.paymentOrdersTable)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId, limit, offset)
	return resp, err
}

func (m *defaultPaymentModel) FindPaymentOrdersByStatus(ctx context.Context, status string, offset, limit int) ([]*model.PaymentOrder, error) {
	var resp []*model.PaymentOrder
	if status == "" {
		query := fmt.Sprintf("select %s from %s order by created_at desc limit ? offset ?", paymentOrderRows, m.paymentOrdersTable)
		err := m.conn.QueryRowsCtx(ctx, &resp, query, limit, offset)
		return resp, err
	}
	query := fmt.Sprintf("select %s from %s where status = ? order by created_at desc limit ? offset ?", paymentOrderRows, m.paymentOrdersTable)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, status, limit, offset)
	return resp, err
}

func (m *defaultPaymentModel) CountPaymentOrdersByUserId(ctx context.Context, userId uint64) (int64, error) {
	var count int64
	query := fmt.Sprintf("select count(*) from %s where user_id = ?", m.paymentOrdersTable)
	err := m.conn.QueryRowCtx(ctx, &count, query, userId)
	return count, err
}

func (m *defaultPaymentModel) CountPaymentOrdersByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	if status == "" {
		query := fmt.Sprintf("select count(*) from %s", m.paymentOrdersTable)
		err := m.conn.QueryRowCtx(ctx, &count, query)
		return count, err
	}
	query := fmt.Sprintf("select count(*) from %s where status = ?", m.paymentOrdersTable)
	err := m.conn.QueryRowCtx(ctx, &count, query, status)
	return count, err
}

// 支付退款相关方法实现

func (m *defaultPaymentModel) InsertPaymentRefund(ctx context.Context, data *model.PaymentRefund) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (refund_id, payment_id, order_id, out_trade_no, out_request_no, user_id, refund_amount, refund_fee, refund_reason, status, refund_status) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.paymentRefundsTable)
	return m.conn.ExecCtx(ctx, query, data.RefundId, data.PaymentId, data.OrderId, data.OutTradeNo, data.OutRequestNo, data.UserId, data.RefundAmount, data.RefundFee, data.RefundReason, data.Status, data.RefundStatus)
}

func (m *defaultPaymentModel) FindPaymentRefundById(ctx context.Context, id uint64) (*model.PaymentRefund, error) {
	var resp model.PaymentRefund
	query := fmt.Sprintf("select %s from %s where id = ? limit 1", paymentRefundRows, m.paymentRefundsTable)
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultPaymentModel) FindPaymentRefundByRefundId(ctx context.Context, refundId string) (*model.PaymentRefund, error) {
	var resp model.PaymentRefund
	query := fmt.Sprintf("select %s from %s where refund_id = ? limit 1", paymentRefundRows, m.paymentRefundsTable)
	err := m.conn.QueryRowCtx(ctx, &resp, query, refundId)
	switch err {
	case nil:
		return &resp, nil
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultPaymentModel) FindPaymentRefundByOutRequestNo(ctx context.Context, outRequestNo string) (*model.PaymentRefund, error) {
	var resp model.PaymentRefund
	query := fmt.Sprintf("select %s from %s where out_request_no = ? limit 1", paymentRefundRows, m.paymentRefundsTable)
	err := m.conn.QueryRowCtx(ctx, &resp, query, outRequestNo)
	switch err {
	case nil:
		return &resp, nil
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultPaymentModel) UpdatePaymentRefund(ctx context.Context, data *model.PaymentRefund) error {
	query := fmt.Sprintf("update %s set %s where refund_id = ?", m.paymentRefundsTable, paymentRefundRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.PaymentId, data.OrderId, data.OutTradeNo, data.OutRequestNo, data.UserId, data.RefundAmount, data.RefundFee, data.RefundReason, data.Status, data.RefundStatus, data.GmtRefund, data.RefundId)
	return err
}

func (m *defaultPaymentModel) UpdatePaymentRefundStatus(ctx context.Context, refundId string, status string) error {
	query := fmt.Sprintf("update %s set status = ? where refund_id = ?", m.paymentRefundsTable)
	_, err := m.conn.ExecCtx(ctx, query, status, refundId)
	return err
}

func (m *defaultPaymentModel) DeletePaymentRefund(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where id = ?", m.paymentRefundsTable)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultPaymentModel) FindPaymentRefundsByPaymentId(ctx context.Context, paymentId string, offset, limit int) ([]*model.PaymentRefund, error) {
	var resp []*model.PaymentRefund
	query := fmt.Sprintf("select %s from %s where payment_id = ? order by created_at desc limit ? offset ?", paymentRefundRows, m.paymentRefundsTable)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, paymentId, limit, offset)
	return resp, err
}

func (m *defaultPaymentModel) FindPaymentRefundsByUserId(ctx context.Context, userId uint64, offset, limit int) ([]*model.PaymentRefund, error) {
	var resp []*model.PaymentRefund
	query := fmt.Sprintf("select %s from %s where user_id = ? order by created_at desc limit ? offset ?", paymentRefundRows, m.paymentRefundsTable)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId, limit, offset)
	return resp, err
}

func (m *defaultPaymentModel) CountPaymentRefundsByPaymentId(ctx context.Context, paymentId string) (int64, error) {
	var count int64
	query := fmt.Sprintf("select count(*) from %s where payment_id = ?", m.paymentRefundsTable)
	err := m.conn.QueryRowCtx(ctx, &count, query, paymentId)
	return count, err
}

func (m *defaultPaymentModel) CountPaymentRefundsByUserId(ctx context.Context, userId uint64) (int64, error) {
	var count int64
	query := fmt.Sprintf("select count(*) from %s where user_id = ?", m.paymentRefundsTable)
	err := m.conn.QueryRowCtx(ctx, &count, query, userId)
	return count, err
}

// 支付通知相关方法实现

func (m *defaultPaymentModel) InsertPaymentNotify(ctx context.Context, data *model.PaymentNotify) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (notify_id, payment_id, notify_type, notify_data, sign, sign_type, verify_status, process_status, client_ip) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.paymentNotifiesTable)
	return m.conn.ExecCtx(ctx, query, data.NotifyId, data.PaymentId, data.NotifyType, data.NotifyData, data.Sign, data.SignType, data.VerifyStatus, data.ProcessStatus, data.ClientIP)
}

func (m *defaultPaymentModel) FindPaymentNotifyById(ctx context.Context, id uint64) (*model.PaymentNotify, error) {
	var resp model.PaymentNotify
	query := fmt.Sprintf("select %s from %s where id = ? limit 1", paymentNotifyRows, m.paymentNotifiesTable)
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultPaymentModel) FindPaymentNotifyByNotifyId(ctx context.Context, notifyId string) (*model.PaymentNotify, error) {
	var resp model.PaymentNotify
	query := fmt.Sprintf("select %s from %s where notify_id = ? limit 1", paymentNotifyRows, m.paymentNotifiesTable)
	err := m.conn.QueryRowCtx(ctx, &resp, query, notifyId)
	switch err {
	case nil:
		return &resp, nil
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultPaymentModel) UpdatePaymentNotify(ctx context.Context, data *model.PaymentNotify) error {
	query := fmt.Sprintf("update %s set %s where notify_id = ?", m.paymentNotifiesTable, paymentNotifyRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.PaymentId, data.NotifyType, data.NotifyData, data.Sign, data.SignType, data.VerifyStatus, data.ProcessStatus, data.ClientIP, data.ErrorMessage, data.ProcessedAt, data.NotifyId)
	return err
}

func (m *defaultPaymentModel) UpdatePaymentNotifyVerifyStatus(ctx context.Context, notifyId string, verifyStatus string) error {
	query := fmt.Sprintf("update %s set verify_status = ? where notify_id = ?", m.paymentNotifiesTable)
	_, err := m.conn.ExecCtx(ctx, query, verifyStatus, notifyId)
	return err
}

func (m *defaultPaymentModel) UpdatePaymentNotifyProcessStatus(ctx context.Context, notifyId string, processStatus string, errorMessage string) error {
	query := fmt.Sprintf("update %s set process_status = ?, error_message = ?, processed_at = now() where notify_id = ?", m.paymentNotifiesTable)
	_, err := m.conn.ExecCtx(ctx, query, processStatus, errorMessage, notifyId)
	return err
}

func (m *defaultPaymentModel) DeletePaymentNotify(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where id = ?", m.paymentNotifiesTable)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultPaymentModel) FindPaymentNotifiesByPaymentId(ctx context.Context, paymentId string, offset, limit int) ([]*model.PaymentNotify, error) {
	var resp []*model.PaymentNotify
	query := fmt.Sprintf("select %s from %s where payment_id = ? order by created_at desc limit ? offset ?", paymentNotifyRows, m.paymentNotifiesTable)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, paymentId, limit, offset)
	return resp, err
}

func (m *defaultPaymentModel) FindPaymentNotifiesByType(ctx context.Context, notifyType string, offset, limit int) ([]*model.PaymentNotify, error) {
	var resp []*model.PaymentNotify
	query := fmt.Sprintf("select %s from %s where notify_type = ? order by created_at desc limit ? offset ?", paymentNotifyRows, m.paymentNotifiesTable)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, notifyType, limit, offset)
	return resp, err
}

func (m *defaultPaymentModel) CountPaymentNotifiesByPaymentId(ctx context.Context, paymentId string) (int64, error) {
	var count int64
	query := fmt.Sprintf("select count(*) from %s where payment_id = ?", m.paymentNotifiesTable)
	err := m.conn.QueryRowCtx(ctx, &count, query, paymentId)
	return count, err
}

func (m *defaultPaymentModel) CountPaymentNotifiesByType(ctx context.Context, notifyType string) (int64, error) {
	var count int64
	query := fmt.Sprintf("select count(*) from %s where notify_type = ?", m.paymentNotifiesTable)
	err := m.conn.QueryRowCtx(ctx, &count, query, notifyType)
	return count, err
}

// 支付配置相关方法实现

func (m *defaultPaymentModel) InsertPaymentConfig(ctx context.Context, data *model.PaymentConfig) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (app_id, app_name, app_private_key, alipay_public_key, gateway_url, is_prod, is_enabled, notify_url, return_url) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.paymentConfigsTable)
	return m.conn.ExecCtx(ctx, query, data.AppId, data.AppName, data.AppPrivateKey, data.AlipayPublicKey, data.GatewayUrl, data.IsProd, data.IsEnabled, data.NotifyUrl, data.ReturnUrl)
}

func (m *defaultPaymentModel) FindPaymentConfigById(ctx context.Context, id uint64) (*model.PaymentConfig, error) {
	var resp model.PaymentConfig
	query := fmt.Sprintf("select %s from %s where id = ? limit 1", paymentConfigRows, m.paymentConfigsTable)
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultPaymentModel) FindPaymentConfigByAppId(ctx context.Context, appId string) (*model.PaymentConfig, error) {
	var resp model.PaymentConfig
	query := fmt.Sprintf("select %s from %s where app_id = ? limit 1", paymentConfigRows, m.paymentConfigsTable)
	err := m.conn.QueryRowCtx(ctx, &resp, query, appId)
	switch err {
	case nil:
		return &resp, nil
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultPaymentModel) UpdatePaymentConfig(ctx context.Context, data *model.PaymentConfig) error {
	query := fmt.Sprintf("update %s set %s where app_id = ?", m.paymentConfigsTable, paymentConfigRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.AppName, data.AppPrivateKey, data.AlipayPublicKey, data.GatewayUrl, data.IsProd, data.IsEnabled, data.NotifyUrl, data.ReturnUrl, data.AppId)
	return err
}

func (m *defaultPaymentModel) DeletePaymentConfig(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where id = ?", m.paymentConfigsTable)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultPaymentModel) FindAllPaymentConfigs(ctx context.Context, offset, limit int) ([]*model.PaymentConfig, error) {
	var resp []*model.PaymentConfig
	query := fmt.Sprintf("select %s from %s order by created_at desc limit ? offset ?", paymentConfigRows, m.paymentConfigsTable)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, limit, offset)
	return resp, err
}

func (m *defaultPaymentModel) CountPaymentConfigs(ctx context.Context) (int64, error) {
	var count int64
	query := fmt.Sprintf("select count(*) from %s", m.paymentConfigsTable)
	err := m.conn.QueryRowCtx(ctx, &count, query)
	return count, err
}

// 字段定义
const (
	paymentOrderRows                 = "id, payment_id, order_id, out_trade_no, user_id, amount, subject, body, status, trade_no, trade_status, buyer_user_id, buyer_logon_id, receipt_amount, product_code, return_url, notify_url, timeout, client_ip, gmt_payment, gmt_close, created_at, updated_at, deleted_at"
	paymentOrderRowsWithPlaceHolder  = "order_id = ?, out_trade_no = ?, user_id = ?, amount = ?, subject = ?, body = ?, status = ?, trade_no = ?, trade_status = ?, buyer_user_id = ?, buyer_logon_id = ?, receipt_amount = ?, product_code = ?, return_url = ?, notify_url = ?, timeout = ?, client_ip = ?, gmt_payment = ?, gmt_close = ?"
	paymentRefundRows                = "id, refund_id, payment_id, order_id, out_trade_no, out_request_no, user_id, refund_amount, refund_fee, refund_reason, status, refund_status, gmt_refund, created_at, updated_at, deleted_at"
	paymentRefundRowsWithPlaceHolder = "payment_id = ?, order_id = ?, out_trade_no = ?, out_request_no = ?, user_id = ?, refund_amount = ?, refund_fee = ?, refund_reason = ?, status = ?, refund_status = ?, gmt_refund = ?"
	paymentNotifyRows                = "id, notify_id, payment_id, notify_type, notify_data, sign, sign_type, verify_status, process_status, client_ip, error_message, processed_at, created_at, updated_at, deleted_at"
	paymentNotifyRowsWithPlaceHolder = "payment_id = ?, notify_type = ?, notify_data = ?, sign = ?, sign_type = ?, verify_status = ?, process_status = ?, client_ip = ?, error_message = ?, processed_at = ?"
	paymentConfigRows                = "id, app_id, app_name, app_private_key, alipay_public_key, gateway_url, is_prod, is_enabled, notify_url, return_url, created_at, updated_at, deleted_at"
	paymentConfigRowsWithPlaceHolder = "app_name = ?, app_private_key = ?, alipay_public_key = ?, gateway_url = ?, is_prod = ?, is_enabled = ?, notify_url = ?, return_url = ?"
)

var ErrNotFound = fmt.Errorf("not found")
