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
	txyCategoryFieldNames          = builder.RawFieldNames(&TxyCategory{})
	txyCategoryRows                = strings.Join(txyCategoryFieldNames, ",")
	txyCategoryRowsExpectAutoSet   = strings.Join(stringx.Remove(txyCategoryFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	txyCategoryRowsWithPlaceHolder = strings.Join(stringx.Remove(txyCategoryFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	txyCategoryModel interface {
		Insert(ctx context.Context, data *TxyCategory) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*TxyCategory, error)
		Update(ctx context.Context, data *TxyCategory) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultTxyCategoryModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TxyCategory struct {
		Id          uint64         `db:"id"`   // 分类主键id
		Name        string         `db:"name"` // 名称
		Seoname     sql.NullString `db:"seoname"`
		Keywords    string         `db:"keywords"`    // 关键词
		Description string         `db:"description"` // 描述
		Sort        uint64         `db:"sort"`        // 排序
		Pid         uint64         `db:"pid"`         // 父级栏目id
		Status      int64          `db:"status"`
		CreatedAt   sql.NullTime   `db:"created_at"` // 创建时间
		UpdatedAt   sql.NullTime   `db:"updated_at"`
		DeletedAt   sql.NullTime   `db:"deleted_at"`
	}
)

func newTxyCategoryModel(conn sqlx.SqlConn) *defaultTxyCategoryModel {
	return &defaultTxyCategoryModel{
		conn:  conn,
		table: "`txy_category`",
	}
}

func (m *defaultTxyCategoryModel) Delete(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultTxyCategoryModel) FindOne(ctx context.Context, id uint64) (*TxyCategory, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", txyCategoryRows, m.table)
	var resp TxyCategory
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

func (m *defaultTxyCategoryModel) Insert(ctx context.Context, data *TxyCategory) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table, txyCategoryRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Name, data.Seoname, data.Keywords, data.Description, data.Sort, data.Pid, data.Status, data.DeletedAt)
	return ret, err
}

func (m *defaultTxyCategoryModel) Update(ctx context.Context, data *TxyCategory) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, txyCategoryRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Name, data.Seoname, data.Keywords, data.Description, data.Sort, data.Pid, data.Status, data.DeletedAt, data.Id)
	return err
}

func (m *defaultTxyCategoryModel) tableName() string {
	return m.table
}
