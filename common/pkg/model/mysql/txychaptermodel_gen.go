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
	txyChapterFieldNames          = builder.RawFieldNames(&TxyChapter{})
	txyChapterRows                = strings.Join(txyChapterFieldNames, ",")
	txyChapterRowsExpectAutoSet   = strings.Join(stringx.Remove(txyChapterFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	txyChapterRowsWithPlaceHolder = strings.Join(stringx.Remove(txyChapterFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	txyChapterModel interface {
		Insert(ctx context.Context, data *TxyChapter) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*TxyChapter, error)
		Update(ctx context.Context, data *TxyChapter) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultTxyChapterModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TxyChapter struct {
		Id          uint64         `db:"id"`          // 主键
		BookId      uint64         `db:"book_id"`     // 所属书籍 ID
		Title       string         `db:"title"`       // 标题
		Slug        string         `db:"slug"`        // URL 唯一标识
		Sort        int64          `db:"sort"`        // 排序值
		IsGroup     int64          `db:"is_group"`    // 是否组
		IsOpen      int64          `db:"is_open"`     // 是否展开
		Content     sql.NullString `db:"content"`     // 内容
		Description string         `db:"description"` // 描述
		CreatedAt   sql.NullTime   `db:"created_at"`  // 添加时间
		UpdatedAt   sql.NullTime   `db:"updated_at"`  // 修改时间
		DeletedAt   sql.NullTime   `db:"deleted_at"`  // 删除时间
	}
)

func newTxyChapterModel(conn sqlx.SqlConn) *defaultTxyChapterModel {
	return &defaultTxyChapterModel{
		conn:  conn,
		table: "`txy_chapter`",
	}
}

func (m *defaultTxyChapterModel) Delete(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultTxyChapterModel) FindOne(ctx context.Context, id uint64) (*TxyChapter, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", txyChapterRows, m.table)
	var resp TxyChapter
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

func (m *defaultTxyChapterModel) Insert(ctx context.Context, data *TxyChapter) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, txyChapterRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.BookId, data.Title, data.Slug, data.Sort, data.IsGroup, data.IsOpen, data.Content, data.Description, data.DeletedAt)
	return ret, err
}

func (m *defaultTxyChapterModel) Update(ctx context.Context, data *TxyChapter) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, txyChapterRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.BookId, data.Title, data.Slug, data.Sort, data.IsGroup, data.IsOpen, data.Content, data.Description, data.DeletedAt, data.Id)
	return err
}

func (m *defaultTxyChapterModel) tableName() string {
	return m.table
}
