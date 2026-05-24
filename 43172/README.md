# 奢侈品二手交易与鉴定平台

基于 Go (Gin) + Vue3 + TypeScript + Redis 构建的在线二手奢侈品交易与鉴定服务平台。

## 项目特性

### 核心功能

1. **用户管理**
   - 支持买家、卖家、鉴定师三种角色注册认证
   - 鉴定师需要提交资质审核
   - 不同角色权限严格隔离

2. **商品管理**
   - 卖家可上传奢侈品商品
   - 支持多品牌分类（包包、首饰、手表、服装、鞋履等）
   - 每件商品支持多角度高清图片上传和详细描述

3. **鉴定服务**
   - 买家下单前可申请第三方鉴定
   - 鉴定师接单后出具电子鉴定报告
   - 鉴定结果关联商品展示

4. **交易流程**
   - 完整的下单→支付→发货→验收→确认链条
   - Redis缓存的库存管理
   - 交易状态实时追踪

5. **真伪保障系统**
   - 鉴定通过的商品显示真伪认证标签
   - 鉴定报告永久存档可追溯
   - PDF鉴定报告导出

6. **评价与信用体系**
   - 交易完成后双向评价
   - 鉴定师有独立评分
   - 用户信用影响交易优先级

7. **数据统计仪表盘**
   - 交易趋势图表
   - 热门品牌排行
   - 鉴定通过率统计
   - 实时数据可视化

### 技术特性

- 边界条件校验
- 异常捕获重试
- 操作日志审计
- 配置热加载
- 数据库事务保证
- 分页搜索排序
- 文件上传下载
- Redis缓存管理
- JWT认证授权

## 技术栈

### 后端
- **语言**: Go 1.21
- **框架**: Gin
- **数据库**: MySQL 8.0
- **缓存**: Redis 7
- **ORM**: GORM
- **认证**: JWT
- **日志**: Logrus
- **PDF**: gopdf
- **配置**: Viper

### 前端
- **框架**: Vue 3.4
- **语言**: TypeScript
- **构建工具**: Vite 5
- **UI组件**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **HTTP客户端**: Axios
- **图表**: ECharts
- **日期处理**: Day.js

## 项目结构

```
luxury-trading-platform/
├── backend/                    # 后端项目
│   ├── cmd/server/            # 入口
│   │   └── main.go
│   ├── configs/               # 配置文件
│   │   └── config.yaml
│   ├── internal/
│   │   ├── cache/             # Redis缓存层
│   │   ├── config/            # 配置管理
│   │   ├── handler/           # HTTP处理器
│   │   ├── middleware/        # 中间件
│   │   ├── model/             # 数据模型
│   │   ├── repository/        # 数据访问层
│   │   ├── router/            # 路由注册
│   │   ├── service/           # 业务逻辑层
│   │   └── utils/             # 工具函数
│   ├── uploads/               # 上传文件目录
│   ├── Dockerfile
│   └── go.mod
├── frontend/                   # 前端项目
│   ├── src/
│   │   ├── api/               # API接口
│   │   ├── layouts/           # 布局组件
│   │   ├── router/            # 路由配置
│   │   ├── stores/            # Pinia状态管理
│   │   ├── styles/            # 全局样式
│   │   ├── types/             # TypeScript类型
│   │   ├── views/             # 页面组件
│   │   ├── App.vue
│   │   └── main.ts
│   ├── Dockerfile
│   ├── nginx.conf
│   ├── index.html
│   ├── package.json
│   └── vite.config.ts
├── docker-compose.yml          # Docker编排
└── .gitignore
```

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- Redis 7+

### 本地开发

#### 1. 启动数据库

```bash
docker-compose up -d mysql redis
```

#### 2. 启动后端

```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

后端服务运行在 http://localhost:8080

#### 3. 启动前端

```bash
cd frontend
npm install
npm run dev
```

前端服务运行在 http://localhost:3000

### Docker部署

```bash
docker-compose up -d
```

访问地址：
- 前端: http://localhost:3000
- 后端API: http://localhost:8080

## API文档

### 认证相关

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/auth/register | 用户注册 |
| POST | /api/v1/auth/register/authenticator | 鉴定师注册 |
| POST | /api/v1/auth/login | 用户登录 |

### 商品相关

| 方法 | 路径 | 描述 | 权限 |
|------|------|------|------|
| GET | /api/v1/products | 商品列表 | 公开 |
| GET | /api/v1/products/:id | 商品详情 | 公开 |
| POST | /api/v1/products | 创建商品 | 卖家 |
| PUT | /api/v1/products/:id | 更新商品 | 卖家 |
| DELETE | /api/v1/products/:id | 删除商品 | 卖家 |
| POST | /api/v1/products/:id/images | 上传图片 | 卖家 |

### 订单相关

| 方法 | 路径 | 描述 | 权限 |
|------|------|------|------|
| POST | /api/v1/orders | 创建订单 | 买家 |
| GET | /api/v1/orders | 订单列表 | 登录用户 |
| POST | /api/v1/orders/:id/pay | 支付订单 | 买家 |
| POST | /api/v1/orders/:id/ship | 发货 | 卖家 |
| POST | /api/v1/orders/:id/confirm | 确认收货 | 买家 |
| POST | /api/v1/orders/:id/cancel | 取消订单 | 订单相关方 |

### 鉴定相关

| 方法 | 路径 | 描述 | 权限 |
|------|------|------|------|
| POST | /api/v1/authentications | 申请鉴定 | 买家 |
| GET | /api/v1/authentications | 鉴定列表 | 登录用户 |
| POST | /api/v1/authentications/:id/accept | 接单 | 鉴定师 |
| POST | /api/v1/authentications/:id/complete | 完成鉴定 | 鉴定师 |
| GET | /api/v1/authentications/:id/report/download | 下载报告 | 登录用户 |

### 评价相关

| 方法 | 路径 | 描述 | 权限 |
|------|------|------|------|
| POST | /api/v1/reviews | 创建评价 | 登录用户 |
| GET | /api/v1/reviews | 评价列表 | 公开 |
| GET | /api/v1/reviews/user/:id/rating | 用户评分 | 公开 |

### 管理相关

| 方法 | 路径 | 描述 | 权限 |
|------|------|------|------|
| GET | /api/v1/admin/users | 用户列表 | 管理员 |
| GET | /api/v1/admin/authenticators | 鉴定师列表 | 管理员 |
| POST | /api/v1/admin/authenticators/:id/approve | 通过鉴定师 | 管理员 |
| POST | /api/v1/admin/authenticators/:id/reject | 拒绝鉴定师 | 管理员 |
| GET | /api/v1/statistics/dashboard | 数据统计 | 管理员 |

## 数据模型

### 用户角色

- **buyer** (买家): 浏览商品、下单、评价
- **seller** (卖家): 发布商品、管理订单、发货
- **authenticator** (鉴定师): 接鉴定任务、出具报告
- **admin** (管理员): 用户管理、鉴定师审核、数据统计

### 核心实体

- User: 用户
- Product: 商品
- ProductImage: 商品图片
- Brand: 品牌
- Order: 订单
- Authentication: 鉴定记录
- Review: 评价
- AuditLog: 审计日志
- AuthenticatorProfile: 鉴定师资质
- SellerProfile: 卖家信息
- BuyerProfile: 买家信息

## 配置说明

### 后端配置 (config.yaml)

```yaml
server:
  port: 8080
  mode: "debug"  # debug/release

database:
  host: "localhost"
  port: 3306
  user: "root"
  password: "root"
  dbname: "luxury_trading"
  charset: "utf8mb4"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "your-secret-key"
  expire_time: 24h

log:
  level: "debug"  # debug/info/warn/error
  file: "logs/app.log"

upload:
  path: "./uploads"
  max_size: 10485760  # 10MB
  types: [".jpg", ".jpeg", ".png", ".gif", ".pdf"]
```

### 环境变量

也支持通过环境变量覆盖配置：

- `DB_HOST` - 数据库主机
- `DB_PORT` - 数据库端口
- `DB_USER` - 数据库用户
- `DB_PASSWORD` - 数据库密码
- `DB_NAME` - 数据库名
- `REDIS_HOST` - Redis主机
- `REDIS_PORT` - Redis端口
- `JWT_SECRET` - JWT密钥

## 安全特性

1. **认证授权**: JWT Token 认证，角色权限控制
2. **密码安全**: bcrypt 加密存储
3. **SQL注入防护**: GORM 参数化查询
4. **XSS防护**: 前端自动转义
5. **CORS配置**: 跨域安全控制
6. **审计日志**: 所有操作记录可追溯

## 缓存策略

- 用户信息: 24小时
- 商品详情: 1小时
- 库存数据: 实时更新
- 统计数据: 5分钟

## License

MIT
