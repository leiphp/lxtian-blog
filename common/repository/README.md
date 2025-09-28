# Repository 通用数据库操作封装

这是一个基于 GORM 的通用数据库操作封装，支持事务操作，可以前后端复用。

## 🚀 特性

- **通用CRUD操作**: 提供基础的增删改查方法
- **事务支持**: 完整的事务管理，支持嵌套事务
- **查询构建器**: 灵活的查询构建器，支持复杂查询
- **模块化设计**: 按服务模块分离，易于维护和扩展
- **类型安全**: 使用 Go 泛型，提供类型安全
- **批量操作**: 支持批量插入、更新、删除
- **统计查询**: 内置常用统计方法
- **分页支持**: 内置分页查询功能

## 📁 目录结构

```
common/repository/
├── base.go              # 基础仓储实现
├── transaction.go       # 事务管理器
├── manager.go          # 仓储管理器
├── interfaces.go       # 基础接口定义
├── simple_example.go   # 使用示例
├── README.md           # 文档说明
├── USAGE.md            # 使用指南
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

## 🔧 快速开始

### 1. 初始化仓储管理器

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
// 获取用户仓储
userRepo := repoManager.GetUserRepository()
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

### 5. 批量操作

```go
operations := []func(ctx context.Context) error{
    func(txCtx context.Context) error {
        return userRepo.BatchUpdateUserStatus(txCtx, []uint64{1,2,3}, 1)
    },
    func(txCtx context.Context) error {
        return articleRepo.BatchUpdateStatus(txCtx, []uint64{1,2,3}, 1)
    },
}

err := repoManager.BatchExecute(ctx, operations)
```

## 📋 支持的仓储

### 用户模块 (User)

- `UserRepository`: 用户基础操作
- `GetByUsername()`: 根据用户名查询
- `GetByEmail()`: 根据邮箱查询
- `GetUsersByRole()`: 根据角色查询用户
- `BatchAssignRole()`: 批量分配角色

### 支付模块 (Payment)

- `PaymentOrderRepository`: 支付订单操作
- `BusinessOrderRepository`: 业务订单操作
- `PaymentRefundRepository`: 支付退款操作
- `GetByPaymentId()`: 根据支付ID查询
- `UpdateTradeInfo()`: 更新交易信息

### Web模块 (Web)

- `ArticleRepository`: 文章操作
- `CategoryRepository`: 分类操作
- `TagRepository`: 标签操作
- `CommentRepository`: 评论操作
- `SearchArticles()`: 文章搜索
- `GetPopularArticles()`: 获取热门文章

## 🔧 配置选项

```go
config := &repository.RepositoryConfig{
    MysqlDataSource: "root:password@tcp(127.0.0.1:3306)/lxtian_blog",
    MaxIdleConns:    10,
    MaxOpenConns:    100,
    MaxLifetime:     3600,
}

repoManager := repository.NewRepositoryManagerWithOptions(config)
```

## 🎯 最佳实践

### 1. 错误处理

```go
user, err := userRepo.GetByID(ctx, 1)
if err != nil {
    if strings.Contains(err.Error(), "not found") {
        // 处理未找到的情况
        return nil, errors.New("用户不存在")
    }
    // 处理其他错误
    return nil, fmt.Errorf("查询用户失败: %w", err)
}
```

### 2. 上下文传递

```go
// 在服务层传递上下文
func (s *UserService) GetUser(ctx context.Context, userID uint64) (*User, error) {
    user, err := s.repoManager.GetUserRepository().GetByID(ctx, userID)
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

## 🛠️ 扩展指南

### 添加新的仓储

1. 在对应模块目录下创建新的仓储文件
2. 实现对应的接口
3. 在 `manager.go` 中注册新的仓储

```go
// 在 manager.go 中添加
type RepositoryManager struct {
    // ... 现有字段
    NewRepository NewRepositoryInterface
}

func (rm *RepositoryManager) initializeRepositories() {
    // ... 现有初始化
    rm.NewRepository = new.NewRepository(rm.db)
}
```

### 自定义查询方法

```go
// 在仓储中添加自定义方法
func (r *userRepository) GetActiveUsersWithRecentActivity(ctx context.Context, days int) ([]*mysql.TxyUser, error) {
    db := r.GetDB(ctx)
    var users []*mysql.TxyUser
    
    err := db.Where("status = ? AND last_login_time > ?", 1, time.Now().Unix()-int64(days*24*3600)).
        Find(&users).Error
    
    return users, err
}
```

## 📝 注意事项

1. **事务管理**: 避免在事务中执行长时间操作
2. **连接池**: 合理配置数据库连接池参数
3. **错误处理**: 统一错误处理策略
4. **性能监控**: 监控数据库查询性能
5. **数据一致性**: 确保事务操作的数据一致性

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目。

## 📄 许可证

本项目采用 MIT 许可证。
