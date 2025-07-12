package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TxyChapterDataModel = (*customTxyChapterDataModel)(nil)

type (
	// TxyChapterDataModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTxyChapterDataModel.
	TxyChapterDataModel interface {
		txyChapterDataModel
		withSession(session sqlx.Session) TxyChapterDataModel
	}

	customTxyChapterDataModel struct {
		*defaultTxyChapterDataModel
	}
)

// NewTxyChapterDataModel returns a model for the database table.
func NewTxyChapterDataModel(conn sqlx.SqlConn) TxyChapterDataModel {
	return &customTxyChapterDataModel{
		defaultTxyChapterDataModel: newTxyChapterDataModel(conn),
	}
}

func (m *customTxyChapterDataModel) withSession(session sqlx.Session) TxyChapterDataModel {
	return NewTxyChapterDataModel(sqlx.NewSqlConnFromSession(session))
}
