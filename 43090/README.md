# 在线拍卖竞价管理系统

一个基于 Go Gin + Vue3 + TypeScript + Redis 的全栈在线拍卖竞价管理系统。

## 功能特性

### 后端功能 (Go + Gin)
- ✅ 用户认证与权限管理（JWT、角色权限）
- ✅ 拍卖品管理（增删改查、上下架、多图上传）
- ✅ 拍卖会管理（创建、编辑、竞价规则配置）
- ✅ 实时竞价功能（Redis缓存、并发控制、加价校验）
- ✅ 自动竞价代理（自动加价、上限设置）
- ✅ 消息通知系统（站内消息、邮件推送）
- ✅ 订单和支付管理（订单跟踪、支付状态）
- ✅ 评价系统（买卖双方互评）
- ✅ 数据统计和CSV导出（成交率、出价统计）

### 前端功能 (Vue3 + TypeScript)
- ✅ 用户登录注册页面
- ✅ 拍卖品列表（搜索、分页、分类筛选）
- ✅ 拍卖品详情与实时竞价
- ✅ 个人中心（资料、出价、订单、消息）
- ✅ 消息中心
- ✅ 管理后台（用户、拍卖会、订单、统计）
- ✅ 响应式设计

## 技术栈

### 后端
- **框架**: Gin v1.9.1
- **数据库**: MySQL 8.0+ (GORM v1.25.5)
- **缓存**: Redis (go-redis v8)
- **认证**: JWT (golang-jwt v5)
- **邮件**: gomail
- **其他**: godotenv, uuid

### 前端
- **框架**: Vue 3.4
- **语言**: TypeScript 5.3
- **构建工具**: Vite 5.0
- **UI组件**: Element Plus 2.5
- **路由**: Vue Router 4.2
- **状态管理**: Pinia 2.1
- **HTTP客户端**: Axios 1.6
- **日期处理**: Day.js 1.11

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- Redis 6.0+

### 1. 启动后端服务

```bash
cd backend

# 复制配置文件
cp .env.example .env

# 修改 .env 配置（数据库、Redis、邮件等）
vim .env

# 安装依赖
go mod download

# 创建数据库
mysql -u root -p -e "CREATE DATABASE auction_system CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 启动服务
go run cmd/server/main.go
```

后端服务将在 `http://localhost:8080` 启动

### 2. 启动前端服务

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端服务将在 `http://localhost:3000` 启动

## 项目结构

```
.
├── backend/                     # 后端服务
│   ├── cmd/
│   │   └── server/              # 服务入口
│   ├── config/                  # 配置管理
│   ├── internal/
│   │   ├── controllers/         # 控制器层
│   │   ├── services/            # 业务逻辑层
│   │   ├── models/              # 数据模型层
│   │   ├── middleware/          # 中间件
│   │   ├── routers/             # 路由
│   │   ├── dto/                 # 数据传输对象
│   │   └── utils/               # 工具函数
│   ├── pkg/
│   │   ├── logger/              # 日志
│   │   ├── response/            # 响应封装
│   │   ├── redis/               # Redis封装
│   │   └── mail/                # 邮件服务
│   └── uploads/                 # 上传文件目录
└── frontend/                    # 前端应用
    ├── src/
    │   ├── api/                 # API接口
    │   ├── components/          # 组件
    │   ├── router/              # 路由
    │   ├── stores/              # 状态管理
    │   ├── types/               # TypeScript类型
    │   ├── utils/               # 工具函数
    │   ├── views/               # 页面
    │   └── App.vue              # 根组件
    └── package.json
```

## API 文档

### 认证接口
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录

### 拍卖品接口
- `GET /api/v1/items` - 获取拍卖品列表
- `GET /api/v1/items/:id` - 获取拍卖品详情
- `POST /api/v1/items` - 创建拍卖品
- `PUT /api/v1/items/:id` - 更新拍卖品
- `POST /api/v1/items/:id/online` - 上架
- `POST /api/v1/items/:id/offline` - 下架
- `POST /api/v1/items/:id/bid` - 出价
- `GET /api/v1/items/:id/bids` - 获取出价记录

### 拍卖会接口
- `GET /api/v1/sessions` - 获取拍卖会列表
- `GET /api/v1/sessions/active` - 获取进行中的拍卖会
- `POST /api/v1/sessions` - 创建拍卖会（管理员）
- `POST /api/v1/sessions/:id/start` - 开始拍卖会
- `POST /api/v1/sessions/:id/end` - 结束拍卖会

### 订单接口
- `POST /api/v1/orders` - 创建订单
- `GET /api/v1/orders/:id` - 获取订单详情
- `POST /api/v1/orders/:id/pay` - 支付订单
- `POST /api/v1/orders/:id/ship` - 发货
- `POST /api/v1/orders/:id/confirm-delivery` - 确认收货

### 管理接口（需要管理员权限）
- `GET /api/v1/admin/users` - 用户列表
- `PUT /api/v1/admin/users/:id/status` - 更新用户状态
- `GET /api/v1/admin/orders` - 订单列表
- `GET /api/v1/admin/statistics/overall` - 数据统计
- `GET /api/v1/admin/export/orders` - 导出订单CSV
- `GET /api/v1/admin/export/bids` - 导出出价CSV

## 核心功能说明

### 实时竞价流程
1. 用户在拍卖进行时出价
2. 系统通过Redis分布式锁防止并发冲突
3. 验证加价幅度（默认最低加价10元）
4. 验证用户余额
5. 更新当前最高出价和出价历史
6. 通过Redis发布/订阅通知所有关注者
7. 通知前一位出价者其出价被超越

### 自动竞价代理
1. 用户设置最高出价上限
2. 当有其他用户出价时，系统自动加价
3. 每次加价幅度为配置的自动加价幅度
4. 直到达到用户设置的上限或没有其他竞争者

### 消息通知
- 出价被超越通知
- 拍卖即将结束提醒
- 竞拍成功通知
- 支付成功通知
- 支持站内消息和邮件双通道

## 配置说明

### 后端配置 (.env)
```env
# 服务器
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# 数据库
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=123456
DB_NAME=auction_system

# Redis
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRE_HOURS=24

# 邮件
MAIL_HOST=smtp.gmail.com
MAIL_PORT=587
MAIL_USER=your-email@gmail.com
MAIL_PASSWORD=your-app-password

# 拍卖规则
AUCTION_MIN_INCREMENT=10
AUCTION_EXTEND_TIME=300
AUCTION_AUTO_BID_INCREMENT=10
```

## 开发说明

### 数据库表
系统自动创建以下表：
- users - 用户表
- categories - 分类表
- auction_items - 拍卖品表
- auction_images - 拍卖品图片表
- auction_sessions - 拍卖会表
- auction_item_sessions - 拍卖品-拍卖会关联表
- bids - 出价记录表
- auto_bids - 自动出价表
- orders - 订单表
- payments - 支付记录表
- reviews - 评价表
- notifications - 通知表
- system_logs - 系统日志表

### Redis Key 设计
- `bid:lock:{item_id}` - 出价分布式锁
- `bid:current:{item_id}` - 当前最高出价
- `bid:history:{item_id}` - 出价历史（Sorted Set）
- `bid:auto:{item_id}:{user_id}` - 用户自动出价配置
- `bid:channel:{item_id}` - 出价消息频道
- `notification:{user_id}` - 通知消息频道

## License
MIT
