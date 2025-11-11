package main

import (
	"encoding/json"
	"fmt"
	"log"

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

func Test_curd(db *gorm.DB) {
	// // 插入测试数据
	// if err := seedData(db); err != nil {
	// 	log.Fatalf("插入测试数据失败: %v", err)
	// }

	// 查询所有用户数据
	users, _ := queryUsers(db)
	data, _ := json.MarshalIndent(users, "", "  ")
	fmt.Println("查询到的用户数据:", string(data))
	
}
