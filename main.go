package main

import (
	"fmt"
	"log"
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
	db, err := connectDB(config)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	fmt.Println("✓ 数据库连接成功！")

	// 执行数据库迁移
	if err := migrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	Test_curd(db)

	// 可以在这里使用 db 进行数据库操作
	_ = db
}
