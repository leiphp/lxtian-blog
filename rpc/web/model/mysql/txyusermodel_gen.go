// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.2

package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	txyUserFieldNames          = builder.RawFieldNames(&TxyUser{})
	txyUserRows                = strings.Join(txyUserFieldNames, ",")
	txyUserRowsExpectAutoSet   = strings.Join(stringx.Remove(txyUserFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	txyUserRowsWithPlaceHolder = strings.Join(stringx.Remove(txyUserFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	txyUserModel interface {
		Insert(ctx context.Context, data *TxyUser) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*TxyUser, error)
		FindOneByUsername(ctx context.Context, username sql.NullString) (*TxyUser, error)
		Update(ctx context.Context, data *TxyUser) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultTxyUserModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TxyUser struct {
		Id            uint64         `db:"id"`              // 主键id
		Uid           uint64         `db:"uid"`             // 关联的本站用户id
		Type          uint64         `db:"type"`            // 类型 1：QQ  2：新浪微博 3：微信 4：人人 5：开心网
		Nickname      string         `db:"nickname"`        // 第三方昵称
		HeadImg       string         `db:"head_img"`        // 头像
		Openid        string         `db:"openid"`          // 第三方用户id
		AccessToken   string         `db:"access_token"`    // access_token token
		Ctime         int64          `db:"ctime"`           // 创建时间
		Mtime         int64          `db:"mtime"`           // 修改时间
		LastLoginTime uint64         `db:"last_login_time"` // 最后登录时间
		LastLoginIp   string         `db:"last_login_ip"`   // 最后登录ip
		LoginTimes    uint64         `db:"login_times"`     // 登录次数
		Status        uint64         `db:"status"`          // 状态
		Email         string         `db:"email"`           // 邮箱
		Username      sql.NullString `db:"username"`        // 用户名
		Password      string         `db:"password"`        // 密码
		IsAdmin       uint64         `db:"is_admin"`        // 是否是admin
		Gold          uint64         `db:"gold"`
		Score         uint64         `db:"score"`
		Conscore      uint64         `db:"conscore"`
	}
)

func newTxyUserModel(conn sqlx.SqlConn) *defaultTxyUserModel {
	return &defaultTxyUserModel{
		conn:  conn,
		table: "`txy_user`",
	}
}

func (m *defaultTxyUserModel) Delete(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultTxyUserModel) FindOne(ctx context.Context, id uint64) (*TxyUser, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", txyUserRows, m.table)
	var resp TxyUser
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTxyUserModel) FindOneByUsername(ctx context.Context, username sql.NullString) (*TxyUser, error) {
	var resp TxyUser
	query := fmt.Sprintf("select %s from %s where `username` = ? limit 1", txyUserRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, username)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTxyUserModel) Insert(ctx context.Context, data *TxyUser) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, txyUserRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Uid, data.Type, data.Nickname, data.HeadImg, data.Openid, data.AccessToken, data.Ctime, data.Mtime, data.LastLoginTime, data.LastLoginIp, data.LoginTimes, data.Status, data.Email, data.Username, data.Password, data.IsAdmin, data.Gold, data.Score, data.Conscore)
	return ret, err
}

func (m *defaultTxyUserModel) Update(ctx context.Context, newData *TxyUser) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, txyUserRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.Uid, newData.Type, newData.Nickname, newData.HeadImg, newData.Openid, newData.AccessToken, newData.Ctime, newData.Mtime, newData.LastLoginTime, newData.LastLoginIp, newData.LoginTimes, newData.Status, newData.Email, newData.Username, newData.Password, newData.IsAdmin, newData.Gold, newData.Score, newData.Conscore, newData.Id)
	return err
}

func (m *defaultTxyUserModel) tableName() string {
	return m.table
}
