package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TxyCommentModel = (*customTxyCommentModel)(nil)

type (
	// TxyCommentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTxyCommentModel.
	TxyCommentModel interface {
		txyCommentModel
		withSession(session sqlx.Session) TxyCommentModel
	}

	customTxyCommentModel struct {
		*defaultTxyCommentModel
	}
)

// NewTxyCommentModel returns a model for the database table.
func NewTxyCommentModel(conn sqlx.SqlConn) TxyCommentModel {
	return &customTxyCommentModel{
		defaultTxyCommentModel: newTxyCommentModel(conn),
	}
}

func (m *customTxyCommentModel) withSession(session sqlx.Session) TxyCommentModel {
	return NewTxyCommentModel(sqlx.NewSqlConnFromSession(session))
}
