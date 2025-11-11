package main

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes() *gin.Engine {
	// 创建 Gin 引擎
	r := gin.Default()

	// 添加 CORS 中间件（如果需要跨域访问）
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "服务运行正常",
		})
	})	
	
	// 测试数据接口
	r.POST("/seed", SeedData) // 插入测试数据

	// 用户相关路由
	r.GET("/users", GetUsers)                           // 查询所有用户
	r.GET("/users/:id/orders", GetUserOrders)            // 查询用户的订单
	r.GET("/users/:id/orders/products", GetUserOrdersWithProducts) // 查询用户的订单及商品

	// 商品相关路由
	r.GET("/products", GetProducts)                     // 查询所有商品
	r.GET("/products/:id", GetProduct)                   // 查询单个商品
	r.GET("/products/:id/orders", GetProductOrders)      // 查询商品被哪些订单购买
	r.GET("/products/:id/stats", GetProductSalesStats)   // 查询商品销售统计

	// 订单相关路由
	r.GET("/orders", GetOrders)                          // 查询所有订单
	r.GET("/orders/:id", GetOrder)                       // 查询单个订单
	r.GET("/orders/:id/products", GetOrderProducts)       // 查询订单包含哪些商品





	return r
}

