# 个人财务管理 API

基于 Go 和 Gin 框架开发的个人财务管理 API 服务。

## 功能特性

- **用户认证**：JWT 认证，用户数据隔离
- **账户管理**：多币种账户的创建、修改、删除
- **分类管理**：收支分类树，支持一级和二级分类
- **交易记录**：收入/支出记录，支持按时间筛选
- **月度统计**：按账户和分类维度汇总当月收支
- **预算管理**：分类月度预算设置和超支提醒
- **数据导出**：按时间范围导出交易记录为 CSV

## 技术栈

- Go 1.21+
- Gin Web 框架
- GORM ORM
- MySQL 数据库
- JWT 认证

## 项目结构

```
finance-api/
├── config/          # 配置管理
├── controllers/     # 控制器层
├── middleware/      # 中间件
├── models/          # 数据模型
├── routes/          # 路由配置
├── services/        # 业务逻辑层
├── utils/           # 工具函数
└── main.go         # 主程序入口
```

## 快速开始

### 1. 环境准备

- Go 1.21+
- MySQL 5.7+

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置环境变量

复制 `.env.example` 为 `.env` 并修改配置：

```bash
cp .env.example .env
```

### 4. 创建数据库

```sql
CREATE DATABASE finance_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 5. 运行服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

## API 文档

### 认证接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /api/auth/register | 用户注册 | 否 |
| POST | /api/auth/login | 用户登录 | 否 |
| GET | /api/auth/profile | 获取用户信息 | 是 |

### 账户管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/accounts | 创建账户 |
| GET | /api/accounts | 获取账户列表 |
| GET | /api/accounts/:id | 获取账户详情 |
| PUT | /api/accounts/:id | 更新账户 |
| DELETE | /api/accounts/:id | 删除账户 |

### 分类管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/categories | 创建分类 |
| GET | /api/categories?type=income/expense | 获取分类列表 |
| GET | /api/categories/:id | 获取分类详情 |
| PUT | /api/categories/:id | 更新分类 |
| DELETE | /api/categories/:id | 删除分类 |

### 交易记录

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/transactions | 创建交易 |
| GET | /api/transactions | 获取交易列表（支持分页筛选） |
| GET | /api/transactions/:id | 获取交易详情 |
| PUT | /api/transactions/:id | 更新交易 |
| DELETE | /api/transactions/:id | 删除交易 |

### 统计功能

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/statistics/monthly?month=YYYY-MM | 获取月度统计 |

### 预算管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/budgets | 创建预算 |
| GET | /api/budgets?month=YYYY-MM | 获取预算列表 |
| GET | /api/budgets/:id | 获取预算详情 |
| PUT | /api/budgets/:id | 更新预算 |
| DELETE | /api/budgets/:id | 删除预算 |

### 数据导出

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/export/transactions?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD | 导出交易CSV |

## 请求示例

### 用户注册

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456","email":"test@example.com"}'
```

### 用户登录

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456"}'
```

### 创建账户

```bash
curl -X POST http://localhost:8080/api/accounts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {token}" \
  -d '{"name":"现金","currency":"CNY","initial_balance":1000,"remark":"日常使用"}'
```

## 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

- `code`: 0 表示成功，非 0 表示错误
- `message`: 错误信息
- `data`: 响应数据

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |
