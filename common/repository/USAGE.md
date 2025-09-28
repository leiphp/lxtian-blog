# Repository 通用数据库操作封装 - 使用指南

## 🎯 项目概述

这是一个基于 GORM 的通用数据库操作封装，支持事务操作，可以前后端复用。主要特点：

- ✅ **通用CRUD操作**: 提供基础的增删改查方法
- ✅ **事务支持**: 完整的事务管理，支持嵌套事务
- ✅ **查询构建器**: 灵活的查询构建器，支持复杂查询
- ✅ **模块化设计**: 按服务模块分离，易于维护和扩展
- ✅ **类型安全**: 使用 Go 泛型，提供类型安全
- ✅ **批量操作**: 支持批量插入、更新、删除

## 📁 目录结构

```
common/repository/
├── base.go              # 基础仓储实现
├── transaction.go       # 事务管理器
├── manager.go          # 仓储管理器
├── interfaces.go       # 接口定义
├── simple_example.go   # 使用示例
├── README.md           # 详细文档
├── USAGE.md            # 使用指南（本文件）
├── user/               # 用户模块仓储
│   ├── txyuser_repo.go      # TxyUser表仓储
│   ├── txyroles_repo.go     # TxyRoles表仓储
│   └── txypermissions_repo.go # TxyPermissions表仓储
├── payment/            # 支付模块仓储
│   ├── paymentorder_repo.go # PaymentOrder表仓储
│   └── txyorder_repo.go     # TxyOrder表仓储
└── web/                # Web模块仓储
    ├── txyarticle_repo.go    # TxyArticle表仓储
    ├── txycategory_repo.go   # TxyCategory表仓储
    ├── txytag_repo.go        # TxyTag表仓储
    └── txycomment_repo.go    # TxyComment表仓储
```

## 🚀 快速开始

### 1. 基本初始化

```go
package main

import (
    "context"
    "lxtian-blog/common/repository"
)

func main() {
    // 使用配置初始化
    mysqlDataSource := "root:password@tcp(127.0.0.1:3306)/lxtian_blog?charset=utf8mb4&parseTime=True&loc=Local"
    repoManager := repository.NewRepositoryManagerWithConfig(mysqlDataSource)
    
    // 或者使用现有数据库连接
    // repoManager := repository.NewRepositoryManager(db)
}
```

### 2. 基础CRUD操作

```go
// 创建基础仓储
userRepo := repository.NewBaseRepository[mysql.TxyUser](repoManager.GetDB())
ctx := context.Background()

// 创建用户
user := &mysql.TxyUser{
    Nickname: "测试用户",
    // ... 其他字段
}
err := userRepo.Create(ctx, user)

// 查询用户
user, err := userRepo.GetByID(ctx, 1)

// 更新用户
user.Nickname = "更新后的昵称"
err := userRepo.Update(ctx, user)

// 删除用户
err := userRepo.Delete(ctx, 1)

// 条件查询
user, err := userRepo.GetByCondition(ctx, map[string]interface{}{
    "username": "testuser",
})

// 分页查询
users, total, err := userRepo.GetList(ctx, map[string]interface{}{
    "status": 1,
}, 1, 10) // 第1页，每页10条
```

### 2.1 使用具体的表仓储

```go
// 用户相关仓储
userRepo := user.NewTxyUserRepository(repoManager.GetDB())
rolesRepo := user.NewTxyRolesRepository(repoManager.GetDB())
permsRepo := user.NewTxyPermissionsRepository(repoManager.GetDB())

// Web相关仓储
articleRepo := web.NewTxyArticleRepository(repoManager.GetDB())
categoryRepo := web.NewTxyCategoryRepository(repoManager.GetDB())
tagRepo := web.NewTxyTagRepository(repoManager.GetDB())
commentRepo := web.NewTxyCommentRepository(repoManager.GetDB())

// 支付相关仓储
paymentRepo := payment.NewPaymentOrderRepository(repoManager.GetDB())
orderRepo := payment.NewTxyOrderRepository(repoManager.GetDB())

// 使用具体的仓储方法
user, err := userRepo.GetByUid(ctx, 12345)
article, err := articleRepo.GetByTitle(ctx, "文章标题")
payment, err := paymentRepo.GetByPaymentId(ctx, "PAY_123456")
```

### 3. 事务操作

```go
err := repoManager.WithTransaction(ctx, func(txCtx context.Context) error {
    // 在事务中执行多个操作
    user := &mysql.TxyUser{Nickname: "事务用户"}
    if err := userRepo.Create(txCtx, user); err != nil {
        return err
    }
    
    article := &mysql.TxyArticle{Title: "事务文章"}
    if err := articleRepo.Create(txCtx, article); err != nil {
        return err
    }
    
    return nil
})
```

### 4. 查询构建器

```go
db := repoManager.GetDB()
qb := repository.NewQueryBuilder(db)

// 构建复杂查询
var users []mysql.TxyUser
err := qb.Select("id", "nickname", "head_img").
    Where("status = ?", 1).
    Where("last_login_time > ?", 1640995200).
    Order("last_login_time DESC").
    Limit(10).
    Execute(&users)
```

## 📋 主要组件说明

### 1. RepositoryManager (仓储管理器)

负责管理所有仓储实例和事务操作。

```go
type RepositoryManager struct {
    db *gorm.DB
    TransactionManager *TransactionManager
}

// 主要方法
func (rm *RepositoryManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
func (rm *RepositoryManager) BatchExecute(ctx context.Context, operations []TransactionFunc) error
func (rm *RepositoryManager) GetDB() *gorm.DB
```

### 2. BaseRepository (基础仓储)

提供通用的CRUD操作方法。

```go
type BaseRepository[T any] interface {
    Create(ctx context.Context, entity *T) error
    CreateBatch(ctx context.Context, entities []*T) error
    GetByID(ctx context.Context, id uint64) (*T, error)
    GetByCondition(ctx context.Context, condition map[string]interface{}) (*T, error)
    GetList(ctx context.Context, condition map[string]interface{}, page, pageSize int) ([]*T, int64, error)
    Update(ctx context.Context, entity *T) error
    UpdateByCondition(ctx context.Context, condition map[string]interface{}, updates map[string]interface{}) error
    Delete(ctx context.Context, id uint64) error
    DeleteByCondition(ctx context.Context, condition map[string]interface{}) error
    Count(ctx context.Context, condition map[string]interface{}) (int64, error)
    Exists(ctx context.Context, condition map[string]interface{}) (bool, error)
    WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
```

### 3. QueryBuilder (查询构建器)

提供灵活的查询构建能力。

```go
type QueryBuilder struct {
    db *gorm.DB
}

// 主要方法
func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder
func (qb *QueryBuilder) Where(query interface{}, args ...interface{}) *QueryBuilder
func (qb *QueryBuilder) WhereIn(field string, values []interface{}) *QueryBuilder
func (qb *QueryBuilder) WhereBetween(field string, start, end interface{}) *QueryBuilder
func (qb *QueryBuilder) Order(value interface{}) *QueryBuilder
func (qb *QueryBuilder) Group(name string) *QueryBuilder
func (qb *QueryBuilder) Having(query interface{}, args ...interface{}) *QueryBuilder
func (qb *QueryBuilder) Join(query string, args ...interface{}) *QueryBuilder
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder
func (qb *QueryBuilder) Page(page, pageSize int) *QueryBuilder
func (qb *QueryBuilder) Execute(dest interface{}) error
func (qb *QueryBuilder) ExecuteFirst(dest interface{}) error
func (qb *QueryBuilder) Count(count *int64) error
func (qb *QueryBuilder) ExecuteRaw(sql string, dest interface{}, args ...interface{}) error
```

### 4. TransactionManager (事务管理器)

管理数据库事务操作。

```go
type TransactionManager struct {
    db *gorm.DB
}

// 主要方法
func (tm *TransactionManager) ExecuteInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
func (tm *TransactionManager) ExecuteWithOptions(ctx context.Context, fn TransactionFunc, options *TransactionOptions) error
func (tm *TransactionManager) BatchExecute(ctx context.Context, operations []TransactionFunc) error
```

## 🎯 实际使用场景

### 1. 在 RPC 服务中使用

```go
// rpc/user/internal/svc/servicecontext.go
type ServiceContext struct {
    Config       config.Config
    DB           sqlx.SqlConn
    Rds          *redis.Redis
    UserModel    model.UserModel
    UserRepo     repository.UserRepository  // 新增
}

func NewServiceContext(c config.Config) *ServiceContext {
    // ... 现有代码
    
    // 初始化仓储
    repoManager := repository.NewRepositoryManagerWithConfig(mysqlDataSource)
    userRepo := repository.NewUserRepository(repoManager.GetDB())
    
    return &ServiceContext{
        // ... 现有字段
        UserRepo: userRepo,
    }
}
```

### 2. 在 Gateway 服务中使用

```go
// gateway/internal/svc/servicecontext.go
type ServiceContext struct {
    Config       config.Config
    UserRpc      user.User
    PaymentRpc   payment.Payment
    WebRpc       web.Web
    UserRepo     repository.UserRepository     // 新增
    PaymentRepo  repository.PaymentOrderRepository // 新增
    ArticleRepo  repository.ArticleRepository  // 新增
}

func NewServiceContext(c config.Config) *ServiceContext {
    // ... 现有代码
    
    // 初始化仓储
    repoManager := repository.NewRepositoryManagerWithConfig(mysqlDataSource)
    
    return &ServiceContext{
        // ... 现有字段
        UserRepo:    repository.NewUserRepository(repoManager.GetDB()),
        PaymentRepo: repository.NewPaymentOrderRepository(repoManager.GetDB()),
        ArticleRepo: repository.NewArticleRepository(repoManager.GetDB()),
    }
}
```

### 3. 在业务逻辑中使用

```go
// gateway/internal/logic/user/getuserlogic.go
func (l *GetUserLogic) GetUser(req *types.GetUserReq) (*types.GetUserResp, error) {
    // 使用仓储进行数据操作
    user, err := l.svcCtx.UserRepo.GetByID(l.ctx, req.UserId)
    if err != nil {
        return nil, err
    }
    
    return &types.GetUserResp{
        Id:       user.Id,
        Nickname: user.Nickname,
        HeadImg:  user.HeadImg,
    }, nil
}
```

## 🔧 配置选项

### 1. 基础配置

```go
config := &repository.RepositoryConfig{
    MysqlDataSource: "root:password@tcp(127.0.0.1:3306)/lxtian_blog",
    MaxIdleConns:    10,    // 最大空闲连接数
    MaxOpenConns:    100,   // 最大打开连接数
    MaxLifetime:     3600,  // 连接最大生存时间（秒）
}

repoManager := repository.NewRepositoryManagerWithOptions(config)
```

### 2. 环境变量配置

```go
// 从环境变量读取配置
mysqlDataSource := os.Getenv("MYSQL_DATA_SOURCE")
if mysqlDataSource == "" {
    mysqlDataSource = "root:password@tcp(127.0.0.1:3306)/lxtian_blog?charset=utf8mb4&parseTime=True&loc=Local"
}

repoManager := repository.NewRepositoryManagerWithConfig(mysqlDataSource)
```

## 📝 最佳实践

### 1. 错误处理

```go
user, err := userRepo.GetByID(ctx, 1)
if err != nil {
    if strings.Contains(err.Error(), "not found") {
        return nil, errors.New("用户不存在")
    }
    return nil, fmt.Errorf("查询用户失败: %w", err)
}
```

### 2. 上下文传递

```go
// 在服务层传递上下文
func (s *UserService) GetUser(ctx context.Context, userID uint64) (*User, error) {
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return nil, err
    }
    return s.convertToUser(user), nil
}
```

### 3. 事务边界

```go
// 在业务逻辑层定义事务边界
func (s *PaymentService) CreateOrder(ctx context.Context, req *CreateOrderReq) error {
    return s.repoManager.WithTransaction(ctx, func(txCtx context.Context) error {
        // 创建业务订单
        order := &mysql.TxyOrder{...}
        if err := s.businessRepo.Create(txCtx, order); err != nil {
            return err
        }
        
        // 创建支付订单
        payment := &model.PaymentOrder{...}
        if err := s.paymentRepo.Create(txCtx, payment); err != nil {
            return err
        }
        
        return nil
    })
}
```

### 4. 性能优化

```go
// 使用批量操作
users := make([]*mysql.TxyUser, 0, 100)
for i := 0; i < 100; i++ {
    users = append(users, &mysql.TxyUser{...})
}
err := userRepo.CreateBatch(ctx, users)

// 使用索引优化查询
qb.Where("user_id = ? AND status = ?", userID, 1) // 确保有复合索引
```

## 🚨 注意事项

1. **事务管理**: 避免在事务中执行长时间操作
2. **连接池**: 合理配置数据库连接池参数
3. **错误处理**: 统一错误处理策略
4. **性能监控**: 监控数据库查询性能
5. **数据一致性**: 确保事务操作的数据一致性

## 🔄 迁移指南

### 从现有 Model 迁移到 Repository

1. **保留现有 Model**: 现有的 go-zero 生成的 Model 可以继续使用
2. **逐步引入 Repository**: 在新功能中使用 Repository，旧功能逐步迁移
3. **统一数据访问**: 最终目标是在所有服务中统一使用 Repository

```go
// 迁移前：使用 go-zero Model
user, err := l.svcCtx.UserModel.FindOne(l.ctx, req.UserId)

// 迁移后：使用 Repository
user, err := l.svcCtx.UserRepo.GetByID(l.ctx, req.UserId)
```

## 📞 技术支持

如果在使用过程中遇到问题，请：

1. 查看 `simple_example.go` 中的示例代码
2. 检查数据库连接配置
3. 确认模型定义是否正确
4. 查看错误日志获取详细信息

## 🎉 总结

这个 Repository 封装提供了：

- **统一的数据库访问接口**
- **完整的事务支持**
- **灵活的查询构建器**
- **模块化的设计**
- **类型安全的操作**

通过使用这个封装，你可以：

- 减少重复代码
- 提高开发效率
- 统一数据访问模式
- 简化事务管理
- 提高代码可维护性
