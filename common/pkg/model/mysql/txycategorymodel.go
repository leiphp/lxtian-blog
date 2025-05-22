package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TxyCategoryModel = (*customTxyCategoryModel)(nil)

type (
	// TxyCategoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTxyCategoryModel.
	TxyCategoryModel interface {
		txyCategoryModel
		withSession(session sqlx.Session) TxyCategoryModel
	}

	customTxyCategoryModel struct {
		*defaultTxyCategoryModel
	}
)

// NewTxyCategoryModel returns a model for the database table.
func NewTxyCategoryModel(conn sqlx.SqlConn) TxyCategoryModel {
	return &customTxyCategoryModel{
		defaultTxyCategoryModel: newTxyCategoryModel(conn),
	}
}

func (m *customTxyCategoryModel) withSession(session sqlx.Session) TxyCategoryModel {
	return NewTxyCategoryModel(sqlx.NewSqlConnFromSession(session))
}
