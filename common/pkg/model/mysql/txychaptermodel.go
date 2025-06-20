package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TxyChapterModel = (*customTxyChapterModel)(nil)

type (
	// TxyChapterModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTxyChapterModel.
	TxyChapterModel interface {
		txyChapterModel
		withSession(session sqlx.Session) TxyChapterModel
	}

	customTxyChapterModel struct {
		*defaultTxyChapterModel
	}
)

// NewTxyChapterModel returns a model for the database table.
func NewTxyChapterModel(conn sqlx.SqlConn) TxyChapterModel {
	return &customTxyChapterModel{
		defaultTxyChapterModel: newTxyChapterModel(conn),
	}
}

func (m *customTxyChapterModel) withSession(session sqlx.Session) TxyChapterModel {
	return NewTxyChapterModel(sqlx.NewSqlConnFromSession(session))
}
