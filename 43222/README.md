# 🌱 家庭菜园规划与园艺管理平台 (Garden Planner)

一个完整的在线家庭菜园规划与园艺管理平台，使用 Go (Gin) + React + TypeScript + PostgreSQL 构建。

## 📋 功能特性

### 1. 菜园地块管理
- 创建多个种植区域，记录土壤类型、光照条件、面积大小
- 支持自定义划分种植区块
- 绑定灌溉设备和传感器数据

### 2. 植物数据库
- 内置常见蔬菜水果的种植信息
- 包括生长周期、浇水频率、施肥需求、病虫害预防
- 支持搜索筛选和多条件组合查询

### 3. 种植日历
- 根据用户所在地区和季节智能推荐种植品种
- 设置浇水施肥提醒
- 支持农历和公历切换显示

### 4. 生长追踪
- 上传植物照片记录生长过程
- 系统自动生成时间轴
- 标记重要节点如发芽、开花、结果
- 导出种植报告

### 5. 病虫害诊断
- 上传病叶照片或描述症状
- 智能匹配可能的问题并给出防治建议
- 支持图片识别和文字描述两种方式

### 6. 社区交流
- 分享种植经验、提问求助
- 关注其他用户
- 支持图文帖子、评论和点赞

### 7. 种子交换市场
- 发布多余种子与他人交换
- 支持在线协商和线下交易确认
- 信用评价体系

### 8. 购物功能
- 对接种子商家和园艺用品店
- 直接购买种子、肥料、工具
- 支持购物车、订单管理、物流追踪

## 🏗️ 技术架构

### 后端 (Go + Gin)
```
backend/
├── config/          # 配置管理
├── database/        # 数据库连接和迁移
├── handlers/        # API处理器
│   ├── auth.go      # 认证
│   ├── plot.go      # 地块管理
│   ├── plant.go     # 植物数据库
│   ├── growth.go    # 生长追踪
│   ├── calendar.go  # 种植日历
│   ├── disease.go   # 病虫害诊断
│   ├── community.go # 社区交流
│   ├── exchange.go  # 种子交换
│   └── shop.go      # 购物功能
├── middleware/      # 中间件
│   ├── auth.go      # JWT认证
│   └── logging.go   # 日志记录
├── models/          # 数据模型
└── main.go          # 入口文件
```

### 前端 (React + TypeScript)
```
frontend/
├── src/
│   ├── api/          # API接口封装
│   ├── components/   # 组件
│   ├── layouts/      # 布局
│   ├── pages/        # 页面
│   │   ├── LoginPage.tsx
│   │   ├── RegisterPage.tsx
│   │   ├── HomePage.tsx
│   │   ├── PlotsPage.tsx
│   │   ├── PlantsPage.tsx
│   │   ├── CalendarPage.tsx
│   │   ├── GrowthPage.tsx
│   │   ├── DiseasePage.tsx
│   │   ├── CommunityPage.tsx
│   │   ├── PostDetailPage.tsx
│   │   ├── ExchangePage.tsx
│   │   ├── ShopPage.tsx
│   │   ├── CartPage.tsx
│   │   ├── OrdersPage.tsx
│   │   └── ProfilePage.tsx
│   ├── store/        # 状态管理 (Zustand)
│   ├── types/        # TypeScript类型定义
│   ├── App.tsx       # 应用入口
│   └── main.tsx      # React入口
```

## 🗄️ 数据库设计

主要数据表：
- **users**: 用户表
- **plots**: 地块表
- **plants**: 植物表
- **planting_records**: 种植记录表
- **growth_logs**: 生长日志表
- **posts**: 帖子表
- **comments**: 评论表
- **likes**: 点赞表
- **follows**: 关注表
- **disease_diagnoses**: 病虫害诊断表
- **seed_exchanges**: 种子交换表
- **exchange_offers**: 交换要约表
- **products**: 商品表
- **carts**: 购物车表
- **orders**: 订单表
- **order_items**: 订单项目表
- **calendar_events**: 日历事件表
- **operation_logs**: 操作日志表

## 🚀 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- PostgreSQL 13+

### 后端设置

```bash
cd backend

# 安装依赖
go mod download

# 配置数据库
# 编辑 .env 文件设置数据库连接信息

# 运行服务
go run main.go
```

后端服务将在 `http://localhost:8080` 启动

### 前端设置

```bash
cd frontend

# 安装依赖
npm install

# 开发模式运行
npm run dev
```

前端服务将在 `http://localhost:3000` 启动

## 🔐 API认证

所有需要认证的API请求需要在Header中添加：
```
Authorization: Bearer <token>
```

## 📝 API端点

### 认证
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录

### 用户
- `GET /api/v1/users/profile` - 获取个人资料
- `PUT /api/v1/users/profile` - 更新个人资料
- `PUT /api/v1/users/password` - 修改密码

### 地块管理
- `POST /api/v1/plots` - 创建地块
- `GET /api/v1/plots` - 获取地块列表
- `GET /api/v1/plots/:id` - 获取地块详情
- `PUT /api/v1/plots/:id` - 更新地块
- `DELETE /api/v1/plots/:id` - 删除地块

### 植物数据库
- `GET /api/v1/plants` - 获取植物列表
- `GET /api/v1/plants/:id` - 获取植物详情
- `POST /api/v1/plants` - 创建植物(专家/管理员)

### 种植记录
- `POST /api/v1/planting-records` - 创建种植记录
- `GET /api/v1/planting-records` - 获取种植记录
- `GET /api/v1/planting-records/:id` - 获取详情
- `PUT /api/v1/planting-records/:id` - 更新记录
- `DELETE /api/v1/planting-records/:id` - 删除记录

### 生长日志
- `POST /api/v1/growth-logs` - 创建生长日志
- `GET /api/v1/growth-logs` - 获取生长日志
- `PUT /api/v1/growth-logs/:id` - 更新日志
- `DELETE /api/v1/growth-logs/:id` - 删除日志

### 文件上传
- `POST /api/v1/upload` - 上传文件

### 种植日历
- `POST /api/v1/calendar/events` - 创建日历事件
- `GET /api/v1/calendar/events` - 获取日历事件
- `GET /api/v1/calendar/recommendations` - 获取种植推荐

### 病虫害诊断
- `POST /api/v1/disease-diagnosis` - 创建诊断
- `GET /api/v1/disease-diagnosis` - 获取诊断列表
- `GET /api/v1/disease-diagnosis/:id` - 获取诊断详情

### 社区
- `POST /api/v1/posts` - 创建帖子
- `GET /api/v1/posts` - 获取帖子列表
- `GET /api/v1/posts/:id` - 获取帖子详情
- `POST /api/v1/posts/:id/like` - 点赞帖子
- `POST /api/v1/posts/:id/comments` - 添加评论

### 种子交换
- `POST /api/v1/seed-exchanges` - 创建交换
- `GET /api/v1/seed-exchanges` - 获取交换列表
- `POST /api/v1/seed-exchanges/:id/offers` - 发起交换要约

### 商城
- `GET /api/v1/products` - 获取商品列表
- `POST /api/v1/cart` - 添加到购物车
- `GET /api/v1/cart` - 获取购物车
- `POST /api/v1/orders` - 创建订单
- `GET /api/v1/orders` - 获取订单列表

## 📄 许可证

MIT License
