package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TxyTagModel = (*customTxyTagModel)(nil)

type (
	// TxyTagModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTxyTagModel.
	TxyTagModel interface {
		txyTagModel
		withSession(session sqlx.Session) TxyTagModel
	}

	customTxyTagModel struct {
		*defaultTxyTagModel
	}
)

// NewTxyTagModel returns a model for the database table.
func NewTxyTagModel(conn sqlx.SqlConn) TxyTagModel {
	return &customTxyTagModel{
		defaultTxyTagModel: newTxyTagModel(conn),
	}
}

func (m *customTxyTagModel) withSession(session sqlx.Session) TxyTagModel {
	return NewTxyTagModel(sqlx.NewSqlConnFromSession(session))
}
