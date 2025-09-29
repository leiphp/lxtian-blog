package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LxtPaymentRefundsModel = (*customLxtPaymentRefundsModel)(nil)

type (
	// LxtPaymentRefundsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLxtPaymentRefundsModel.
	LxtPaymentRefundsModel interface {
		lxtPaymentRefundsModel
		withSession(session sqlx.Session) LxtPaymentRefundsModel
	}

	customLxtPaymentRefundsModel struct {
		*defaultLxtPaymentRefundsModel
	}
)

// NewLxtPaymentRefundsModel returns a model for the database table.
func NewLxtPaymentRefundsModel(conn sqlx.SqlConn) LxtPaymentRefundsModel {
	return &customLxtPaymentRefundsModel{
		defaultLxtPaymentRefundsModel: newLxtPaymentRefundsModel(conn),
	}
}

func (m *customLxtPaymentRefundsModel) withSession(session sqlx.Session) LxtPaymentRefundsModel {
	return NewLxtPaymentRefundsModel(sqlx.NewSqlConnFromSession(session))
}
