package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TxyColumnModel = (*customTxyColumnModel)(nil)

type (
	// TxyColumnModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTxyColumnModel.
	TxyColumnModel interface {
		txyColumnModel
		withSession(session sqlx.Session) TxyColumnModel
	}

	customTxyColumnModel struct {
		*defaultTxyColumnModel
	}
)

// NewTxyColumnModel returns a model for the database table.
func NewTxyColumnModel(conn sqlx.SqlConn) TxyColumnModel {
	return &customTxyColumnModel{
		defaultTxyColumnModel: newTxyColumnModel(conn),
	}
}

func (m *customTxyColumnModel) withSession(session sqlx.Session) TxyColumnModel {
	return NewTxyColumnModel(sqlx.NewSqlConnFromSession(session))
}
