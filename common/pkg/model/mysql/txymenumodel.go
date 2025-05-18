package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TxyMenuModel = (*customTxyMenuModel)(nil)

type (
	// TxyMenuModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTxyMenuModel.
	TxyMenuModel interface {
		txyMenuModel
		withSession(session sqlx.Session) TxyMenuModel
	}

	customTxyMenuModel struct {
		*defaultTxyMenuModel
	}
)

// NewTxyMenuModel returns a model for the database table.
func NewTxyMenuModel(conn sqlx.SqlConn) TxyMenuModel {
	return &customTxyMenuModel{
		defaultTxyMenuModel: newTxyMenuModel(conn),
	}
}

func (m *customTxyMenuModel) withSession(session sqlx.Session) TxyMenuModel {
	return NewTxyMenuModel(sqlx.NewSqlConnFromSession(session))
}
