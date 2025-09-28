package repository

import (
	"context"
	"fmt"
	"log"
	"time"
)

// SimpleExample 简单使用示例
func SimpleExample() {
	// 1. 初始化仓储管理器
	mysqlDataSource := "root:password@tcp(127.0.0.1:3306)/lxtian_blog?charset=utf8mb4&parseTime=True&loc=Local"
	repoManager := NewRepositoryManagerWithConfig(mysqlDataSource)

	ctx := context.Background()

	// 2. 基础CRUD操作示例
	fmt.Println("=== 基础CRUD操作示例 ===")

	// 创建基础仓储
	// userRepo := NewBaseRepository[mysql.TxyUser](repoManager.GetDB())

	// 创建用户
	// user := &mysql.TxyUser{
	//     Nickname: "测试用户",
	//     // ... 其他字段
	// }
	// err := userRepo.Create(ctx, user)
	// if err != nil {
	//     log.Printf("创建用户失败: %v", err)
	// }

	// 查询用户
	// user, err := userRepo.GetByID(ctx, 1)
	// if err != nil {
	//     log.Printf("查询用户失败: %v", err)
	// } else {
	//     fmt.Printf("查询到用户: %+v\n", user)
	// }

	// 3. 查询构建器示例
	fmt.Println("=== 查询构建器示例 ===")

	// 使用查询构建器
	db := repoManager.GetDB()
	qb := NewQueryBuilder(db)

	// 构建查询
	var results []map[string]interface{}
	err := qb.Select("id", "nickname", "head_img").
		Where("status = ?", 1).
		Where("last_login_time > ?", time.Now().Unix()-86400*7). // 最近7天
		Order("last_login_time DESC").
		Limit(10).
		Execute(&results)

	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("查询结果数量: %d\n", len(results))
	}

	// 4. 事务操作示例
	fmt.Println("=== 事务操作示例 ===")

	err = repoManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// 在事务中执行多个操作

		// 操作1：模拟创建用户
		// user := &mysql.TxyUser{Nickname: "事务用户"}
		// if err := userRepo.Create(txCtx, user); err != nil {
		//     return fmt.Errorf("创建用户失败: %w", err)
		// }

		// 操作2：模拟创建文章
		// article := &mysql.TxyArticle{Title: "事务文章"}
		// if err := articleRepo.Create(txCtx, article); err != nil {
		//     return fmt.Errorf("创建文章失败: %w", err)
		// }

		fmt.Println("事务操作执行中...")
		time.Sleep(100 * time.Millisecond) // 模拟操作耗时

		return nil
	})

	if err != nil {
		log.Printf("事务操作失败: %v", err)
	} else {
		fmt.Println("事务操作成功")
	}

	// 5. 批量操作示例
	fmt.Println("=== 批量操作示例 ===")

	operations := []TransactionFunc{
		func(txCtx context.Context) error {
			// 批量更新操作1
			fmt.Println("执行批量操作1...")
			time.Sleep(50 * time.Millisecond)
			return nil
		},
		func(txCtx context.Context) error {
			// 批量更新操作2
			fmt.Println("执行批量操作2...")
			time.Sleep(50 * time.Millisecond)
			return nil
		},
	}

	err = repoManager.BatchExecute(ctx, operations)
	if err != nil {
		log.Printf("批量操作失败: %v", err)
	} else {
		fmt.Println("批量操作成功")
	}
}

// DatabaseConnectionExample 数据库连接示例
func DatabaseConnectionExample() {
	// 配置数据库连接
	config := &RepositoryConfig{
		MysqlDataSource: "root:password@tcp(127.0.0.1:3306)/lxtian_blog",
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		MaxLifetime:     3600,
	}

	repoManager := NewRepositoryManagerWithOptions(config)
	db := repoManager.GetDB()

	// 测试数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("获取数据库连接失败: %v", err)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		log.Printf("数据库连接测试失败: %v", err)
		return
	}

	fmt.Println("数据库连接测试成功")

	// 获取连接池状态
	fmt.Printf("最大打开连接数: %d\n", sqlDB.Stats().MaxOpenConnections)
	fmt.Printf("当前打开连接数: %d\n", sqlDB.Stats().OpenConnections)
	fmt.Printf("当前空闲连接数: %d\n", sqlDB.Stats().Idle)
}

// QueryBuilderExample 查询构建器详细示例
func QueryBuilderExample() {
	mysqlDataSource := "root:password@tcp(127.0.0.1:3306)/lxtian_blog?charset=utf8mb4&parseTime=True&loc=Local"
	repoManager := NewRepositoryManagerWithConfig(mysqlDataSource)

	db := repoManager.GetDB()
	qb := NewQueryBuilder(db)

	// 1. 基础查询
	fmt.Println("=== 基础查询示例 ===")

	var users []map[string]interface{}
	err := qb.Select("id", "nickname", "status").
		Where("status = ?", 1).
		Execute(&users)

	if err != nil {
		log.Printf("基础查询失败: %v", err)
	} else {
		fmt.Printf("查询到 %d 个用户\n", len(users))
	}

	// 2. 复杂条件查询
	fmt.Println("=== 复杂条件查询示例 ===")

	var articles []map[string]interface{}
	err = qb.Select("id", "title", "view_count", "created_at").
		Where("status = ?", 1).
		Where("view_count > ?", 100).
		Where("created_at > ?", time.Now().Unix()-86400*30). // 最近30天
		Order("view_count DESC").
		Limit(20).
		Execute(&articles)

	if err != nil {
		log.Printf("复杂查询失败: %v", err)
	} else {
		fmt.Printf("查询到 %d 篇文章\n", len(articles))
	}

	// 3. 分页查询
	fmt.Println("=== 分页查询示例 ===")

	page := 1
	pageSize := 10

	var comments []map[string]interface{}
	err = qb.Select("id", "content", "user_id", "created_at").
		Where("status = ?", 1).
		Order("created_at DESC").
		Page(page, pageSize).
		Execute(&comments)

	if err != nil {
		log.Printf("分页查询失败: %v", err)
	} else {
		fmt.Printf("第%d页查询到 %d 条评论\n", page, len(comments))
	}

	// 4. 统计查询
	fmt.Println("=== 统计查询示例 ===")

	var count int64
	err = qb.Count(&count)

	if err != nil {
		log.Printf("统计查询失败: %v", err)
	} else {
		fmt.Printf("总记录数: %d\n", count)
	}

	// 5. 原生SQL查询
	fmt.Println("=== 原生SQL查询示例 ===")

	var stats []map[string]interface{}
	err = qb.ExecuteRaw(`
		SELECT 
			DATE(FROM_UNIXTIME(created_at)) as date,
			COUNT(*) as count
		FROM txy_article 
		WHERE created_at > ? 
		GROUP BY DATE(FROM_UNIXTIME(created_at))
		ORDER BY date DESC
		LIMIT 7
	`, &stats, time.Now().Unix()-86400*7)

	if err != nil {
		log.Printf("原生SQL查询失败: %v", err)
	} else {
		fmt.Printf("最近7天的文章统计: %+v\n", stats)
	}
}

// TransactionExample 事务操作详细示例
func TransactionExample() {
	mysqlDataSource := "root:password@tcp(127.0.0.1:3306)/lxtian_blog?charset=utf8mb4&parseTime=True&loc=Local"
	repoManager := NewRepositoryManagerWithConfig(mysqlDataSource)

	ctx := context.Background()

	// 1. 简单事务
	fmt.Println("=== 简单事务示例 ===")

	err := repoManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// 模拟在事务中执行多个操作
		fmt.Println("开始事务操作...")

		// 操作1
		fmt.Println("执行操作1: 创建用户")
		time.Sleep(100 * time.Millisecond)

		// 操作2
		fmt.Println("执行操作2: 创建文章")
		time.Sleep(100 * time.Millisecond)

		// 操作3
		fmt.Println("执行操作3: 发送通知")
		time.Sleep(100 * time.Millisecond)

		fmt.Println("事务操作完成")
		return nil
	})

	if err != nil {
		log.Printf("简单事务失败: %v", err)
	} else {
		fmt.Println("简单事务成功")
	}

	// 2. 事务回滚示例
	fmt.Println("=== 事务回滚示例 ===")

	err = repoManager.WithTransaction(ctx, func(txCtx context.Context) error {
		fmt.Println("开始可能失败的事务...")

		// 操作1
		fmt.Println("执行操作1: 成功")
		time.Sleep(50 * time.Millisecond)

		// 操作2 - 模拟失败
		fmt.Println("执行操作2: 失败")
		return fmt.Errorf("模拟操作失败")
	})

	if err != nil {
		log.Printf("事务回滚成功: %v", err)
	} else {
		fmt.Println("事务回滚失败")
	}

	// 3. 批量事务操作
	fmt.Println("=== 批量事务操作示例 ===")

	operations := []TransactionFunc{
		func(txCtx context.Context) error {
			fmt.Println("批量操作1: 更新用户状态")
			time.Sleep(50 * time.Millisecond)
			return nil
		},
		func(txCtx context.Context) error {
			fmt.Println("批量操作2: 更新文章状态")
			time.Sleep(50 * time.Millisecond)
			return nil
		},
		func(txCtx context.Context) error {
			fmt.Println("批量操作3: 清理过期数据")
			time.Sleep(50 * time.Millisecond)
			return nil
		},
	}

	err = repoManager.BatchExecute(ctx, operations)
	if err != nil {
		log.Printf("批量事务失败: %v", err)
	} else {
		fmt.Println("批量事务成功")
	}
}
