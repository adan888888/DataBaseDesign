# 代理键 vs 自然键：ID 和学号的选择

## 目录
- [问题背景](#问题背景)
- [基本概念](#基本概念)
- [两种方案对比](#两种方案对比)
- [实际应用场景](#实际应用场景)
- [最佳实践建议](#最佳实践建议)
- [示例说明](#示例说明)

---

## 问题背景

在设计数据库表时，经常会遇到这样的问题：

**如果表已经有自动增长的 ID 作为主键，还需要像学号这样的业务字段吗？**

例如：
- 学生表：有 `id`（自增主键），还需要 `学号` 吗？
- 订单表：有 `id`（自增主键），还需要 `订单号` 吗？
- 用户表：有 `id`（自增主键），还需要 `用户名` 吗？

---

## 基本概念

### 代理键（Surrogate Key）

**定义**：数据库系统自动生成的、没有业务意义的唯一标识符。

**特点**：
- ✅ 自动生成（如自增ID、UUID）
- ✅ 没有业务含义
- ✅ 稳定不变
- ✅ 简单高效

**常见类型**：
- 自增整数（AUTO_INCREMENT）
- UUID（通用唯一标识符）
- GUID（全局唯一标识符）

**示例**：
```sql
CREATE TABLE student (
    id INT PRIMARY KEY AUTO_INCREMENT,  -- 代理键
    name VARCHAR(50),
    age INT
);
```

### 自然键（Natural Key）

**定义**：具有业务意义的、能够唯一标识实体的字段。

**特点**：
- ✅ 有业务含义
- ✅ 用户可见和理解
- ✅ 可能变化（业务规则改变）
- ✅ 可能包含多个字段（复合键）

**示例**：
```sql
CREATE TABLE student (
    student_no VARCHAR(20) PRIMARY KEY,  -- 自然键（学号）
    name VARCHAR(50),
    age INT
);
```

---

## 两种方案对比

### 方案一：只用 ID（代理键）

```sql
CREATE TABLE student (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50),
    age INT
);
```

**优点**：
- ✅ **性能好**：整数类型，索引效率高
- ✅ **稳定**：不会因为业务规则改变而改变
- ✅ **简单**：不需要考虑业务规则
- ✅ **高效**：自增ID，插入速度快

**缺点**：
- ❌ **无业务含义**：用户无法理解ID的含义
- ❌ **不友好**：用户查询时需要记住ID
- ❌ **不直观**：URL、日志中显示ID不够友好

### 方案二：只用学号（自然键）

```sql
CREATE TABLE student (
    student_no VARCHAR(20) PRIMARY KEY,
    name VARCHAR(50),
    age INT
);
```

**优点**：
- ✅ **有业务含义**：学号对用户有意义
- ✅ **用户友好**：用户可以直接使用学号查询
- ✅ **直观**：URL、日志中显示学号更清晰

**缺点**：
- ❌ **可能变化**：业务规则改变时可能需要修改
- ❌ **性能略差**：字符串类型，索引效率略低于整数
- ❌ **复杂度高**：需要保证学号的唯一性和规则

### 方案三：ID + 学号（推荐）⭐

```sql
CREATE TABLE student (
    id INT PRIMARY KEY AUTO_INCREMENT,      -- 代理键（主键）
    student_no VARCHAR(20) UNIQUE NOT NULL,  -- 自然键（唯一索引）
    name VARCHAR(50),
    age INT
);
```

**优点**：
- ✅ **兼顾性能**：ID作为主键，性能最优
- ✅ **业务友好**：学号作为业务标识，用户友好
- ✅ **灵活性强**：可以同时使用两种方式
- ✅ **稳定性好**：ID不变，学号可以修改

**缺点**：
- ⚠️ **存储略多**：多一个字段
- ⚠️ **需要维护**：需要保证学号的唯一性

---

## 实际应用场景

### 场景一：学生管理系统

**推荐方案**：ID（主键）+ 学号（唯一索引）

```sql
CREATE TABLE student (
    id INT PRIMARY KEY AUTO_INCREMENT,
    student_no VARCHAR(20) UNIQUE NOT NULL COMMENT '学号',
    name VARCHAR(50) NOT NULL,
    email VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 查询示例
SELECT * FROM student WHERE id = 1;           -- 内部查询（快）
SELECT * FROM student WHERE student_no = '2024001';  -- 用户查询（友好）
```

**原因**：
- 学号是业务标识，用户需要看到和使用
- ID作为主键，保证性能和稳定性
- 学号可以修改（如转学、重新分配），但ID不变

### 场景二：订单系统

**推荐方案**：ID（主键）+ 订单号（唯一索引）

```sql
CREATE TABLE orders (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_no VARCHAR(32) UNIQUE NOT NULL COMMENT '订单号',
    user_id INT NOT NULL,
    amount DECIMAL(10,2),
    status VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**原因**：
- 订单号需要展示给用户（如：ORD20241110001）
- ID用于内部关联和性能优化
- 订单号有业务规则（如：日期+序号）

### 场景三：用户系统

**推荐方案**：ID（主键）+ 用户名（唯一索引）+ 手机号（唯一索引）

```sql
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL COMMENT '用户名',
    phone VARCHAR(20) UNIQUE COMMENT '手机号',
    email VARCHAR(100) UNIQUE COMMENT '邮箱',
    password VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**原因**：
- 用户可以通过用户名、手机号、邮箱登录
- ID作为主键，用于关联其他表
- 多个业务标识字段都需要唯一性约束

---

## 最佳实践建议

### ✅ 推荐做法

1. **主键使用代理键（ID）**
   - 使用自增整数或UUID
   - 保证性能和稳定性
   - 用于表间关联

2. **业务标识使用自然键（学号/订单号等）**
   - 添加唯一索引（UNIQUE）
   - 用于用户查询和展示
   - 可以修改，但不频繁

3. **建立唯一索引**
   - 对业务标识字段建立唯一索引
   - 保证数据唯一性
   - 提高查询性能

### ❌ 不推荐做法

1. **只用自然键作为主键**
   ```sql
   -- 不推荐：学号作为主键
   CREATE TABLE student (
       student_no VARCHAR(20) PRIMARY KEY,  -- 性能较差
       ...
   );
   ```

2. **业务标识字段没有唯一约束**
   ```sql
   -- 不推荐：学号可以重复
   CREATE TABLE student (
       id INT PRIMARY KEY,
       student_no VARCHAR(20),  -- 缺少UNIQUE约束
       ...
   );
   ```

3. **过度使用代理键**
   ```sql
   -- 不推荐：简单的配置表不需要ID
   CREATE TABLE config (
       id INT PRIMARY KEY AUTO_INCREMENT,
       key VARCHAR(50) UNIQUE,  -- key本身就可以作为主键
       value TEXT
   );
   ```

---

## 示例说明

### 完整示例：学生管理系统

```sql
-- 学生表
CREATE TABLE student (
    id INT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    student_no VARCHAR(20) UNIQUE NOT NULL COMMENT '学号',
    name VARCHAR(50) NOT NULL COMMENT '姓名',
    gender ENUM('男', '女') COMMENT '性别',
    birth_date DATE COMMENT '出生日期',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_student_no (student_no),  -- 学号索引
    INDEX idx_name (name)               -- 姓名索引
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学生表';

-- 选课表（关联表）
CREATE TABLE course_selection (
    id INT PRIMARY KEY AUTO_INCREMENT,
    student_id INT NOT NULL COMMENT '学生ID（关联student.id）',
    course_id INT NOT NULL COMMENT '课程ID',
    score DECIMAL(5,2) COMMENT '成绩',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (student_id) REFERENCES student(id),
    INDEX idx_student_id (student_id),
    INDEX idx_course_id (course_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='选课表';
```

### 使用场景对比

#### 场景1：内部系统查询（使用ID）
```sql
-- 通过ID查询（性能最优）
SELECT * FROM student WHERE id = 1001;

-- 关联查询（使用ID）
SELECT s.name, cs.score 
FROM student s
JOIN course_selection cs ON s.id = cs.student_id
WHERE s.id = 1001;
```

#### 场景2：用户查询（使用学号）
```sql
-- 通过学号查询（用户友好）
SELECT * FROM student WHERE student_no = '2024001';

-- 用户登录后查询自己的信息
SELECT * FROM student WHERE student_no = '2024001';
```

#### 场景3：API接口设计

```go
// RESTful API 设计示例

// 使用ID（内部）
GET /api/v1/students/1001

// 使用学号（用户友好）
GET /api/v1/students/by-student-no/2024001

// 或者统一使用ID，但返回数据包含学号
GET /api/v1/students/1001
Response: {
    "id": 1001,
    "student_no": "2024001",
    "name": "张三",
    ...
}
```

---

## 总结

### 核心原则

1. **主键使用代理键（ID）**
   - ✅ 性能最优
   - ✅ 稳定不变
   - ✅ 用于表间关联

2. **业务标识使用自然键（学号/订单号等）**
   - ✅ 用户友好
   - ✅ 有业务含义
   - ✅ 添加唯一索引

3. **两者结合使用**
   - ✅ 兼顾性能和用户体验
   - ✅ 灵活性强
   - ✅ 最佳实践

### 回答原问题

**如果表已经有自动增长的 ID 作为主键，还需要学号吗？**

**答案：通常需要！**

**原因**：
- ID 作为主键，用于性能和关联
- 学号作为业务标识，用于用户查询和展示
- 两者结合使用，各司其职

**特殊情况**：
- 简单的配置表、日志表：可以只用ID
- 纯内部系统：可以只用ID
- 面向用户的系统：建议ID + 业务标识

---

## 常见问题 FAQ

### Q1: ID 和学号可以都作为主键吗？

**A:** 不可以。一个表只能有一个主键。建议：
- ID 作为主键（PRIMARY KEY）
- 学号作为唯一索引（UNIQUE INDEX）

### Q2: 如果学号会变化怎么办？

**A:** 这正是使用ID作为主键的优势：
- ID 不变，用于关联其他表
- 学号可以修改，不影响关联关系
- 修改学号时，只需更新学号字段

### Q3: 什么时候可以只用ID？

**A:** 以下情况可以只用ID：
- 纯内部系统，用户不需要看到标识
- 简单的配置表、日志表
- 临时表、中间表
- 性能要求极高的场景

### Q4: UUID 和自增ID如何选择？

**A:** 
- **自增ID**：性能好，但可能暴露业务信息（如订单数量）
- **UUID**：更安全，但性能略差，存储空间大
- 根据实际需求选择

---

*最后更新时间：2025年*

