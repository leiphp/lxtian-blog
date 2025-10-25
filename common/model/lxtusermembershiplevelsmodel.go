package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LxtUserMembershipLevelsModel = (*customLxtUserMembershipLevelsModel)(nil)

type (
	// LxtUserMembershipLevelsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLxtUserMembershipLevelsModel.
	LxtUserMembershipLevelsModel interface {
		lxtUserMembershipLevelsModel
		withSession(session sqlx.Session) LxtUserMembershipLevelsModel
	}

	customLxtUserMembershipLevelsModel struct {
		*defaultLxtUserMembershipLevelsModel
	}
)

// NewLxtUserMembershipLevelsModel returns a model for the database table.
func NewLxtUserMembershipLevelsModel(conn sqlx.SqlConn) LxtUserMembershipLevelsModel {
	return &customLxtUserMembershipLevelsModel{
		defaultLxtUserMembershipLevelsModel: newLxtUserMembershipLevelsModel(conn),
	}
}

func (m *customLxtUserMembershipLevelsModel) withSession(session sqlx.Session) LxtUserMembershipLevelsModel {
	return NewLxtUserMembershipLevelsModel(sqlx.NewSqlConnFromSession(session))
}
