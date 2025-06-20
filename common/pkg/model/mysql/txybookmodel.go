package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TxyBookModel = (*customTxyBookModel)(nil)

type (
	// TxyBookModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTxyBookModel.
	TxyBookModel interface {
		txyBookModel
		withSession(session sqlx.Session) TxyBookModel
	}

	customTxyBookModel struct {
		*defaultTxyBookModel
	}
)

// NewTxyBookModel returns a model for the database table.
func NewTxyBookModel(conn sqlx.SqlConn) TxyBookModel {
	return &customTxyBookModel{
		defaultTxyBookModel: newTxyBookModel(conn),
	}
}

func (m *customTxyBookModel) withSession(session sqlx.Session) TxyBookModel {
	return NewTxyBookModel(sqlx.NewSqlConnFromSession(session))
}
