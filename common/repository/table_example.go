package repository

import (
	"context"
	"fmt"
	"log"
	"time"
)

// TableRepositoryExample 按表分开的仓储使用示例
func TableRepositoryExample() {
	// 1. 初始化仓储管理器
	mysqlDataSource := "root:password@tcp(127.0.0.1:3306)/lxtian_blog?charset=utf8mb4&parseTime=True&loc=Local"
	repoManager := NewRepositoryManagerWithConfig(mysqlDataSource)

	ctx := context.Background()

	fmt.Println("=== 按表分开的仓储使用示例 ===")

	// 2. 用户模块仓储示例
	fmt.Println("--- 用户模块仓储示例 ---")

	// 创建用户相关仓储
	// userRepo := user.NewTxyUserRepository(repoManager.GetDB())
	// rolesRepo := user.NewTxyRolesRepository(repoManager.GetDB())
	// permsRepo := user.NewTxyPermissionsRepository(repoManager.GetDB())

	// 用户表操作示例
	// user, err := userRepo.GetByUid(ctx, 12345)
	// if err != nil {
	//     log.Printf("查询用户失败: %v", err)
	// } else {
	//     fmt.Printf("查询到用户: %+v\n", user)
	// }

	// 更新最后登录信息
	// err = userRepo.UpdateLastLogin(ctx, 12345, uint64(time.Now().Unix()), "192.168.1.1")
	// if err != nil {
	//     log.Printf("更新登录信息失败: %v", err)
	// } else {
	//     fmt.Println("更新登录信息成功")
	// }

	// 获取活跃用户数量
	// activeCount, err := userRepo.GetActiveUserCount(ctx, 7) // 最近7天
	// if err != nil {
	//     log.Printf("获取活跃用户数量失败: %v", err)
	// } else {
	//     fmt.Printf("最近7天活跃用户数量: %d\n", activeCount)
	// }

	// 3. Web模块仓储示例
	fmt.Println("--- Web模块仓储示例 ---")

	// 创建Web相关仓储
	// articleRepo := web.NewTxyArticleRepository(repoManager.GetDB())
	// categoryRepo := web.NewTxyCategoryRepository(repoManager.GetDB())
	// tagRepo := web.NewTxyTagRepository(repoManager.GetDB())
	// commentRepo := web.NewTxyCommentRepository(repoManager.GetDB())

	// 文章表操作示例
	// articles, err := articleRepo.GetPublishedArticles(ctx, 1, 10) // 第1页，每页10条
	// if err != nil {
	//     log.Printf("获取已发布文章失败: %v", err)
	// } else {
	//     fmt.Printf("获取到 %d 篇已发布文章\n", len(articles))
	// }

	// 获取热门文章
	// popularArticles, err := articleRepo.GetPopularArticles(ctx, 5)
	// if err != nil {
	//     log.Printf("获取热门文章失败: %v", err)
	// } else {
	//     fmt.Printf("获取到 %d 篇热门文章\n", len(popularArticles))
	// }

	// 更新文章浏览量
	// err = articleRepo.UpdateViewCount(ctx, 1)
	// if err != nil {
	//     log.Printf("更新文章浏览量失败: %v", err)
	// } else {
	//     fmt.Println("更新文章浏览量成功")
	// }

	// 分类表操作示例
	// categories, err := categoryRepo.GetTree(ctx)
	// if err != nil {
	//     log.Printf("获取分类树失败: %v", err)
	// } else {
	//     fmt.Printf("获取到 %d 个分类\n", len(categories))
	// }

	// 标签表操作示例
	// popularTags, err := tagRepo.GetPopularTags(ctx, 10)
	// if err != nil {
	//     log.Printf("获取热门标签失败: %v", err)
	// } else {
	//     fmt.Printf("获取到 %d 个热门标签\n", len(popularTags))
	// }

	// 评论表操作示例
	// comments, total, err := commentRepo.GetByArticleId(ctx, 1, 1, 10)
	// if err != nil {
	//     log.Printf("获取文章评论失败: %v", err)
	// } else {
	//     fmt.Printf("文章评论总数: %d，当前页评论数: %d\n", total, len(comments))
	// }

	// 4. 支付模块仓储示例
	fmt.Println("--- 支付模块仓储示例 ---")

	// 创建支付相关仓储
	// paymentRepo := payment.NewPaymentOrderRepository(repoManager.GetDB())
	// orderRepo := payment.NewTxyOrderRepository(repoManager.GetDB())

	// 支付订单表操作示例
	// payment, err := paymentRepo.GetByPaymentId(ctx, "PAY_123456")
	// if err != nil {
	//     log.Printf("查询支付订单失败: %v", err)
	// } else {
	//     fmt.Printf("查询到支付订单: %+v\n", payment)
	// }

	// 更新支付状态
	// err = paymentRepo.UpdateStatus(ctx, "PAY_123456", "PAID")
	// if err != nil {
	//     log.Printf("更新支付状态失败: %v", err)
	// } else {
	//     fmt.Println("更新支付状态成功")
	// }

	// 获取过期订单
	// expiredOrders, err := paymentRepo.GetExpiredOrders(ctx)
	// if err != nil {
	//     log.Printf("获取过期订单失败: %v", err)
	// } else {
	//     fmt.Printf("获取到 %d 个过期订单\n", len(expiredOrders))
	// }

	// 业务订单表操作示例
	// order, err := orderRepo.GetByOrderSn(ctx, "ORDER_123456")
	// if err != nil {
	//     log.Printf("查询业务订单失败: %v", err)
	// } else {
	//     fmt.Printf("查询到业务订单: %+v\n", order)
	// }

	// 统计用户总金额
	// totalAmount, err := orderRepo.GetTotalAmountByUserId(ctx, 12345)
	// if err != nil {
	//     log.Printf("统计用户总金额失败: %v", err)
	// } else {
	//     fmt.Printf("用户总消费金额: %.2f\n", totalAmount)
	// }

	fmt.Println("=== 示例完成 ===")
}

// ComplexQueryExample 复杂查询示例
func ComplexQueryExample() {
	mysqlDataSource := "root:password@tcp(127.0.0.1:3306)/lxtian_blog?charset=utf8mb4&parseTime=True&loc=Local"
	repoManager := NewRepositoryManagerWithConfig(mysqlDataSource)

	db := repoManager.GetDB()
	qb := NewQueryBuilder(db)
	ctx := context.Background()

	fmt.Println("=== 复杂查询示例 ===")

	// 1. 多表关联查询示例
	fmt.Println("--- 多表关联查询示例 ---")

	// 查询用户及其角色信息
	var userRoles []map[string]interface{}
	err := qb.Select("u.id", "u.nickname", "r.name as role_name", "r.description as role_desc").
		From("txy_user u").
		Join("LEFT JOIN txy_user_roles ur ON u.id = ur.user_id").
		Join("LEFT JOIN txy_roles r ON ur.role_id = r.id").
		Where("u.status = ?", 1).
		Where("r.status = ?", 1).
		Limit(10).
		Execute(&userRoles)

	if err != nil {
		log.Printf("多表关联查询失败: %v", err)
	} else {
		fmt.Printf("查询到 %d 个用户角色信息\n", len(userRoles))
	}

	// 2. 聚合查询示例
	fmt.Println("--- 聚合查询示例 ---")

	// 统计各分类下的文章数量
	var categoryStats []map[string]interface{}
	err = qb.ExecuteRaw(`
		SELECT 
			c.id,
			c.name,
			c.slug,
			COUNT(a.id) as article_count,
			SUM(a.view_count) as total_views
		FROM txy_category c
		LEFT JOIN txy_article a ON c.id = a.category_id AND a.status = 1
		WHERE c.status = 1
		GROUP BY c.id, c.name, c.slug
		ORDER BY article_count DESC
		LIMIT 10
	`, &categoryStats)

	if err != nil {
		log.Printf("聚合查询失败: %v", err)
	} else {
		fmt.Printf("查询到 %d 个分类统计信息\n", len(categoryStats))
	}

	// 3. 时间范围查询示例
	fmt.Println("--- 时间范围查询示例 ---")

	// 查询最近30天的文章
	var recentArticles []map[string]interface{}
	thirtyDaysAgo := time.Now().Unix() - 30*24*3600

	err = qb.Select("id", "title", "view_count", "created_at").
		From("txy_article").
		Where("status = ?", 1).
		Where("created_at > ?", thirtyDaysAgo).
		Order("view_count DESC").
		Limit(20).
		Execute(&recentArticles)

	if err != nil {
		log.Printf("时间范围查询失败: %v", err)
	} else {
		fmt.Printf("查询到 %d 篇最近30天的文章\n", len(recentArticles))
	}

	// 4. 子查询示例
	fmt.Println("--- 子查询示例 ---")

	// 查询评论数最多的文章
	var topArticles []map[string]interface{}
	err = qb.ExecuteRaw(`
		SELECT 
			a.id,
			a.title,
			a.view_count,
			a.comment_count,
			c.name as category_name
		FROM txy_article a
		LEFT JOIN txy_category c ON a.category_id = c.id
		WHERE a.status = 1
		AND a.comment_count = (
			SELECT MAX(comment_count) 
			FROM txy_article 
			WHERE status = 1
		)
		ORDER BY a.view_count DESC
		LIMIT 10
	`, &topArticles)

	if err != nil {
		log.Printf("子查询失败: %v", err)
	} else {
		fmt.Printf("查询到 %d 篇评论最多的文章\n", len(topArticles))
	}

	fmt.Println("=== 复杂查询示例完成 ===")
}

// TransactionWithMultipleTables 多表事务操作示例
func TransactionWithMultipleTables() {
	mysqlDataSource := "root:password@tcp(127.0.0.1:3306)/lxtian_blog?charset=utf8mb4&parseTime=True&loc=Local"
	repoManager := NewRepositoryManagerWithConfig(mysqlDataSource)

	ctx := context.Background()

	fmt.Println("=== 多表事务操作示例 ===")

	// 在事务中同时操作多个表
	err := repoManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// 创建基础仓储实例
		// userRepo := user.NewTxyUserRepository(repoManager.GetDB())
		// articleRepo := web.NewTxyArticleRepository(repoManager.GetDB())
		// commentRepo := web.NewTxyCommentRepository(repoManager.GetDB())

		fmt.Println("开始多表事务操作...")

		// 操作1：创建用户
		// user := &mysql.TxyUser{
		//     Nickname: "事务用户",
		//     // ... 其他字段
		// }
		// if err := userRepo.Create(txCtx, user); err != nil {
		//     return fmt.Errorf("创建用户失败: %w", err)
		// }
		fmt.Println("操作1：创建用户（模拟）")

		// 操作2：创建文章
		// article := &mysql.TxyArticle{
		//     Title: "事务文章",
		//     // ... 其他字段
		// }
		// if err := articleRepo.Create(txCtx, article); err != nil {
		//     return fmt.Errorf("创建文章失败: %w", err)
		// }
		fmt.Println("操作2：创建文章（模拟）")

		// 操作3：创建评论
		// comment := &mysql.TxyComment{
		//     Content: "事务评论",
		//     // ... 其他字段
		// }
		// if err := commentRepo.Create(txCtx, comment); err != nil {
		//     return fmt.Errorf("创建评论失败: %w", err)
		// }
		fmt.Println("操作3：创建评论（模拟）")

		// 模拟一些处理时间
		time.Sleep(100 * time.Millisecond)

		fmt.Println("多表事务操作完成")
		return nil
	})

	if err != nil {
		log.Printf("多表事务操作失败: %v", err)
	} else {
		fmt.Println("多表事务操作成功")
	}

	fmt.Println("=== 多表事务操作示例完成 ===")
}
