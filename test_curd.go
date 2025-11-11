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

// queryProductOrders 查询商品被哪些订单购买（多对多关系演示）
// 通过 OrderItem 中间表实现：Product -> OrderItem -> Order
func queryProductOrders(db *gorm.DB, productID uint) {
	var product Product
	if err := db.Debug().
		// Preload("OrderItems") - 预加载商品的订单明细
		// 作用：查询该商品出现在哪些订单明细中
		Preload("OrderItems").
		// Preload("OrderItems.Order") - 嵌套预加载订单明细所属的订单
		// 作用：通过订单明细获取购买该商品的所有订单
		// 这就是多对多关系的查询：商品 -> 订单明细 -> 订单
		Preload("OrderItems.Order").
		Preload("OrderItems.Order.User"). // 可选：加载订单的用户信息
		First(&product, productID).Error; err != nil {
		fmt.Printf("查询商品失败: %v\n", err)
		return
	}

	fmt.Printf("\n=== 商品：%s (ID: %d) ===\n", product.Name, product.ID)
	fmt.Printf("商品编号: %s\n", product.ProductNo)
	fmt.Printf("价格: %.2f\n", product.Price)
	fmt.Printf("库存: %d\n", product.Stock)
	fmt.Printf("销量: %d\n", product.Sales)

	fmt.Printf("\n该商品出现在 %d 个订单明细中\n", len(product.OrderItems))

	// 统计订单数量（去重）
	orderMap := make(map[uint]Order)
	for _, item := range product.OrderItems {
		if item.Order.ID != 0 {
			orderMap[item.Order.ID] = item.Order
		}
	}

	fmt.Printf("该商品被 %d 个订单购买\n\n", len(orderMap))

	// JSON 格式输出
	productx, _ := json.Marshal(product)
	fmt.Println("\n订单数据 (JSON):")
	fmt.Println(string(productx))
}

// queryOrderProducts 查询订单包含哪些商品（多对多关系演示）
// 通过 OrderItem 中间表实现：Order -> OrderItem -> Product
func queryOrderProducts(db *gorm.DB, orderID uint) {
	var order Order
	if err := db.Debug().
		// Preload("OrderItems") - 预加载订单的商品明细
		// 作用：查询该订单包含哪些商品明细
		Preload("OrderItems").
		// Preload("OrderItems.Product") - 嵌套预加载商品明细对应的商品
		// 作用：通过订单明细获取订单包含的所有商品
		// 这就是多对多关系的查询：订单 -> 订单明细 -> 商品
		Preload("OrderItems.Product").
		Preload("User"). // 可选：加载订单的用户信息
		First(&order, orderID).Error; err != nil {
		fmt.Printf("查询订单失败: %v\n", err)
		return
	}

	fmt.Printf("\n=== 订单：%s (ID: %d) ===\n", order.OrderNo, order.ID)
	fmt.Printf("订单状态: %d\n", order.Status)
	fmt.Printf("订单金额: %.2f\n", order.TotalAmount)
	fmt.Printf("实付金额: %.2f\n", order.PayAmount)
	if order.User.ID != 0 {
		fmt.Printf("用户: %s (ID: %d)\n", order.User.Username, order.User.ID)
	}

	fmt.Printf("\n该订单包含 %d 个商品明细\n\n", len(order.OrderItems))

	// 显示每个商品的信息
	for i, item := range order.OrderItems {
		fmt.Printf("商品 %d:\n", i+1)
		fmt.Printf("  商品名称: %s\n", item.ProductName)
		fmt.Printf("  购买数量: %d\n", item.Quantity)
		fmt.Printf("  单价: %.2f\n", item.Price)
		fmt.Printf("  小计: %.2f\n", item.Subtotal)

		// 如果商品信息已加载，显示更多信息
		if item.Product.ID != 0 {
			fmt.Printf("  商品编号: %s\n", item.Product.ProductNo)
			fmt.Printf("  商品库存: %d\n", item.Product.Stock)
		}
		fmt.Println()
	}
	// JSON 格式输出
	orderx, _ := json.Marshal(order)
	fmt.Println("\n订单数据 (JSON):")
	fmt.Println(string(orderx))
}

// queryUserOrdersWithProducts 查询用户的所有订单及每个订单的商品（演示多对多关系）
// 关系说明：
//   - User -> Order：一对多（一个用户可以有多个订单）
//   - Order -> Product：多对多（一个订单可以有多个商品，一个商品可以出现在多个订单中）
//
// 示例：订单1（iPhone 15, iPhone 16），订单2（iPhone 15, iPhone 16, MacBook Air 13）
func queryUserOrdersWithProducts(db *gorm.DB, userID uint) {
	var user User
	if err := db.Debug().
		// Preload("Orders") - 预加载用户的所有订单
		Preload("Orders").
		// Preload("Orders.OrderItems") - 预加载每个订单的商品明细
		Preload("Orders.OrderItems").
		// Preload("Orders.OrderItems.Product") - 预加载商品明细对应的商品信息
		Preload("Orders.OrderItems.Product").
		First(&user, userID).Error; err != nil {
		fmt.Printf("查询用户失败: %v\n", err)
		return
	}

	fmt.Printf("\n用户: %s (ID: %d)\n", user.Username, user.ID)
	fmt.Printf("该用户共有 %d 个订单\n\n", len(user.Orders))

	// 遍历每个订单
	for i, order := range user.Orders {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("订单 %d: %s\n", i+1, order.OrderNo)
		fmt.Printf("订单状态: %d | 订单金额: %.2f 元 | 实付金额: %.2f 元\n",
			order.Status, order.TotalAmount, order.PayAmount)
		fmt.Printf("包含商品数量: %d\n\n", len(order.OrderItems))

		// 遍历订单中的每个商品
		for j, item := range order.OrderItems {
			fmt.Printf("  商品 %d:\n", j+1)
			fmt.Printf("    - 商品名称: %s\n", item.ProductName)
			fmt.Printf("    - 购买数量: %d\n", item.Quantity)
			fmt.Printf("    - 单价: %.2f 元\n", item.Price)
			fmt.Printf("    - 小计: %.2f 元\n", item.Subtotal)

			// 如果商品信息已加载，显示商品编号
			if item.Product.ID != 0 {
				fmt.Printf("    - 商品编号: %s\n", item.Product.ProductNo)
			}
			fmt.Println()
		}
	}

	// 统计：哪些商品出现在多个订单中
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("商品统计（哪些商品出现在多个订单中）:")

	productOrderMap := make(map[string][]string) // 商品名称 -> 订单号列表

	for _, order := range user.Orders {
		for _, item := range order.OrderItems {
			productName := item.ProductName
			if _, exists := productOrderMap[productName]; !exists {
				productOrderMap[productName] = []string{}
			}
			// 检查订单号是否已存在（避免重复）
			exists := false
			for _, orderNo := range productOrderMap[productName] {
				if orderNo == order.OrderNo {
					exists = true
					break
				}
			}
			if !exists {
				productOrderMap[productName] = append(productOrderMap[productName], order.OrderNo)
			}
		}
	}

	for productName, orderNos := range productOrderMap {
		if len(orderNos) > 1 {
			fmt.Printf("  ✓ %s 出现在 %d 个订单中: %v\n", productName, len(orderNos), orderNos)
		} else {
			fmt.Printf("  - %s 只出现在 1 个订单中: %v\n", productName, orderNos)
		}
	}
}

// queryProductSalesStats 统计商品的销售情况（多对多关系统计）
func queryProductSalesStats(db *gorm.DB, productID uint) {
	var product Product
	if err := db.Debug().
		Preload("OrderItems").
		Preload("OrderItems.Order").
		First(&product, productID).Error; err != nil {
		fmt.Printf("查询商品失败: %v\n", err)
		return
	}

	fmt.Printf("\n=== 商品销售统计：%s ===\n", product.Name)

	// 统计总销量
	totalQuantity := 0
	totalAmount := 0.0
	orderCount := 0
	orderMap := make(map[uint]bool)

	for _, item := range product.OrderItems {
		totalQuantity += item.Quantity
		totalAmount += item.Subtotal
		if !orderMap[item.OrderID] {
			orderMap[item.OrderID] = true
			orderCount++
		}
	}

	fmt.Printf("总销量: %d 件\n", totalQuantity)
	fmt.Printf("总销售额: %.2f 元\n", totalAmount)
	fmt.Printf("订单数量: %d 个\n", orderCount)
	if orderCount > 0 {
		fmt.Printf("平均订单金额: %.2f 元\n", totalAmount/float64(orderCount))
	}

	// JSON 格式输出
	productx, _ := json.Marshal(product)
	fmt.Println(string(productx))
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
	// fmt.Println("\n=== 查询第一个用户的订单 ===")
	// queryUserOrders(db, 1)

	// 多对多关系查询演示
	fmt.Println("\n=== 多对多关系查询演示 ===")
	/*	用户（User）
		└── 订单1（Order）
				├── iPhone 15  (Product)
				└── iPhone 16  (Product)

		└── 订单2（Order）
				├── iPhone 15      (Product) ← 同一个商品出现在不同订单
				├── iPhone 16      (Product) ← 同一个商品出现在不同订单
				└── MacBook Air 13 (Product)*/
	// 演示：一个用户有多个订单，每个订单包含多个商品
	// 例如：订单1（iPhone 15, iPhone 16），订单2（iPhone 15, iPhone 16, MacBook Air 13）
	//fmt.Println("\n=== 演示：一个用户的所有订单及商品 ===")
	//queryUserOrdersWithProducts(db, 1) // 查询用户ID=1的所有订单和商品

	// 1. 查询商品被哪些订单购买
	//fmt.Println("\n1. 查询商品被哪些订单购买（Product -> OrderItem -> Order）")
	//queryProductOrders(db, 3) // 查询商品ID=1的订单

	// 2. 查询订单包含哪些商品
	//fmt.Println("\n2. 查询订单包含哪些商品（Order -> OrderItem -> Product）")
	//queryOrderProducts(db, 1) // 查询订单ID=1的商品

	// 3. 统计商品销售情况
	//fmt.Println("\n3. 统计商品销售情况")
	//queryProductSalesStats(db, 3) // 统计商品ID=1的销售情况

	var product []Product
	db.Debug().Find(&product)
	marshal, _ := json.MarshalIndent(product, "", " ")
	fmt.Println(string(marshal))
}
