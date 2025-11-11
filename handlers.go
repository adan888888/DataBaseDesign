package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// GetUsers 查询所有用户
// GET /users
func GetUsers(c *gin.Context) {
	users, err := queryUsers(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "查询成功",
		Data:    users,
	})
}

// GetUserOrders 查询指定用户的订单
// GET /users/:id/orders
func GetUserOrders(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的用户ID",
		})
		return
	}

	var user User
	if err := db.Debug().
		Preload("Orders").
		Preload("Addresses").
		Preload("Orders.OrderItems").
		First(&user, uint(userID)).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "查询成功",
		Data:    user,
	})
}

// GetUserOrdersWithProducts 查询用户的所有订单及每个订单的商品
// GET /users/:id/orders/products
func GetUserOrdersWithProducts(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的用户ID",
		})
		return
	}

	var user User
	if err := db.Debug().
		Preload("Orders").
		Preload("Orders.OrderItems").
		Preload("Orders.OrderItems.Product").
		First(&user, uint(userID)).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "查询成功",
		Data:    user,
	})
}

// GetProductOrders 查询商品被哪些订单购买
// GET /products/:id/orders
func GetProductOrders(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的商品ID",
		})
		return
	}

	var product Product
	if err := db.Debug().
		Preload("OrderItems").
		Preload("OrderItems.Order").
		Preload("OrderItems.Order.User").
		First(&product, uint(productID)).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "商品不存在",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "查询成功",
		Data:    product,
	})
}

// GetOrderProducts 查询订单包含哪些商品
// GET /orders/:id/products
func GetOrderProducts(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的订单ID",
		})
		return
	}

	var order Order
	if err := db.Debug().
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Preload("User").
		First(&order, uint(orderID)).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "订单不存在",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "查询成功",
		Data:    order,
	})
}

// GetProductSalesStats 统计商品的销售情况
// GET /products/:id/stats
func GetProductSalesStats(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的商品ID",
		})
		return
	}

	var product Product
	if err := db.Debug().
		Preload("OrderItems").
		Preload("OrderItems.Order").
		First(&product, uint(productID)).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "商品不存在",
		})
		return
	}

	// 统计销售数据
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

	stats := map[string]interface{}{
		"product":        product,
		"total_quantity": totalQuantity,
		"total_amount":    totalAmount,
		"order_count":    orderCount,
		"average_amount": 0.0,
	}

	if orderCount > 0 {
		stats["average_amount"] = totalAmount / float64(orderCount)
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "查询成功",
		Data:    stats,
	})
}

// GetProducts 查询所有商品
// GET /products
func GetProducts(c *gin.Context) {
	var products []Product
	if err := db.Debug().Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "查询成功",
		Data:    products,
	})
}

// GetProduct 查询单个商品
// GET /products/:id
func GetProduct(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的商品ID",
		})
		return
	}

	var product Product
	if err := db.Debug().First(&product, uint(productID)).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "商品不存在",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "查询成功",
		Data:    product,
	})
}

// GetOrders 查询所有订单
// GET /orders
func GetOrders(c *gin.Context) {
	var orders []Order
	if err := db.Debug().
		Preload("User").
		Preload("OrderItems").
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "查询成功",
		Data:    orders,
	})
}

// GetOrder 查询单个订单
// GET /orders/:id
func GetOrder(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的订单ID",
		})
		return
	}

	var order Order
	if err := db.Debug().
		Preload("User").
		Preload("Address").
		Preload("OrderItems").
		Preload("OrderItems.Product").
		First(&order, uint(orderID)).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "订单不存在",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "查询成功",
		Data:    order,
	})
}

// SeedData 插入测试数据接口
// POST /seed
func SeedData(c *gin.Context) {
	if err := seedData(); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "测试数据插入成功",
	})
}

