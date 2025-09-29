package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LxtPaymentOrdersModel = (*customLxtPaymentOrdersModel)(nil)

type (
	// LxtPaymentOrdersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLxtPaymentOrdersModel.
	LxtPaymentOrdersModel interface {
		lxtPaymentOrdersModel
		withSession(session sqlx.Session) LxtPaymentOrdersModel
	}

	customLxtPaymentOrdersModel struct {
		*defaultLxtPaymentOrdersModel
	}
)

// NewLxtPaymentOrdersModel returns a model for the database table.
func NewLxtPaymentOrdersModel(conn sqlx.SqlConn) LxtPaymentOrdersModel {
	return &customLxtPaymentOrdersModel{
		defaultLxtPaymentOrdersModel: newLxtPaymentOrdersModel(conn),
	}
}

func (m *customLxtPaymentOrdersModel) withSession(session sqlx.Session) LxtPaymentOrdersModel {
	return NewLxtPaymentOrdersModel(sqlx.NewSqlConnFromSession(session))
}
