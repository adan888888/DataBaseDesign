package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 从环境变量读取配置，如果没有则使用默认值
	config := DBConfig{
		Host:     getEnv("DB_HOST", "127.0.0.1"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "mima123"),
		Database: getEnv("DB_NAME", "table_design"),
	}

	// 连接数据库
	var err error
	db, err = connectDB(config)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	fmt.Println("✓ 数据库连接成功！")

	// 执行数据库迁移
	if err := migrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 设置 Gin 模式（开发模式会显示更多调试信息）
	ginMode := getEnv("GIN_MODE", gin.DebugMode)
	gin.SetMode(ginMode)

	// 设置路由
	r := SetupRoutes()

	// 获取端口号，默认 8080
	port := getEnv("PORT", "8080")

	fmt.Printf("✓ 服务器启动成功！\n")
	fmt.Printf("✓ 访问地址: http://localhost:%s\n", port)
	fmt.Printf("✓ API 文档:\n")
	fmt.Printf("  - 健康检查: GET http://localhost:%s/health\n", port)
	fmt.Printf("  - 查询所有用户: GET http://localhost:%s/users\n", port)
	fmt.Printf("  - 查询用户订单: GET http://localhost:%s/users/:id/orders\n", port)
	fmt.Printf("  - 查询用户订单及商品: GET http://localhost:%s/users/:id/orders/products\n", port)
	fmt.Printf("  - 查询所有商品: GET http://localhost:%s/products\n", port)
	fmt.Printf("  - 查询商品订单: GET http://localhost:%s/products/:id/orders\n", port)
	fmt.Printf("  - 查询商品统计: GET http://localhost:%s/products/:id/stats\n", port)
	fmt.Printf("  - 查询所有订单: GET http://localhost:%s/orders\n", port)
	fmt.Printf("  - 查询订单商品: GET http://localhost:%s/orders/:id/products\n", port)
	fmt.Printf("  - 插入测试数据: POST http://localhost:%s/seed\n", port)

	// 启动服务器
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
		os.Exit(1)
	}
}
