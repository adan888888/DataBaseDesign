package main

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表
type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement;comment:用户ID" json:"id"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex;not null;comment:用户名" json:"username"`
	Phone     string         `gorm:"type:varchar(20);uniqueIndex;comment:手机号" json:"phone"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;comment:邮箱" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null;comment:密码(加密)" json:"-"`
	Nickname  string         `gorm:"type:varchar(50);comment:昵称" json:"nickname"`
	Avatar    string         `gorm:"type:varchar(255);comment:头像URL" json:"avatar"`
	Status    int8           `gorm:"type:tinyint;default:1;comment:状态(1:正常 0:禁用)" json:"status"`
	CreatedAt time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"-"`

	// 关联关系
	Addresses []Address `gorm:"foreignKey:UserID;references:ID" json:"addresses,omitempty"`
	Orders    []Order   `gorm:"foreignKey:UserID;references:ID" json:"orders,omitempty"`
}

// Address 地址表
type Address struct {
	ID            uint           `gorm:"primaryKey;autoIncrement;comment:地址ID" json:"id"`
	UserID        uint           `gorm:"not null;index;comment:用户ID" json:"user_id"`
	ReceiverName  string         `gorm:"type:varchar(50);not null;comment:收货人姓名" json:"receiver_name"`
	ReceiverPhone string         `gorm:"type:varchar(20);not null;comment:收货人电话" json:"receiver_phone"`
	Province      string         `gorm:"type:varchar(50);not null;comment:省份" json:"province"`
	City          string         `gorm:"type:varchar(50);not null;comment:城市" json:"city"`
	District      string         `gorm:"type:varchar(50);not null;comment:区/县" json:"district"`
	Detail        string         `gorm:"type:varchar(255);not null;comment:详细地址" json:"detail"`
	PostalCode    string         `gorm:"type:varchar(10);comment:邮政编码" json:"postal_code"`
	IsDefault     bool           `gorm:"type:tinyint(1);default:0;comment:是否默认地址(1:是 0:否)" json:"is_default"`
	CreatedAt     time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index;comment:删除时间" json:"-"`

	// 关联关系
	User   User    `gorm:"foreignKey:UserID;references:ID" json:"-"`    // 隐藏反向关联，避免 JSON 输出冗余
	Orders []Order `gorm:"foreignKey:AddressID;references:ID" json:"-"` // 隐藏反向关联，避免 JSON 输出冗余
}

// Order 订单表
type Order struct {
	ID             uint           `gorm:"primaryKey;autoIncrement;comment:订单ID" json:"id"`
	OrderNo        string         `gorm:"type:varchar(32);uniqueIndex;not null;comment:订单号" json:"order_no"`
	UserID         uint           `gorm:"not null;index;comment:用户ID" json:"user_id"`
	AddressID      uint           `gorm:"not null;index;comment:收货地址ID" json:"address_id"`
	TotalAmount    float64        `gorm:"type:decimal(10,2);not null;default:0.00;comment:订单总金额" json:"total_amount"`
	DiscountAmount float64        `gorm:"type:decimal(10,2);default:0.00;comment:优惠金额" json:"discount_amount"`
	PayAmount      float64        `gorm:"type:decimal(10,2);not null;default:0.00;comment:实付金额" json:"pay_amount"`
	Status         int8           `gorm:"type:tinyint;default:0;index;comment:订单状态(0:待支付 1:已支付 2:已发货 3:已完成 4:已取消)" json:"status"`
	PayMethod      string         `gorm:"type:varchar(20);comment:支付方式" json:"pay_method"`
	PayTime        *time.Time     `gorm:"comment:支付时间" json:"pay_time"`
	ShipTime       *time.Time     `gorm:"comment:发货时间" json:"ship_time"`
	CompleteTime   *time.Time     `gorm:"comment:完成时间" json:"complete_time"`
	Remark         string         `gorm:"type:varchar(500);comment:订单备注" json:"remark"`
	CreatedAt      time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index;comment:删除时间" json:"-"`

	// 关联关系
	User       User        `gorm:"foreignKey:UserID;references:ID" json:"-"`                      // 隐藏反向关联，避免 JSON 输出冗余
	Address    Address     `gorm:"foreignKey:AddressID;references:ID" json:"-"`                   // 隐藏反向关联，避免 JSON 输出冗余
	OrderItems []OrderItem `gorm:"foreignKey:OrderID;references:ID" json:"order_items,omitempty"` //如果值为零值，则不输出（对结构体无效）
}

// Product 商品表
type Product struct {
	ID          uint           `gorm:"primaryKey;autoIncrement;comment:商品ID" json:"id"`
	ProductNo   string         `gorm:"type:varchar(50);uniqueIndex;not null;comment:商品编号" json:"product_no"`
	Name        string         `gorm:"type:varchar(200);not null;index;comment:商品名称" json:"name"`
	Description string         `gorm:"type:text;comment:商品描述" json:"description"`
	CategoryID  uint           `gorm:"index;comment:分类ID" json:"category_id"`
	Price       float64        `gorm:"type:decimal(10,2);not null;default:0.00;comment:商品价格" json:"price"`
	Stock       int            `gorm:"type:int;default:0;comment:库存数量" json:"stock"`
	Sales       int            `gorm:"type:int;default:0;comment:销量" json:"sales"`
	Image       string         `gorm:"type:varchar(500);comment:商品主图" json:"image"`
	Images      string         `gorm:"type:text;comment:商品图片(JSON数组)" json:"images"`
	Status      int8           `gorm:"type:tinyint;default:1;index;comment:状态(1:上架 0:下架)" json:"status"`
	Sort        int            `gorm:"type:int;default:0;comment:排序" json:"sort"`
	CreatedAt   time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index;comment:删除时间" json:"-"`

	// 关联关系
	OrderItems []OrderItem `gorm:"foreignKey:ProductID;references:ID" json:"order_items,omitempty"`
}

// OrderItem 订单商品明细表
type OrderItem struct {
	ID           uint           `gorm:"primaryKey;autoIncrement;comment:明细ID" json:"id"`
	OrderID      uint           `gorm:"not null;index;comment:订单ID" json:"order_id"`
	ProductID    uint           `gorm:"not null;index;comment:商品ID" json:"product_id"`
	ProductName  string         `gorm:"type:varchar(200);not null;comment:商品名称(快照)" json:"product_name"`
	ProductImage string         `gorm:"type:varchar(500);comment:商品图片(快照)" json:"product_image"`
	Price        float64        `gorm:"type:decimal(10,2);not null;comment:商品单价(快照)" json:"price"`
	Quantity     int            `gorm:"type:int;not null;default:1;comment:购买数量" json:"quantity"`
	Subtotal     float64        `gorm:"type:decimal(10,2);not null;comment:小计金额" json:"subtotal"`
	CreatedAt    time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index;comment:删除时间" json:"-"`

	// 关联关系
	Order   Order   `gorm:"foreignKey:OrderID;references:ID" json:"-"`   // 隐藏反向关联，避免 JSON 输出冗余
	Product Product `gorm:"foreignKey:ProductID;references:ID" json:"-"` // 隐藏反向关联，避免 JSON 输出冗余
}
