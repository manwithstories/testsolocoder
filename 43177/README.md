# 在线专业维修服务平台

一个基于 Go (Gin) + Vue3 + TypeScript + Redis 的全栈维修服务平台，支持家电维修、数码维修、汽车维修三大服务领域。

## 项目结构

```
.
├── backend/                 # 后端 Go 项目
│   ├── cmd/                 # 入口文件
│   │   └── main.go
│   ├── internal/            # 内部代码
│   │   ├── models/          # 数据模型
│   │   ├── handlers/        # 处理器
│   │   ├── middleware/      # 中间件
│   │   ├── routes/          # 路由
│   │   ├── services/        # 服务层
│   │   └── utils/           # 工具函数
│   ├── pkg/                 # 公共包
│   │   ├── config/          # 配置
│   │   ├── logger/          # 日志
│   │   └── redis/           # Redis
│   └── go.mod
└── frontend/                # 前端 Vue3 项目
    ├── src/
    │   ├── api/             # API 请求
    │   ├── components/      # 组件
    │   ├── router/          # 路由
    │   ├── store/           # 状态管理
    │   ├── types/           # 类型定义
    │   ├── views/           # 页面
    │   └── assets/          # 静态资源
    └── package.json
```

## 核心功能

### 1. 用户管理
- 支持客户、技师、平台管理员三种角色注册登录
- JWT 认证机制
- 技师资质证书审核
- 用户评分和历史订单

### 2. 维修项目管理
- 家电维修、数码维修、汽车维修三大分类
- 标准报价区间和预计维修时长
- 技师可给出具体报价

### 3. 工单管理
- 客户提交维修需求
- 系统自动匹配附近技师
- 工单状态流转（待分配→待接单→已接单→已到达→维修中→已完成）
- 工单取消和退款处理
- 超时自动关闭

### 4. 评价管理
- 客户评分和评价
- 技师回复评价
- 差评触发平台介入机制

### 5. 配件管理
- 技师申请常用配件库存
- 平台统一采购配送
- 配件使用关联工单记录
- 库存不足自动预警

### 6. 财务管理
- 技师提现申请
- 平台抽成（可配置15%）
- 月度收入报表
- 技师绩效分析

## 快速开始

### 后端

```bash
cd backend

# 安装依赖
go mod tidy

# 运行服务
go run cmd/main.go
```

后端服务运行在 `http://localhost:8080`

### 前端

```bash
cd frontend

# 安装依赖
npm install

# 开发模式
npm run dev

# 构建
npm run build
```

前端开发服务运行在 `http://localhost:3000`

## API 文档

### 认证相关
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录

### 用户相关
- `GET /api/user/profile` - 获取用户信息
- `PUT /api/user/profile` - 更新用户信息
- `POST /api/user/certificate` - 提交技师证书

### 工单相关
- `POST /api/orders` - 创建工单
- `GET /api/orders` - 获取工单列表
- `GET /api/orders/:id` - 获取工单详情
- `POST /api/orders/:id/accept` - 接单
- `POST /api/orders/:id/arrive` - 到达现场
- `POST /api/orders/:id/start` - 开始维修
- `POST /api/orders/:id/complete` - 完工
- `POST /api/orders/:id/cancel` - 取消工单
- `POST /api/orders/:id/refund` - 申请退款

### 管理后台
- `GET /api/admin/dashboard` - 数据概览
- `GET /api/admin/users` - 用户列表
- `POST /api/admin/technicians/:id/verify` - 技师审核
- `POST /api/admin/refunds/:id/approve` - 退款通过
- `POST /api/admin/refunds/:id/reject` - 退款拒绝

## 数据库模型

- User - 用户表
- TechnicianProfile - 技师资料表
- Category - 服务分类表
- ServiceItem - 服务项目表
- Order - 工单表
- OrderLog - 工单日志表
- Review - 评价表
- Part - 配件表
- PartRequest - 配件申请表
- PartRequestItem - 配件申请明细表
- PartUsage - 配件使用记录表
- WithdrawRequest - 提现申请表
- Transaction - 交易记录表
- MonthlyReport - 月度报表表

## 技术栈

### 后端
- Go 1.21
- Gin Web Framework
- GORM (SQLite)
- Redis
- JWT 认证
- Logrus 日志

### 前端
- Vue 3
- TypeScript
- Vite
- Pinia
- Vue Router
- Element Plus
- Axios
- ECharts

## 默认账号

系统初始化时会自动创建以下管理员账号（需在代码中设置）：

- 用户名: admin
- 密码: admin123
- 角色: 管理员

## 注意事项

1. Redis 为可选依赖，未连接时系统仍可运行
2. 数据库使用 SQLite，首次运行会自动创建表结构
3. 平台抽成比例默认为 15%，可在配置中修改
4. 工单超时自动关闭时间为 24 小时
