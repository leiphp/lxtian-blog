package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TxyArticleModel = (*customTxyArticleModel)(nil)

type (
	// TxyArticleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTxyArticleModel.
	TxyArticleModel interface {
		txyArticleModel
		withSession(session sqlx.Session) TxyArticleModel
	}

	customTxyArticleModel struct {
		*defaultTxyArticleModel
	}
)

// NewTxyArticleModel returns a model for the database table.
func NewTxyArticleModel(conn sqlx.SqlConn) TxyArticleModel {
	return &customTxyArticleModel{
		defaultTxyArticleModel: newTxyArticleModel(conn),
	}
}

func (m *customTxyArticleModel) withSession(session sqlx.Session) TxyArticleModel {
	return NewTxyArticleModel(sqlx.NewSqlConnFromSession(session))
}
