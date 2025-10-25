package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LxtUserMembershipsModel = (*customLxtUserMembershipsModel)(nil)

type (
	// LxtUserMembershipsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLxtUserMembershipsModel.
	LxtUserMembershipsModel interface {
		lxtUserMembershipsModel
		withSession(session sqlx.Session) LxtUserMembershipsModel
	}

	customLxtUserMembershipsModel struct {
		*defaultLxtUserMembershipsModel
	}
)

// NewLxtUserMembershipsModel returns a model for the database table.
func NewLxtUserMembershipsModel(conn sqlx.SqlConn) LxtUserMembershipsModel {
	return &customLxtUserMembershipsModel{
		defaultLxtUserMembershipsModel: newLxtUserMembershipsModel(conn),
	}
}

func (m *customLxtUserMembershipsModel) withSession(session sqlx.Session) LxtUserMembershipsModel {
	return NewLxtUserMembershipsModel(sqlx.NewSqlConnFromSession(session))
}
