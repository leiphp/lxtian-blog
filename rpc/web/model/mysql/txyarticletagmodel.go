package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TxyArticleTagModel = (*customTxyArticleTagModel)(nil)

type (
	// TxyArticleTagModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTxyArticleTagModel.
	TxyArticleTagModel interface {
		txyArticleTagModel
		withSession(session sqlx.Session) TxyArticleTagModel
	}

	customTxyArticleTagModel struct {
		*defaultTxyArticleTagModel
	}
)

// NewTxyArticleTagModel returns a model for the database table.
func NewTxyArticleTagModel(conn sqlx.SqlConn) TxyArticleTagModel {
	return &customTxyArticleTagModel{
		defaultTxyArticleTagModel: newTxyArticleTagModel(conn),
	}
}

func (m *customTxyArticleTagModel) withSession(session sqlx.Session) TxyArticleTagModel {
	return NewTxyArticleTagModel(sqlx.NewSqlConnFromSession(session))
}
