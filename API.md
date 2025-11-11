# API 文档

## 基础信息

- 基础 URL: `http://localhost:8080`
- 默认端口: `8080` (可通过环境变量 `PORT` 修改)

## 统一响应格式

```json
{
  "code": 200,
  "message": "查询成功",
  "data": {}
}
```

## API 端点

### 健康检查

#### GET /health
检查服务是否正常运行

**响应示例:**
```json
{
  "status": "ok",
  "message": "服务运行正常"
}
```

---

### 用户相关 API

#### GET /api/users
查询所有用户

**响应示例:**
```json
{
  "code": 200,
  "message": "查询成功",
  "data": [
    {
      "id": 1,
      "username": "zhangsan",
      "phone": "13800138001",
      "email": "zhangsan@example.com",
      ...
    }
  ]
}
```

#### GET /api/users/:id/orders
查询指定用户的订单（包含订单明细）

**路径参数:**
- `id` (uint): 用户ID

**示例:**
```
GET /api/users/1/orders
```

**响应示例:**
```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "id": 1,
    "username": "zhangsan",
    "orders": [
      {
        "id": 1,
        "order_no": "ORD20251111...",
        "order_items": [...]
      }
    ]
  }
}
```

#### GET /api/users/:id/orders/products
查询用户的所有订单及每个订单的商品（多对多关系演示）

**路径参数:**
- `id` (uint): 用户ID

**示例:**
```
GET /api/users/1/orders/products
```

**响应示例:**
```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "id": 1,
    "username": "zhangsan",
    "orders": [
      {
        "id": 1,
        "order_no": "ORD20251111...",
        "order_items": [
          {
            "id": 1,
            "product_name": "iPhone 15",
            "quantity": 1,
            "product": {
              "id": 1,
              "name": "iPhone 15",
              "stock": 100
            }
          }
        ]
      }
    ]
  }
}
```

---

### 商品相关 API

#### GET /api/products
查询所有商品

**响应示例:**
```json
{
  "code": 200,
  "message": "查询成功",
  "data": [
    {
      "id": 1,
      "product_no": "PROD20251111...",
      "name": "iPhone 15 Pro",
      "price": 7999.00,
      "stock": 100,
      ...
    }
  ]
}
```

#### GET /api/products/:id
查询单个商品

**路径参数:**
- `id` (uint): 商品ID

**示例:**
```
GET /api/products/1
```

#### GET /api/products/:id/orders
查询商品被哪些订单购买（多对多关系演示）

**路径参数:**
- `id` (uint): 商品ID

**示例:**
```
GET /api/products/1/orders
```

**响应示例:**
```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "id": 1,
    "name": "iPhone 15 Pro",
    "order_items": [
      {
        "id": 1,
        "order": {
          "id": 1,
          "order_no": "ORD20251111...",
          "user": {
            "id": 1,
            "username": "zhangsan"
          }
        }
      }
    ]
  }
}
```

#### GET /api/products/:id/stats
统计商品的销售情况

**路径参数:**
- `id` (uint): 商品ID

**示例:**
```
GET /api/products/1/stats
```

**响应示例:**
```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "product": {...},
    "total_quantity": 10,
    "total_amount": 79990.00,
    "order_count": 5,
    "average_amount": 15998.00
  }
}
```

---

### 订单相关 API

#### GET /api/orders
查询所有订单

**响应示例:**
```json
{
  "code": 200,
  "message": "查询成功",
  "data": [
    {
      "id": 1,
      "order_no": "ORD20251111...",
      "user_id": 1,
      "total_amount": 7999.00,
      "order_items": [...]
    }
  ]
}
```

#### GET /api/orders/:id
查询单个订单详情

**路径参数:**
- `id` (uint): 订单ID

**示例:**
```
GET /api/orders/1
```

**响应示例:**
```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "id": 1,
    "order_no": "ORD20251111...",
    "user": {...},
    "address": {...},
    "order_items": [
      {
        "id": 1,
        "product_name": "iPhone 15 Pro",
        "quantity": 1,
        "product": {
          "id": 1,
          "name": "iPhone 15 Pro",
          "stock": 100
        }
      }
    ]
  }
}
```

#### GET /api/orders/:id/products
查询订单包含哪些商品（多对多关系演示）

**路径参数:**
- `id` (uint): 订单ID

**示例:**
```
GET /api/orders/1/products
```

**响应示例:**
```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "id": 1,
    "order_no": "ORD20251111...",
    "order_items": [
      {
        "id": 1,
        "product_name": "iPhone 15 Pro",
        "quantity": 1,
        "product": {
          "id": 1,
          "name": "iPhone 15 Pro",
          "stock": 100
        }
      }
    ]
  }
}
```

---

## 错误响应

### 400 Bad Request
```json
{
  "code": 400,
  "message": "无效的用户ID"
}
```

### 404 Not Found
```json
{
  "code": 404,
  "message": "用户不存在"
}
```

### 500 Internal Server Error
```json
{
  "code": 500,
  "message": "查询用户数据失败: ..."
}
```

---

## 使用示例

### 使用 curl

```bash
# 查询所有用户
curl http://localhost:8080/api/users

# 查询用户ID为1的订单
curl http://localhost:8080/api/users/1/orders

# 查询商品ID为1的订单
curl http://localhost:8080/api/products/1/orders

# 查询商品ID为1的销售统计
curl http://localhost:8080/api/products/1/stats
```

### 使用浏览器

直接在浏览器中访问：
- `http://localhost:8080/api/users`
- `http://localhost:8080/api/users/1/orders`
- `http://localhost:8080/api/products/1/orders`

### 使用 Postman

1. 创建新的 GET 请求
2. 输入 URL: `http://localhost:8080/api/users`
3. 点击 Send

---

## 环境变量

可以通过环境变量配置服务：

- `PORT`: 服务端口（默认: 8080）
- `GIN_MODE`: Gin 模式（debug/release/test，默认: debug）
- `DB_HOST`: 数据库主机（默认: 127.0.0.1）
- `DB_PORT`: 数据库端口（默认: 3306）
- `DB_USER`: 数据库用户（默认: root）
- `DB_PASSWORD`: 数据库密码（默认: mima123）
- `DB_NAME`: 数据库名称（默认: table_design）

---

## 启动服务

```bash
# 直接运行
go run .

# 或者编译后运行
go build .
./DataBaseDesign

# 指定端口运行
PORT=3000 go run .
```

