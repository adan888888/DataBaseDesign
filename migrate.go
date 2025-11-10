package main

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// migrate 数据库迁移函数
func migrate(db *gorm.DB) error {
	// 自动迁移所有模型
	err := db.AutoMigrate(
		&User{},
		&Address{},
		&Product{},
		&Order{},
		&OrderItem{},
	)
	if err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	fmt.Println("✓ 数据库表创建成功！")
	return nil
}

// generateOrderNo 生成订单号
// 格式: ORD + 年月日 + 6位随机数
func generateOrderNo() string {
	now := time.Now()
	return fmt.Sprintf("ORD%s%06d", now.Format("20060102"), now.Unix()%1000000)
}

// 订单状态常量
const (
	OrderStatusPending   int8 = 0 // 待支付
	OrderStatusPaid      int8 = 1 // 已支付
	OrderStatusShipped   int8 = 2 // 已发货
	OrderStatusCompleted int8 = 3 // 已完成
	OrderStatusCancelled int8 = 4 // 已取消
)

// 用户状态常量
const (
	UserStatusNormal int8 = 1 // 正常
	UserStatusBanned int8 = 0 // 禁用
)

// 商品状态常量
const (
	ProductStatusOnSale  int8 = 1 // 上架
	ProductStatusOffSale int8 = 0 // 下架
)

// generateProductNo 生成商品编号
// 格式: PROD + 年月日 + 6位随机数
func generateProductNo() string {
	now := time.Now()
	return fmt.Sprintf("PROD%s%06d", now.Format("20060102"), now.Unix()%1000000)
}

