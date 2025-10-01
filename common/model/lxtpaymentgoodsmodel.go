package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LxtPaymentGoodsModel = (*customLxtPaymentGoodsModel)(nil)

type (
	// LxtPaymentGoodsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLxtPaymentGoodsModel.
	LxtPaymentGoodsModel interface {
		lxtPaymentGoodsModel
		withSession(session sqlx.Session) LxtPaymentGoodsModel
	}

	customLxtPaymentGoodsModel struct {
		*defaultLxtPaymentGoodsModel
	}
)

// NewLxtPaymentGoodsModel returns a model for the database table.
func NewLxtPaymentGoodsModel(conn sqlx.SqlConn) LxtPaymentGoodsModel {
	return &customLxtPaymentGoodsModel{
		defaultLxtPaymentGoodsModel: newLxtPaymentGoodsModel(conn),
	}
}

func (m *customLxtPaymentGoodsModel) withSession(session sqlx.Session) LxtPaymentGoodsModel {
	return NewLxtPaymentGoodsModel(sqlx.NewSqlConnFromSession(session))
}
