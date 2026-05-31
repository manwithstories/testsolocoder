# 养蜂管理与蜂蜜交易平台

一个基于 Go Gin + React + TypeScript 的在线养蜂管理与蜂蜜交易平台，为蜂农、蜂蜜加工厂和蜂蜜爱好者提供服务。

## 功能特性

### 1. 蜂场管理
- 蜂箱创建与位置坐标标注
- 蜂群健康状况记录
- 产蜜周期管理
- 蜂箱按区域或品种分组管理

### 2. 蜂群健康监控
- 蜂王状态记录
- 工蜂数量变化追踪
- 病虫害自动预警通知
- 季节性管理建议

### 3. 蜂蜜采收
- 采收记录管理（来源蜂箱、日期、产量、品种）
- 自动更新蜂箱状态
- 自动更新库存数量

### 4. 库存管理
- 蜂蜜按品种和批次分类
- 检测报告绑定
- 保质期管理
- 库存预警与临期预警

### 5. 交易市场
- 蜂农发布蜂蜜信息
- 买家下单购买
- 物流跟踪
- 交易互评与信誉系统

### 6. 检测预约
- 预约检测机构
- 检测报告与批次绑定
- 检测结果决定上架等级

### 7. 蜂农社区
- 发帖分享养蜂经验
- 讨论问题
- 标签分类（病虫害防治、采收技巧等）

### 8. 数据分析
- 蜂场产量统计
- 销售趋势分析
- 病虫害发生率分析
- 报表导出

## 技术栈

### 后端
- Go 1.21+
- Gin Web Framework
- GORM ORM
- PostgreSQL
- JWT 认证
- bcrypt 密码加密

### 前端
- React 18
- TypeScript
- Ant Design
- ECharts 图表
- Zustand 状态管理
- Axios HTTP 客户端
- React Router

## 项目结构

```
.
├── backend/                 # 后端代码
│   ├── config/             # 配置管理
│   ├── database/           # 数据库连接
│   ├── handlers/           # 业务处理器
│   ├── middleware/          # 中间件
│   ├── models/             # 数据模型
│   ├── routes/             # 路由定义
│   ├── utils/              # 工具函数
│   └── main.go             # 主入口
├── frontend/               # 前端代码
│   ├── src/
│   │   ├── api/            # API 调用
│   │   ├── components/     # 组件
│   │   ├── pages/          # 页面
│   │   ├── store/          # 状态管理
│   │   ├── types/          # TypeScript 类型
│   │   └── utils/          # 工具函数
│   ├── package.json
│   └── vite.config.ts
├── database/               # 数据库相关
│   └── schema.sql          # 数据库 Schema
└── README.md
```

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- PostgreSQL 12+

### 后端启动

```bash
cd backend

# 安装依赖
go mod download

# 配置环境变量（修改 .env 文件）
cp .env.example .env

# 启动服务
go run main.go
```

后端服务将在 `http://localhost:8080` 启动。

### 前端启动

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端服务将在 `http://localhost:3000` 启动。

### 数据库初始化

```bash
# 创建数据库
createdb beehive_platform

# 导入 Schema
psql -d beehive_platform -f database/schema.sql
```

或者，首次启动后端时会自动执行数据库迁移。

## 用户角色

- **蜂农 (beekeeper)**: 管理蜂箱、记录健康状况、采收蜂蜜、管理库存、发布商品、处理订单
- **买家 (buyer)**: 浏览商品、下单购买、评价订单
- **检测机构 (inspector)**: 接受检测预约、提交检测结果

## API 文档

### 认证
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `GET /api/v1/auth/profile` - 获取用户信息
- `PUT /api/v1/auth/profile` - 更新用户信息

### 蜂箱管理
- `POST /api/v1/beehives` - 创建蜂箱
- `GET /api/v1/beehives` - 蜂箱列表
- `GET /api/v1/beehives/:id` - 蜂箱详情
- `PUT /api/v1/beehives/:id` - 更新蜂箱
- `DELETE /api/v1/beehives/:id` - 删除蜂箱

### 健康记录
- `POST /api/v1/health-records` - 创建健康记录
- `GET /api/v1/health-records` - 健康记录列表
- `GET /api/v1/disease-warnings` - 病害预警
- `GET /api/v1/seasonal-tips` - 季节性建议

### 采收管理
- `POST /api/v1/harvests` - 创建采收记录
- `GET /api/v1/harvests` - 采收记录列表
- `GET /api/v1/harvests/:id` - 采收记录详情

### 库存管理
- `GET /api/v1/inventory` - 库存列表
- `GET /api/v1/inventory/:id` - 库存详情
- `PUT /api/v1/inventory/:id` - 更新库存
- `GET /api/v1/inventory/alerts` - 库存预警

### 商品管理
- `POST /api/v1/products` - 发布商品
- `GET /api/v1/products` - 商品列表
- `GET /api/v1/products/:id` - 商品详情
- `GET /api/v1/my-products` - 我的商品
- `PUT /api/v1/products/:id` - 更新商品
- `DELETE /api/v1/products/:id` - 删除商品

### 订单管理
- `POST /api/v1/orders` - 创建订单
- `GET /api/v1/orders` - 订单列表
- `GET /api/v1/orders/:id` - 订单详情
- `PUT /api/v1/orders/:id/pay` - 支付订单
- `PUT /api/v1/orders/:id/ship` - 发货
- `PUT /api/v1/orders/:id/deliver` - 确认收货
- `PUT /api/v1/orders/:id/complete` - 完成订单
- `PUT /api/v1/orders/:id/cancel` - 取消订单
- `PUT /api/v1/orders/:id/rate` - 评价订单

### 检测管理
- `POST /api/v1/inspections` - 预约检测
- `GET /api/v1/inspections` - 检测列表
- `GET /api/v1/inspections/:id` - 检测详情
- `PUT /api/v1/inspections/:id/assign` - 接受预约
- `PUT /api/v1/inspections/:id/result` - 提交结果
- `PUT /api/v1/inspections/:id/cancel` - 取消预约

### 社区
- `POST /api/v1/posts` - 发布帖子
- `GET /api/v1/posts` - 帖子列表
- `GET /api/v1/posts/:id` - 帖子详情
- `PUT /api/v1/posts/:id` - 更新帖子
- `DELETE /api/v1/posts/:id` - 删除帖子
- `POST /api/v1/posts/:id/like` - 点赞帖子
- `POST /api/v1/comments` - 发表评论
- `GET /api/v1/posts/:id/comments` - 评论列表

### 数据分析
- `GET /api/v1/analytics/overview` - 概览统计
- `GET /api/v1/analytics/production` - 产量统计
- `GET /api/v1/analytics/disease` - 病害统计
- `GET /api/v1/analytics/sales` - 销售统计
- `GET /api/v1/analytics/export` - 导出报表

## 配置说明

后端配置通过 `.env` 文件管理：

```env
# 服务器配置
SERVER_PORT=8080
SERVER_READ_TIMEOUT=10
SERVER_WRITE_TIMEOUT=10

# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=beehive_platform
DB_SSLMODE=disable

# JWT 配置
JWT_SECRET=your-secret-key
JWT_EXPIRE_HOUR=24

# 上传配置
UPLOAD_PATH=./uploads
UPLOAD_MAX_SIZE=50

# 日志配置
LOG_PATH=./logs
LOG_MAX_SIZE=100
LOG_MAX_BACKUPS=3
LOG_MAX_AGE=28
```

## 数据校验

系统对以下数据进行边界校验：
- 蜂箱坐标范围：纬度 -90~90，经度 -180~180
- 蜂蜜保质期：默认采收后 2 年
- 库存数量：非负数
- 订单数量：大于 0 且不超过库存
- 价格：大于 0

## 事务处理

订单创建使用数据库事务确保数据一致性：
1. 检查库存是否充足
2. 创建订单
3. 扣减商品库存
4. 扣减库存数量
5. 发送通知

## 日志记录

系统记录以下日志：
- 操作日志：用户操作记录
- 错误日志：异常情况记录
- 访问日志：HTTP 请求日志

日志文件存储在 `logs/` 目录，按大小和日期自动切割。

## 许可证

MIT License
