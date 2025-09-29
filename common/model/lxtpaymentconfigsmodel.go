package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LxtPaymentConfigsModel = (*customLxtPaymentConfigsModel)(nil)

type (
	// LxtPaymentConfigsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLxtPaymentConfigsModel.
	LxtPaymentConfigsModel interface {
		lxtPaymentConfigsModel
		withSession(session sqlx.Session) LxtPaymentConfigsModel
	}

	customLxtPaymentConfigsModel struct {
		*defaultLxtPaymentConfigsModel
	}
)

// NewLxtPaymentConfigsModel returns a model for the database table.
func NewLxtPaymentConfigsModel(conn sqlx.SqlConn) LxtPaymentConfigsModel {
	return &customLxtPaymentConfigsModel{
		defaultLxtPaymentConfigsModel: newLxtPaymentConfigsModel(conn),
	}
}

func (m *customLxtPaymentConfigsModel) withSession(session sqlx.Session) LxtPaymentConfigsModel {
	return NewLxtPaymentConfigsModel(sqlx.NewSqlConnFromSession(session))
}
