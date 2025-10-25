package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LxtUserMembershipRenewalsModel = (*customLxtUserMembershipRenewalsModel)(nil)

type (
	// LxtUserMembershipRenewalsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLxtUserMembershipRenewalsModel.
	LxtUserMembershipRenewalsModel interface {
		lxtUserMembershipRenewalsModel
		withSession(session sqlx.Session) LxtUserMembershipRenewalsModel
	}

	customLxtUserMembershipRenewalsModel struct {
		*defaultLxtUserMembershipRenewalsModel
	}
)

// NewLxtUserMembershipRenewalsModel returns a model for the database table.
func NewLxtUserMembershipRenewalsModel(conn sqlx.SqlConn) LxtUserMembershipRenewalsModel {
	return &customLxtUserMembershipRenewalsModel{
		defaultLxtUserMembershipRenewalsModel: newLxtUserMembershipRenewalsModel(conn),
	}
}

func (m *customLxtUserMembershipRenewalsModel) withSession(session sqlx.Session) LxtUserMembershipRenewalsModel {
	return NewLxtUserMembershipRenewalsModel(sqlx.NewSqlConnFromSession(session))
}
