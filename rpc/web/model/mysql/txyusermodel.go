package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TxyUserModel = (*customTxyUserModel)(nil)

type (
	// TxyUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTxyUserModel.
	TxyUserModel interface {
		txyUserModel
		withSession(session sqlx.Session) TxyUserModel
	}

	customTxyUserModel struct {
		*defaultTxyUserModel
	}
)

// NewTxyUserModel returns a model for the database table.
func NewTxyUserModel(conn sqlx.SqlConn) TxyUserModel {
	return &customTxyUserModel{
		defaultTxyUserModel: newTxyUserModel(conn),
	}
}

func (m *customTxyUserModel) withSession(session sqlx.Session) TxyUserModel {
	return NewTxyUserModel(sqlx.NewSqlConnFromSession(session))
}
