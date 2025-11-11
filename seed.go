package main

import (
	"fmt"
	"time"
)

// seedData 插入测试数据
func seedData() error {
	fmt.Println("开始插入测试数据...")

	// 1. 插入用户数据
	users := []User{
		{
			Username: "zhangsan",
			Phone:    "13800138001",
			Email:    "zhangsan@example.com",
			Password: "hashed_password_123",
			Nickname: "张三",
			Avatar:   "https://example.com/avatar/zhangsan.jpg",
			Status:   UserStatusNormal,
		},
		{
			Username: "lisi",
			Phone:    "13800138002",
			Email:    "lisi@example.com",
			Password: "hashed_password_456",
			Nickname: "李四",
			Avatar:   "https://example.com/avatar/lisi.jpg",
			Status:   UserStatusNormal,
		},
		{
			Username: "wangwu",
			Phone:    "13800138003",
			Email:    "wangwu@example.com",
			Password: "hashed_password_789",
			Nickname: "王五",
			Avatar:   "https://example.com/avatar/wangwu.jpg",
			Status:   UserStatusNormal,
		},
	}

	if err := db.Create(&users).Error; err != nil {
		return fmt.Errorf("插入用户数据失败: %v", err)
	}
	
	// 验证 ID 是否被正确填充
	for i, user := range users {
		if user.ID == 0 {
			return fmt.Errorf("用户 %d 的 ID 未正确填充", i)
		}
	}
	fmt.Printf("✓ 成功插入 %d 个用户 (ID: %d, %d, %d)\n", len(users), users[0].ID, users[1].ID, users[2].ID)

	// 2. 插入地址数据
	addresses := []Address{
		{
			UserID:        users[0].ID,
			ReceiverName:  "张三",
			ReceiverPhone: "13800138001",
			Province:      "北京市",
			City:          "北京市",
			District:      "海淀区",
			Detail:        "中关村大街1号",
			PostalCode:    "100080",
			IsDefault:     true,
		},
		{
			UserID:        users[0].ID,
			ReceiverName:  "张三",
			ReceiverPhone: "13800138001",
			Province:      "北京市",
			City:          "北京市",
			District:      "朝阳区",
			Detail:        "三里屯路2号",
			PostalCode:    "100027",
			IsDefault:     false,
		},
		{
			UserID:        users[1].ID,
			ReceiverName:  "李四",
			ReceiverPhone: "13800138002",
			Province:      "上海市",
			City:          "上海市",
			District:      "浦东新区",
			Detail:        "陆家嘴环路1000号",
			PostalCode:    "200120",
			IsDefault:     true,
		},
	}

	if err := db.Create(&addresses).Error; err != nil {
		return fmt.Errorf("插入地址数据失败: %v", err)
	}
	
	// 验证地址 ID 是否被正确填充
	for i, address := range addresses {
		if address.ID == 0 {
			return fmt.Errorf("地址 %d 的 ID 未正确填充", i)
		}
	}
	fmt.Printf("✓ 成功插入 %d 个地址 (ID: %d, %d, %d)\n", len(addresses), addresses[0].ID, addresses[1].ID, addresses[2].ID)

	// 3. 插入商品数据
	products := []Product{
		{
			ProductNo:   generateProductNo(1),
			Name:        "iPhone 15 Pro",
			Description: "苹果最新款手机，A17 Pro芯片，6.1英寸屏幕",
			CategoryID:  1,
			Price:       7999.00,
			Stock:       100,
			Sales:       0,
			Image:       "https://example.com/images/iphone15pro.jpg",
			Status:      ProductStatusOnSale,
			Sort:        1,
		},
		{
			ProductNo:   generateProductNo(2),
			Name:        "AirPods Pro",
			Description: "苹果无线降噪耳机，主动降噪，空间音频",
			CategoryID:  2,
			Price:       1899.00,
			Stock:       200,
			Sales:       0,
			Image:       "https://example.com/images/airpodspro.jpg",
			Status:      ProductStatusOnSale,
			Sort:        2,
		},
		{
			ProductNo:   generateProductNo(3),
			Name:        "MacBook Pro 14英寸",
			Description: "苹果笔记本电脑，M3芯片，14英寸Liquid Retina XDR显示屏",
			CategoryID:  3,
			Price:       14999.00,
			Stock:       50,
			Sales:       0,
			Image:       "https://example.com/images/macbookpro14.jpg",
			Status:      ProductStatusOnSale,
			Sort:        3,
		},
		{
			ProductNo:   generateProductNo(4),
			Name:        "手机保护壳",
			Description: "iPhone 15 Pro专用保护壳，防摔防刮",
			CategoryID:  4,
			Price:       99.00,
			Stock:       500,
			Sales:       0,
			Image:       "https://example.com/images/phonecase.jpg",
			Status:      ProductStatusOnSale,
			Sort:        4,
		},
	}

	if err := db.Create(&products).Error; err != nil {
		return fmt.Errorf("插入商品数据失败: %v", err)
	}
	
	// 验证商品 ID 是否被正确填充
	for i, product := range products {
		if product.ID == 0 {
			return fmt.Errorf("商品 %d 的 ID 未正确填充", i)
		}
	}
	fmt.Printf("✓ 成功插入 %d 个商品 (ID: %d, %d, %d, %d)\n", len(products), products[0].ID, products[1].ID, products[2].ID, products[3].ID)

	// 4. 插入订单数据(用户、关联商品、地址)
	now := time.Now()
	payTime := now.Add(10 * time.Minute)
	orders := []Order{
		{
			OrderNo:       generateOrderNo(),
			UserID:        users[0].ID,
			AddressID:     addresses[0].ID,
			TotalAmount:   7999.00,
			DiscountAmount: 0.00,
			PayAmount:     7999.00,
			Status:        OrderStatusPaid,
			PayMethod:     "支付宝",
			PayTime:       &payTime,
			Remark:        "请尽快发货",
		},
		{
			OrderNo:       generateOrderNo(),
			UserID:        users[0].ID,
			AddressID:     addresses[0].ID,
			TotalAmount:   1998.00,
			DiscountAmount: 99.00,
			PayAmount:     1899.00,
			Status:        OrderStatusPending,
			PayMethod:     "",
			Remark:        "",
		},
		{
			OrderNo:       generateOrderNo(),
			UserID:        users[1].ID,
			AddressID:     addresses[2].ID,
			TotalAmount:   14999.00,
			DiscountAmount: 0.00,
			PayAmount:     14999.00,
			Status:        OrderStatusShipped,
			PayMethod:     "微信支付",
			PayTime:       &payTime,
			ShipTime:      &payTime,
			Remark:        "公司地址，工作日配送",
		},
	}

	if err := db.Create(&orders).Error; err != nil {
		return fmt.Errorf("插入订单数据失败: %v", err)
	}
	
	// 验证订单 ID 是否被正确填充
	for i, order := range orders {
		if order.ID == 0 {
			return fmt.Errorf("订单 %d 的 ID 未正确填充", i)
		}
	}
	fmt.Printf("✓ 成功插入 %d 个订单 (ID: %d, %d, %d)\n", len(orders), orders[0].ID, orders[1].ID, orders[2].ID)

	// 5. 插入订单明细数据
	orderItems := []OrderItem{
		// 订单1：iPhone 15 Pro × 1
		{
			OrderID:      orders[0].ID,
			ProductID:    products[0].ID,
			ProductName:  products[0].Name,
			ProductImage: products[0].Image,
			Price:        products[0].Price,
			Quantity:     1,
			Subtotal:     products[0].Price * 1,
		},
		// 订单2：AirPods Pro × 1 + 手机保护壳 × 1
		{
			OrderID:      orders[1].ID,
			ProductID:    products[1].ID,
			ProductName:  products[1].Name,
			ProductImage: products[1].Image,
			Price:        products[1].Price,
			Quantity:     1,
			Subtotal:     products[1].Price * 1,
		},
		{
			OrderID:      orders[1].ID,
			ProductID:    products[3].ID,
			ProductName:  products[3].Name,
			ProductImage: products[3].Image,
			Price:        products[3].Price,
			Quantity:     1,
			Subtotal:     products[3].Price * 1,
		},
		// 订单3：MacBook Pro × 1
		{
			OrderID:      orders[2].ID,
			ProductID:    products[2].ID,
			ProductName:  products[2].Name,
			ProductImage: products[2].Image,
			Price:        products[2].Price,
			Quantity:     1,
			Subtotal:     products[2].Price * 1,
		},
	}

	if err := db.Create(&orderItems).Error; err != nil {
		return fmt.Errorf("插入订单明细数据失败: %v", err)
	}
	fmt.Printf("✓ 成功插入 %d 个订单明细\n", len(orderItems))

	fmt.Println("✓ 所有测试数据插入完成！")
	return nil
}

