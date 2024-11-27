package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TxyChatModel = (*customTxyChatModel)(nil)

type (
	// TxyChatModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTxyChatModel.
	TxyChatModel interface {
		txyChatModel
		withSession(session sqlx.Session) TxyChatModel
	}

	customTxyChatModel struct {
		*defaultTxyChatModel
	}
)

// NewTxyChatModel returns a model for the database table.
func NewTxyChatModel(conn sqlx.SqlConn) TxyChatModel {
	return &customTxyChatModel{
		defaultTxyChatModel: newTxyChatModel(conn),
	}
}

func (m *customTxyChatModel) withSession(session sqlx.Session) TxyChatModel {
	return NewTxyChatModel(sqlx.NewSqlConnFromSession(session))
}
