package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TxyPermissionsModel = (*customTxyPermissionsModel)(nil)

type (
	// TxyPermissionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTxyPermissionsModel.
	TxyPermissionsModel interface {
		txyPermissionsModel
		withSession(session sqlx.Session) TxyPermissionsModel
	}

	customTxyPermissionsModel struct {
		*defaultTxyPermissionsModel
	}
)

// NewTxyPermissionsModel returns a model for the database table.
func NewTxyPermissionsModel(conn sqlx.SqlConn) TxyPermissionsModel {
	return &customTxyPermissionsModel{
		defaultTxyPermissionsModel: newTxyPermissionsModel(conn),
	}
}

func (m *customTxyPermissionsModel) withSession(session sqlx.Session) TxyPermissionsModel {
	return NewTxyPermissionsModel(sqlx.NewSqlConnFromSession(session))
}
