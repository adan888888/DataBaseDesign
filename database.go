package main

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// db 全局数据库实例
var db *gorm.DB

// DBConfig 数据库配置结构体
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// connectDB 连接数据库并返回 GORM 实例
func connectDB(config DBConfig) (*gorm.DB, error) {
	// 构建 DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.Database)

	fmt.Printf("正在连接数据库: %s@%s:%s/%s\n", config.User, config.Host, config.Port, config.Database)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v\n\n常见问题排查:\n"+
			"1. 检查数据库服务是否已启动\n"+
			"2. 检查用户名和密码是否正确\n"+
			"3. 检查数据库是否存在: %s\n"+
			"4. 检查用户是否有访问权限\n"+
			"5. 如果是远程连接，检查防火墙和 MySQL 的 bind-address 配置\n"+
			"6. 检查 MySQL 用户是否允许从当前 IP 连接", err, config.Database)
	}

	// 测试数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %v", err)
	}

	// Ping 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接测试失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)   // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)  // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(0) // 设置连接可复用的最大时间（0 表示不限制）

	return db, nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

