package repository

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// baseRepository 基础仓储实现
type baseRepository[T any] struct {
	db *gorm.DB
}

// NewBaseRepository 创建基础仓储
func NewBaseRepository[T any](db *gorm.DB) *baseRepository[T] {
	return &baseRepository[T]{db: db}
}

// Create 创建单条记录
func (r *baseRepository[T]) Create(ctx context.Context, entity *T) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("failed to create entity: %w", err)
	}
	return nil
}

// CreateBatch 批量创建记录
func (r *baseRepository[T]) CreateBatch(ctx context.Context, entities []*T) error {
	if len(entities) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).CreateInBatches(entities, 100).Error; err != nil {
		return fmt.Errorf("failed to create entities batch: %w", err)
	}
	return nil
}

// GetByID 根据ID获取记录
func (r *baseRepository[T]) GetByID(ctx context.Context, id uint64) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("entity not found with id %d", id)
		}
		return nil, fmt.Errorf("failed to get entity by id: %w", err)
	}
	return &entity, nil
}

// GetByCondition 根据条件获取单条记录
func (r *baseRepository[T]) GetByCondition(ctx context.Context, condition map[string]interface{}) (*T, error) {
	var entity T
	query := r.db.WithContext(ctx)

	for key, value := range condition {
		query = query.Where(key, value)
	}

	if err := query.First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("entity not found with condition: %v", condition)
		}
		return nil, fmt.Errorf("failed to get entity by condition: %w", err)
	}
	return &entity, nil
}

// GetList 根据条件获取列表（支持分页）
func (r *baseRepository[T]) GetList(ctx context.Context, condition map[string]interface{}, page, pageSize int, orderBy string, keywords string, fields ...string) ([]*T, int64, error) {
	var entities []*T
	var total int64

	query := r.db.WithContext(ctx).Debug()

	// 构建查询条件
	for key, value := range condition {
		query = query.Where(key, value)
	}

	// 移除兼容模式：不再从 condition 中读取 keywords

	if keywords != "" && len(fields) > 0 {
		kw := "%" + keywords + "%"
		parts := make([]string, 0, len(fields))
		vals := make([]interface{}, 0, len(fields))
		for _, f := range fields {
			f = strings.TrimSpace(f)
			if f == "" {
				continue
			}
			// 使用反引号包裹字段名，避免 MySQL 保留关键字冲突
			parts = append(parts, fmt.Sprintf("`%s` LIKE ?", f))
			vals = append(vals, kw)
		}
		if len(parts) > 0 {
			where := strings.Join(parts, " OR ")
			query = query.Where(where, vals...)
		}
	}

	// 获取总数
	if err := query.Model(new(T)).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count entities: %w", err)
	}

	// 分页查询
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	orderClause := orderBy
	if strings.TrimSpace(orderClause) == "" {
		orderClause = "id desc"
	}
	if err := query.Order(orderClause).Find(&entities).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get entities list: %w", err)
	}

	return entities, total, nil
}

// Update 更新记录
func (r *baseRepository[T]) Update(ctx context.Context, entity *T) error {
	if err := r.db.WithContext(ctx).Debug().Save(entity).Error; err != nil {
		return fmt.Errorf("failed to update entity: %w", err)
	}
	return nil
}

// UpdateByCondition 根据条件更新记录
func (r *baseRepository[T]) UpdateByCondition(ctx context.Context, condition map[string]interface{}, updates map[string]interface{}) error {
	query := r.db.WithContext(ctx).Model(new(T))

	// 构建查询条件
	for key, value := range condition {
		query = query.Where(key, value)
	}

	if err := query.Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update entities by condition: %w", err)
	}
	return nil
}

// Delete 根据ID删除记录（软删除，如果模型有DeletedAt字段）
func (r *baseRepository[T]) Delete(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).Delete(new(T), id).Error; err != nil {
		return fmt.Errorf("failed to delete entity by id: %w", err)
	}
	return nil
}

// DeleteByCondition 根据条件删除记录（软删除，如果模型有DeletedAt字段）
func (r *baseRepository[T]) DeleteByCondition(ctx context.Context, condition map[string]interface{}) error {
	query := r.db.WithContext(ctx)

	// 构建查询条件
	for key, value := range condition {
		query = query.Where(key, value)
	}

	if err := query.Delete(new(T)).Error; err != nil {
		return fmt.Errorf("failed to delete entities by condition: %w", err)
	}
	return nil
}

// ForceDelete 根据ID物理删除记录（永久删除，忽略DeletedAt字段）
func (r *baseRepository[T]) ForceDelete(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).Unscoped().Delete(new(T), id).Error; err != nil {
		return fmt.Errorf("failed to force delete entity by id: %w", err)
	}
	return nil
}

// ForceDeleteByCondition 根据条件物理删除记录（永久删除，忽略DeletedAt字段）
func (r *baseRepository[T]) ForceDeleteByCondition(ctx context.Context, condition map[string]interface{}) error {
	query := r.db.WithContext(ctx).Unscoped()

	// 构建查询条件
	for key, value := range condition {
		query = query.Where(key, value)
	}

	if err := query.Delete(new(T)).Error; err != nil {
		return fmt.Errorf("failed to force delete entities by condition: %w", err)
	}
	return nil
}

// Count 统计记录数
func (r *baseRepository[T]) Count(ctx context.Context, condition map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(new(T))

	// 构建查询条件
	for key, value := range condition {
		query = query.Where(key, value)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count entities: %w", err)
	}
	return count, nil
}

// Exists 检查记录是否存在
func (r *baseRepository[T]) Exists(ctx context.Context, condition map[string]interface{}) (bool, error) {
	count, err := r.Count(ctx, condition)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// WithTransaction 事务操作
func (r *baseRepository[T]) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建新的上下文，包含事务信息
		txCtx := context.WithValue(ctx, "tx", tx)
		return fn(txCtx)
	})
}

// GetDB 获取数据库连接（用于复杂查询）
func (r *baseRepository[T]) GetDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok {
		return tx
	}
	return r.db.WithContext(ctx)
}

// QueryBuilder 查询构建器
type QueryBuilder struct {
	db *gorm.DB
}

// NewQueryBuilder 创建查询构建器
func NewQueryBuilder(db *gorm.DB) *QueryBuilder {
	return &QueryBuilder{db: db}
}

// Select 指定查询字段
func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder {
	qb.db = qb.db.Select(fields)
	return qb
}

// Where 添加查询条件
func (qb *QueryBuilder) Where(query interface{}, args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Where(query, args...)
	return qb
}

// WhereIn 添加IN查询条件
func (qb *QueryBuilder) WhereIn(field string, values []interface{}) *QueryBuilder {
	qb.db = qb.db.Where(fmt.Sprintf("%s IN ?", field), values)
	return qb
}

// WhereBetween 添加BETWEEN查询条件
func (qb *QueryBuilder) WhereBetween(field string, start, end interface{}) *QueryBuilder {
	qb.db = qb.db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", field), start, end)
	return qb
}

// Order 添加排序
func (qb *QueryBuilder) Order(value interface{}) *QueryBuilder {
	qb.db = qb.db.Order(value)
	return qb
}

// Group 添加分组
func (qb *QueryBuilder) Group(name string) *QueryBuilder {
	qb.db = qb.db.Group(name)
	return qb
}

// Having 添加HAVING条件
func (qb *QueryBuilder) Having(query interface{}, args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Having(query, args...)
	return qb
}

// Join 添加JOIN
func (qb *QueryBuilder) Join(query string, args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Joins(query, args...)
	return qb
}

// Limit 添加限制
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.db = qb.db.Limit(limit)
	return qb
}

// Offset 添加偏移
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.db = qb.db.Offset(offset)
	return qb
}

// Page 添加分页
func (qb *QueryBuilder) Page(page, pageSize int) *QueryBuilder {
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		qb.db = qb.db.Offset(offset).Limit(pageSize)
	}
	return qb
}

// Execute 执行查询
func (qb *QueryBuilder) Execute(dest interface{}) error {
	return qb.db.Find(dest).Error
}

// ExecuteFirst 执行查询并返回第一条记录
func (qb *QueryBuilder) ExecuteFirst(dest interface{}) error {
	return qb.db.First(dest).Error
}

// Count 统计记录数
func (qb *QueryBuilder) Count(count *int64) error {
	return qb.db.Count(count).Error
}

// ExecuteRaw 执行原生SQL
func (qb *QueryBuilder) ExecuteRaw(sql string, dest interface{}, args ...interface{}) error {
	return qb.db.Raw(sql, args...).Scan(dest).Error
}

// GetDB 获取数据库连接
func (qb *QueryBuilder) GetDB() *gorm.DB {
	return qb.db
}
