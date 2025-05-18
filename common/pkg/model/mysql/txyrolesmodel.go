package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TxyRolesModel = (*customTxyRolesModel)(nil)

type (
	// TxyRolesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTxyRolesModel.
	TxyRolesModel interface {
		txyRolesModel
		withSession(session sqlx.Session) TxyRolesModel
	}

	customTxyRolesModel struct {
		*defaultTxyRolesModel
	}
)

// NewTxyRolesModel returns a model for the database table.
func NewTxyRolesModel(conn sqlx.SqlConn) TxyRolesModel {
	return &customTxyRolesModel{
		defaultTxyRolesModel: newTxyRolesModel(conn),
	}
}

func (m *customTxyRolesModel) withSession(session sqlx.Session) TxyRolesModel {
	return NewTxyRolesModel(sqlx.NewSqlConnFromSession(session))
}
