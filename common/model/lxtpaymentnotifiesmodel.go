package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LxtPaymentNotifiesModel = (*customLxtPaymentNotifiesModel)(nil)

type (
	// LxtPaymentNotifiesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLxtPaymentNotifiesModel.
	LxtPaymentNotifiesModel interface {
		lxtPaymentNotifiesModel
		withSession(session sqlx.Session) LxtPaymentNotifiesModel
	}

	customLxtPaymentNotifiesModel struct {
		*defaultLxtPaymentNotifiesModel
	}
)

// NewLxtPaymentNotifiesModel returns a model for the database table.
func NewLxtPaymentNotifiesModel(conn sqlx.SqlConn) LxtPaymentNotifiesModel {
	return &customLxtPaymentNotifiesModel{
		defaultLxtPaymentNotifiesModel: newLxtPaymentNotifiesModel(conn),
	}
}

func (m *customLxtPaymentNotifiesModel) withSession(session sqlx.Session) LxtPaymentNotifiesModel {
	return NewLxtPaymentNotifiesModel(sqlx.NewSqlConnFromSession(session))
}
