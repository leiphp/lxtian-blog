package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LxtUserMembershipPermissionsModel = (*customLxtUserMembershipPermissionsModel)(nil)

type (
	// LxtUserMembershipPermissionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLxtUserMembershipPermissionsModel.
	LxtUserMembershipPermissionsModel interface {
		lxtUserMembershipPermissionsModel
		withSession(session sqlx.Session) LxtUserMembershipPermissionsModel
	}

	customLxtUserMembershipPermissionsModel struct {
		*defaultLxtUserMembershipPermissionsModel
	}
)

// NewLxtUserMembershipPermissionsModel returns a model for the database table.
func NewLxtUserMembershipPermissionsModel(conn sqlx.SqlConn) LxtUserMembershipPermissionsModel {
	return &customLxtUserMembershipPermissionsModel{
		defaultLxtUserMembershipPermissionsModel: newLxtUserMembershipPermissionsModel(conn),
	}
}

func (m *customLxtUserMembershipPermissionsModel) withSession(session sqlx.Session) LxtUserMembershipPermissionsModel {
	return NewLxtUserMembershipPermissionsModel(sqlx.NewSqlConnFromSession(session))
}
