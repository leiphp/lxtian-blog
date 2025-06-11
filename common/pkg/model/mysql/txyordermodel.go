package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TxyOrderModel = (*customTxyOrderModel)(nil)

type (
	// TxyOrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTxyOrderModel.
	TxyOrderModel interface {
		txyOrderModel
		withSession(session sqlx.Session) TxyOrderModel
	}

	customTxyOrderModel struct {
		*defaultTxyOrderModel
	}
)

// NewTxyOrderModel returns a model for the database table.
func NewTxyOrderModel(conn sqlx.SqlConn) TxyOrderModel {
	return &customTxyOrderModel{
		defaultTxyOrderModel: newTxyOrderModel(conn),
	}
}

func (m *customTxyOrderModel) withSession(session sqlx.Session) TxyOrderModel {
	return NewTxyOrderModel(sqlx.NewSqlConnFromSession(session))
}
