package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LxtUserMembershipTypesModel = (*customLxtUserMembershipTypesModel)(nil)

type (
	// LxtUserMembershipTypesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLxtUserMembershipTypesModel.
	LxtUserMembershipTypesModel interface {
		lxtUserMembershipTypesModel
		withSession(session sqlx.Session) LxtUserMembershipTypesModel
	}

	customLxtUserMembershipTypesModel struct {
		*defaultLxtUserMembershipTypesModel
	}
)

// NewLxtUserMembershipTypesModel returns a model for the database table.
func NewLxtUserMembershipTypesModel(conn sqlx.SqlConn) LxtUserMembershipTypesModel {
	return &customLxtUserMembershipTypesModel{
		defaultLxtUserMembershipTypesModel: newLxtUserMembershipTypesModel(conn),
	}
}

func (m *customLxtUserMembershipTypesModel) withSession(session sqlx.Session) LxtUserMembershipTypesModel {
	return NewLxtUserMembershipTypesModel(sqlx.NewSqlConnFromSession(session))
}
