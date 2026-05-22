# 酒店管理系统 (Hotel Management System)

基于 Go + Gin + Vue3 + TypeScript + PostgreSQL 的全栈酒店管理系统，专为中小型酒店和民宿设计。

## 功能特性

### 🏨 客房管理
- 房型管理：添加、编辑、删除房型，设置基础价格、床位数、容纳人数、配套设施
- 房间管理：按楼层管理房间，设置房号、价格、状态
- 批量导入：支持Excel批量导入房间数据
- 房间状态：实时维护空闲、已入住、已预订、维修中状态

### 📅 预订管理
- 在线预订：选择房型、房间、入住/离店日期，自动计算价格
- 预订确认：前台确认预订，生成预订号
- 预订修改：修改入住日期、换房等
- 预订取消：支持免费取消期校验，超时取消收取违约金
- 并发控制：使用数据库事务和悲观锁防止重复预订

### 🛏️ 入住管理
- 入住登记：记录客人身份信息，支持从预订直接入住
- 续住服务：延长入住时间，自动计算续住费用
- 退房结算：自动计算房费、额外费用，多方式收款
- 房间状态：自动更新房间状态

### 💰 支付记录
- 多支付方式：支持现金、微信、支付宝、银行卡、转账
- 支付类型：预付、补付、退款、押金
- 支付流水：完整的支付记录和凭证
- 订单关联：每笔支付关联对应预订或入住订单

### 🎖️ 会员体系
- 会员注册：手机号注册，自动生成会员号
- 会员等级：普通会员、银卡、金卡、钻石会员
- 会员折扣：不同等级享受不同折扣率
- 积分系统：消费累计积分，积分可抵扣现金（100积分=1元）
- 自动升级：积分达标自动升级会员等级

### 📊 统计报表
- 入住率趋势：按日/周/月统计入住率
- 营收明细：房费收入、其他收入、总营收
- 日期筛选：支持按日期范围查询
- Excel导出：报表数据一键导出Excel

### 🏠 房间状态看板
- 实时展示：所有房间状态一目了然
- 颜色区分：空闲(绿)、已入住(蓝)、已预订(黄)、维修中(红)
- 楼层筛选：按楼层查看房间状态
- 快速操作：点击房间可快速办理入住/退房

### 🔐 权限管理
- JWT认证：安全的Token认证机制
- RBAC权限：管理员、前台、普通用户三级角色
- 菜单权限：根据角色动态显示菜单
- 操作日志：记录关键操作

## 技术栈

### 后端
- **语言**: Go 1.21
- **框架**: Gin v1.9.1
- **ORM**: GORM v1.25.5
- **数据库**: PostgreSQL 13+
- **认证**: JWT (golang-jwt/v5)
- **密码加密**: bcrypt
- **日志**: Logrus
- **Excel处理**: Excelize v2
- **参数校验**: validator v10

### 前端
- **框架**: Vue 3.4
- **语言**: TypeScript 5.3
- **构建工具**: Vite 5.0
- **状态管理**: Pinia 2.1
- **路由**: Vue Router 4.2
- **UI组件**: Element Plus 2.4
- **图表**: ECharts 5.4
- **HTTP客户端**: Axios 1.6
- **日期处理**: Day.js 1.11

## 项目结构

```
hotel-management-system/
├── backend/                     # 后端代码
│   ├── cmd/
│   │   └── server/
│   │       └── main.go         # 应用入口
│   ├── internal/
│   │   ├── config/             # 配置管理
│   │   ├── database/           # 数据库连接和初始化
│   │   ├── middleware/         # 中间件（认证、日志、CORS等）
│   │   ├── model/              # 数据模型
│   │   ├── repository/         # 数据访问层
│   │   ├── service/            # 业务逻辑层
│   │   ├── handler/            # HTTP处理器
│   │   ├── dto/                # 数据传输对象
│   │   ├── utils/              # 工具函数
│   │   └── pkg/                # 内部包（JWT、RBAC、日志）
│   ├── pkg/
│   │   └── excel/              # Excel处理工具
│   ├── scripts/
│   │   └── init.sql            # 数据库初始化脚本
│   ├── .env                    # 环境变量配置
│   ├── go.mod
│   └── go.sum
└── frontend/                   # 前端代码
    ├── src/
    │   ├── api/                # API接口
    │   ├── components/         # 公共组件
    │   ├── layout/             # 布局组件
    │   ├── router/             # 路由配置
    │   ├── store/              # Pinia状态管理
    │   ├── types/              # TypeScript类型定义
    │   ├── utils/              # 工具函数
    │   ├── views/              # 页面组件
    │   ├── styles/             # 全局样式
    │   ├── App.vue
    │   └── main.ts
    ├── public/
    ├── index.html
    ├── package.json
    ├── tsconfig.json
    └── vite.config.ts
```

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- PostgreSQL 13+

### 1. 数据库准备

```sql
-- 创建数据库
CREATE DATABASE hotel_db;

-- 执行初始化脚本
psql -d hotel_db -f backend/scripts/init.sql
```

### 2. 后端启动

```bash
cd backend

# 安装依赖
go mod download

# 配置环境变量（修改 .env 文件）
# DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME

# 启动服务
go run ./cmd/server/
```

后端服务将在 `http://localhost:8080` 启动

### 3. 前端启动

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端服务将在 `http://localhost:5173` 启动

### 默认账户

| 角色 | 用户名 | 密码 |
|------|--------|------|
| 管理员 | admin | admin123 |
| 前台 | frontdesk | frontdesk123 |

## API文档

### 认证接口
- `POST /api/auth/login` - 登录
- `POST /api/auth/register` - 注册
- `GET /api/auth/profile` - 获取当前用户信息

### 用户管理（管理员）
- `GET /api/users` - 用户列表
- `POST /api/users` - 创建用户
- `PUT /api/users/:id` - 更新用户
- `DELETE /api/users/:id` - 删除用户

### 房型管理
- `GET /api/room-types` - 房型列表
- `POST /api/room-types` - 创建房型
- `PUT /api/room-types/:id` - 更新房型
- `DELETE /api/room-types/:id` - 删除房型

### 房间管理
- `GET /api/rooms` - 房间列表
- `GET /api/rooms/available` - 可预订房间
- `POST /api/rooms` - 创建房间
- `PUT /api/rooms/:id` - 更新房间
- `DELETE /api/rooms/:id` - 删除房间
- `POST /api/rooms/batch-import` - 批量导入房间
- `GET /api/rooms/status-board` - 房间状态看板

### 预订管理
- `GET /api/bookings` - 预订列表
- `GET /api/bookings/:id` - 预订详情
- `POST /api/bookings` - 创建预订
- `PUT /api/bookings/:id` - 修改预订
- `POST /api/bookings/:id/confirm` - 确认预订
- `POST /api/bookings/:id/cancel` - 取消预订
- `POST /api/bookings/calculate-price` - 计算预订价格

### 入住管理
- `GET /api/checkins` - 入住列表
- `GET /api/checkins/:id` - 入住详情
- `POST /api/checkins` - 办理入住
- `POST /api/checkins/:id/checkout` - 退房结算
- `POST /api/checkins/:id/extend` - 续住
- `POST /api/checkins/:id/extra-charge` - 添加额外消费

### 支付管理
- `GET /api/payments` - 支付记录列表
- `GET /api/payments/:id` - 支付详情
- `POST /api/payments` - 创建支付
- `POST /api/payments/refund` - 退款
- `GET /api/payments/order` - 订单支付流水
- `GET /api/payments/:id/voucher` - 支付凭证

### 会员管理
- `GET /api/members` - 会员列表
- `GET /api/members/:id` - 会员详情
- `POST /api/members` - 注册会员
- `PUT /api/members/:id` - 更新会员
- `DELETE /api/members/:id` - 删除会员
- `POST /api/members/:id/points/use` - 使用积分
- `POST /api/members/:id/points/recharge` - 充值积分
- `GET /api/members/:id/discount` - 获取会员折扣

### 会员等级管理
- `GET /api/member-levels` - 等级列表
- `POST /api/member-levels` - 创建等级
- `PUT /api/member-levels/:id` - 更新等级
- `DELETE /api/member-levels/:id` - 删除等级

### 统计报表
- `GET /api/reports/occupancy-rate` - 入住率统计
- `GET /api/reports/revenue` - 营收统计
- `GET /api/reports/export` - 导出Excel报表

### 看板
- `GET /api/dashboard/stats` - 看板统计数据
- `GET /api/dashboard/board` - 房间状态看板
- `GET /api/dashboard/floor/:floor` - 楼层房间列表

## 核心业务流程

### 预订流程
1. 用户选择房型和入住/离店日期
2. 系统查询可用房间并展示价格
3. 填写客人信息，提交预订
4. 系统使用事务锁检查房间可用性
5. 计算总价（会员自动应用折扣）
6. 生成预订号，设置免费取消期限
7. 前台确认预订

### 入住流程
1. 客人到店，前台查询预订或直接办理
2. 登记客人身份信息
3. 收取押金
4. 更新房间状态为已入住
5. 生成入住记录

### 退房流程
1. 客人退房，查询入住记录
2. 计算总费用（房费+额外消费-押金）
3. 选择支付方式，完成支付
4. 生成支付记录
5. 会员消费累计积分
6. 更新房间状态为空闲

### 会员积分流程
1. 消费时自动计算积分（金额 × 等级倍率）
2. 积分累计到会员账户
3. 检查是否达到升级条件，自动升级
4. 下次消费可使用积分抵扣（100积分=1元）

## 部署建议

### 后端部署
```bash
# 编译
go build -o hotel-server ./cmd/server/

# 后台运行
nohup ./hotel-server > server.log 2>&1 &
```

### 前端部署
```bash
# 构建
npm run build

# 将 dist 目录部署到 Nginx 或其他静态文件服务器
```

### Nginx 配置示例
```nginx
server {
    listen 80;
    server_name hotel.example.com;

    # 前端静态文件
    location / {
        root /var/www/hotel/dist;
        try_files $uri $uri/ /index.html;
    }

    # API代理
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 许可证

MIT License
