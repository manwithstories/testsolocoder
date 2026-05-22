# 在线活动报名与票务管理系统

基于 Go + Gin + Vue3 + TypeScript 开发的活动票务管理系统。

## 功能特性

### 1. 用户管理
- 注册/登录（JWT认证）
- 管理员/普通用户角色权限
- 用户信息管理

### 2. 活动管理
- 活动创建、编辑、删除
- 活动发布/取消
- 海报图片上传
- 活动详情展示

### 3. 票务管理
- 多种票型（普通票、VIP票、早鸟票）
- 独立库存和价格
- Redis原子扣减防止超卖
- 售罄自动下架

### 4. 报名购票
- 浏览活动选择票型
- 批量购票
- 优惠券抵扣
- 事务保证数据一致性

### 5. 签到管理
- 唯一签到二维码
- 扫码签到
- 签到时间记录

### 6. 统计报表
- 按活动统计
- 按票型统计
- 按时间段统计
- 分页搜索排序
- Excel导出

## 技术栈

### 后端
- Go 1.21
- Gin 框架
- GORM ORM
- MySQL 数据库
- Redis 缓存（库存扣减）
- JWT 认证
- Logrus 日志
- Excelize Excel导出

### 前端
- Vue 3 + TypeScript
- Element Plus UI
- Pinia 状态管理
- Vue Router 路由
- Axios HTTP客户端
- ECharts 图表
- Day.js 日期处理

## 项目结构

```
ticket-system/
├── backend/                    # 后端项目
│   ├── config/                 # 配置管理
│   ├── internal/
│   │   ├── common/             # 公共模块
│   │   │   ├── exception/      # 异常处理
│   │   │   └── response/       # 响应封装
│   │   ├── controller/         # 控制器层
│   │   ├── dto/                # 数据传输对象
│   │   ├── jwt/                # JWT认证
│   │   ├── logger/             # 日志
│   │   ├── middleware/         # 中间件
│   │   ├── models/             # 数据模型
│   │   ├── redis/              # Redis客户端
│   │   ├── router/             # 路由注册
│   │   ├── seeder/             # 数据种子
│   │   ├── service/            # 业务逻辑层
│   │   └── util/               # 工具函数
│   ├── main.go                 # 入口文件
│   ├── go.mod
│   ├── .env.development
│   └── .env.example
└── frontend/                   # 前端项目
    ├── src/
    │   ├── api/                # API接口
    │   ├── layout/             # 布局组件
    │   ├── router/             # 路由配置
    │   ├── store/              # 状态管理
    │   ├── styles/             # 样式文件
    │   ├── views/              # 页面组件
    │   ├── App.vue
    │   └── main.ts
    ├── index.html
    ├── package.json
    ├── tsconfig.json
    └── vite.config.ts
```

## 快速开始

### 环境要求
- Go >= 1.21
- Node.js >= 16
- MySQL >= 5.7
- Redis >= 6.0

### 数据库准备

创建数据库：
```sql
CREATE DATABASE ticket_system DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 后端启动

```bash
cd backend

# 安装依赖
go mod download

# 复制配置文件
cp .env.example .env.development

# 修改配置文件中的数据库和Redis连接信息

# 启动服务
go run main.go
```

后端服务启动在 http://localhost:8080

默认管理员账户：
- 用户名：admin
- 密码：admin123

### 前端启动

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端服务启动在 http://localhost:5173

### 生产构建

```bash
# 后端构建
cd backend
go build -o ticket-server main.go

# 前端构建
cd frontend
npm run build
```

## API 接口

### 认证接口
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录

### 活动接口
- `GET /api/activities` - 获取活动列表
- `GET /api/activities/:id` - 获取活动详情
- `POST /api/activities` - 创建活动（管理员）
- `PUT /api/activities/:id` - 更新活动（管理员）
- `PUT /api/activities/:id/status` - 更新活动状态（管理员）
- `DELETE /api/activities/:id` - 删除活动（管理员）

### 票型接口
- `GET /api/ticket-types` - 获取票型列表
- `POST /api/ticket-types` - 创建票型（管理员）
- `PUT /api/ticket-types/:id` - 更新票型（管理员）
- `DELETE /api/ticket-types/:id` - 删除票型（管理员）

### 订单接口
- `GET /api/orders` - 获取订单列表
- `GET /api/orders/:id` - 获取订单详情
- `POST /api/orders` - 创建订单
- `POST /api/orders/:id/pay` - 支付订单
- `POST /api/orders/:id/cancel` - 取消订单

### 签到接口
- `POST /api/checkins` - 签到
- `GET /api/checkins` - 获取签到列表
- `GET /api/checkins/statistics` - 获取签到统计

### 统计接口（管理员）
- `GET /api/statistics/activities` - 按活动统计
- `GET /api/statistics/ticket-types` - 按票型统计
- `GET /api/statistics/daily` - 每日统计
- `GET /api/statistics/export` - 导出Excel

## 核心设计

### 库存扣减流程
1. 使用 Redis `DECRBY` 原子操作扣减库存
2. 扣减失败直接返回库存不足
3. 扣减成功后创建订单
4. 订单创建失败或取消时回滚库存

### 订单事务流程
1. Redis 分布式锁防止重复提交
2. Redis 原子扣减库存
3. 数据库事务创建订单、订单项、签到记录
4. 更新优惠券使用状态
5. 更新票型已售数量
6. 提交事务或回滚

### 权限控制
- 公开接口：活动列表、活动详情
- 用户接口：登录后可访问（个人订单、购票）
- 管理员接口：需要管理员角色（活动管理、票务管理、统计报表等）

## License

MIT
