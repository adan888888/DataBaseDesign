package main

import (
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

// queryUsers 查询所有用户数据
func queryUsers(db *gorm.DB) ([]User, error) {
	var users []User
	if err := db.Debug().Find(&users).Error; err != nil {
		return nil, fmt.Errorf("查询用户数据失败: %v", err)
	}
	return users, nil
}

// queryUserOrders 查询指定用户的订单（包含订单明细）
func queryUserOrders(db *gorm.DB, userID uint) {
	var user User
	if err := db.Debug().
		// Preload("Orders") - 预加载用户的订单数据
		// 作用：在查询用户时，同时查询该用户的所有订单
		// 如果不使用 Preload，user.Orders 将为空（需要额外查询）
		// SQL 执行：SELECT * FROM orders WHERE user_id = ?
		Preload("Orders").
		Preload("Addresses").
		// Preload("Orders.OrderItems") - 嵌套预加载订单的商品明细
		// 作用：在加载订单时，同时加载每个订单的商品明细
		// "Orders.OrderItems" 表示：先加载 Orders，再加载每个 Order 的 OrderItems
		// SQL 执行：SELECT * FROM order_items WHERE order_id IN (?, ?, ...)
		// 这样可以在一次查询中获取：用户 -> 订单 -> 订单明细 的完整数据
		Preload("Orders.OrderItems").
		First(&user, userID).Error; err != nil {
		return
	}

	fmt.Printf("用户: %s (ID: %d)\n", user.Username, user.ID)
	fmt.Printf("订单数量: %d\n", len(user.Orders))

	for i, order := range user.Orders {
		fmt.Printf("\n订单 %d:\n", i+1)
		fmt.Printf("  订单号: %s\n", order.OrderNo)
		fmt.Printf("  订单状态: %d\n", order.Status)
		fmt.Printf("  订单金额: %.2f\n", order.TotalAmount)
		fmt.Printf("  实付金额: %.2f\n", order.PayAmount)
		fmt.Printf("  商品明细数量: %d\n", len(order.OrderItems))

		for j, item := range order.OrderItems {
			fmt.Printf("    商品 %d: %s × %d = %.2f\n",
				j+1, item.ProductName, item.Quantity, item.Subtotal)
		}
	}

	// JSON 格式输出
	orderData, _ := json.Marshal(user)
	fmt.Println("\n订单数据 (JSON):")
	fmt.Println(string(orderData))
}

func TestCurd(db *gorm.DB) {
	// // 插入测试数据
	// if err := seedData(db); err != nil {
	// 	log.Fatalf("插入测试数据失败: %v", err)
	// }

	// 查询所有用户数据
	// users, _ := queryUsers(db)
	// data, _ := json.MarshalIndent(users, "", "  ")
	// fmt.Println("查询到的用户数据:", string(data))

	// 查询第一个用户的订单
	fmt.Println("\n=== 查询第一个用户的订单 ===")
	queryUserOrders(db, 1)
}
