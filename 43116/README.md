# 汽车租赁与预订管理系统

一个完整的在线汽车租赁与预订管理Web应用，专为汽车租赁公司和自驾游平台设计。

## 技术栈

### 后端
- **Go** (Golang 1.21+)
- **Gin** - Web框架
- **GORM** - ORM框架
- **MySQL** - 关系型数据库
- **Redis** - 缓存和消息队列
- **JWT** - 身份认证
- **Excelize** - Excel导出

### 前端
- **Vue 3** - 前端框架
- **TypeScript** - 类型安全
- **Element Plus** - UI组件库
- **Pinia** - 状态管理
- **Vue Router** - 路由
- **Axios** - HTTP客户端
- **ECharts** - 图表库

## 功能特性

### 1. 车辆管理
- 车辆CRUD操作
- 批量图片上传
- 车辆状态管理（可用、出租中、维护中、已停用）
- 车辆信息包括：品牌、型号、座位数、变速箱、日租金、车牌号等

### 2. 预订系统
- 选择取车/还车日期时间
- 选择取车/还车门店
- 自动计算租金总额
- 优惠码折扣支持
- 预订后自动更新车辆状态
- 取车前可取消预订

### 3. 动态定价
- 节假日价格倍率
- 周末价格设置
- 不同车型不同基础价格
- 灵活的定价规则配置

### 4. 门店管理
- 多城市、多门店支持
- 门店营业时间和联系电话
- 按城市筛选门店
- 门店地址和位置信息

### 5. 用户管理
- 用户注册/登录
- JWT Token认证
- 驾驶证信息上传
- 实名认证审核
- 用户状态管理

### 6. 订单管理
- 用户查看历史和当前订单
- 管理员查看所有订单
- 按状态、日期范围、车型筛选
- Excel报表导出

### 7. 评价系统
- 还车后对车辆打分和评价
- 评分影响车辆可信度
- 评价管理和隐藏功能

### 8. 统计仪表盘
- 本月营收统计
- 订单数量统计
- 车辆利用率
- 热门车型排行
- 营收趋势图表

### 9. 消息通知
- 预订成功通知
- 取车提醒
- 还车提醒
- 站内信和邮件通知
- Redis消息队列

### 10. 维护管理
- 车辆维护计划
- 维护中的车辆不可预订
- 维护完成后自动恢复可用状态
- 维护记录管理

## 项目结构

```
.
├── backend/                    # 后端代码
│   ├── cmd/server/             # 入口文件
│   │   └── main.go
│   ├── internal/
│   │   ├── config/             # 配置文件
│   │   │   ├── config.go       # 配置加载
│   │   │   ├── database.go     # 数据库连接
│   │   │   └── redis.go        # Redis连接
│   │   ├── handler/            # API处理器
│   │   ├── middleware/         # 中间件
│   │   ├── model/              # 数据模型
│   │   ├── repository/         # 数据访问层
│   │   ├── router/             # 路由配置
│   │   ├── service/            # 业务逻辑层
│   │   └── utils/              # 工具函数
│   ├── config.yaml             # 配置文件
│   └── go.mod
├── frontend/                   # 前端代码
│   ├── src/
│   │   ├── api/                # API接口
│   │   ├── assets/             # 静态资源
│   │   ├── components/         # 组件
│   │   ├── layouts/            # 布局组件
│   │   ├── router/             # 路由配置
│   │   ├── stores/             # 状态管理
│   │   ├── views/              # 页面组件
│   │   │   └── admin/          # 管理员页面
│   │   └── types/              # 类型定义
│   ├── index.html
│   ├── package.json
│   ├── tsconfig.json
│   └── vite.config.ts
└── nginx/                      # Nginx配置
```

## 快速开始

### 前置要求
- Go 1.21+
- Node.js 18+
- MySQL 5.7+
- Redis 6.0+

### 后端启动

1. 进入后端目录
```bash
cd backend
```

2. 安装依赖
```bash
go mod download
```

3. 修改配置文件 `config.yaml` 中的数据库和Redis连接信息

4. 运行服务
```bash
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动

### 前端启动

1. 进入前端目录
```bash
cd frontend
```

2. 安装依赖
```bash
npm install
```

3. 启动开发服务器
```bash
npm run dev
```

访问 `http://localhost:3000`

### 默认账户

- 管理员：`admin` / `admin123`
- 用户：自行注册

## API 接口

### 认证相关
- `POST /api/register` - 用户注册
- `POST /api/login` - 用户登录
- `POST /api/refresh-token` - 刷新Token

### 车辆相关
- `GET /api/cars` - 获取车辆列表
- `GET /api/cars/:id` - 获取车辆详情
- `POST /api/cars` - 创建车辆（管理员）
- `PUT /api/cars/:id` - 更新车辆（管理员）
- `DELETE /api/cars/:id` - 删除车辆（管理员）

### 预订相关
- `POST /api/bookings` - 创建预订
- `GET /api/bookings` - 获取预订列表
- `PUT /api/bookings/:id/cancel` - 取消预订

### 订单相关
- `GET /api/orders` - 获取订单列表
- `GET /api/admin/orders/export` - 导出订单Excel

### 其他接口
详见代码中的路由配置。

## 数据库表结构

- `users` - 用户表
- `roles` - 角色表
- `permissions` - 权限表
- `cars` - 车辆表
- `car_images` - 车辆图片表
- `cities` - 城市表
- `stores` - 门店表
- `bookings` - 预订表
- `orders` - 订单表
- `promo_codes` - 优惠码表
- `pricing_rules` - 定价规则表
- `reviews` - 评价表
- `maintenance_plans` - 维护计划表
- `messages` - 消息表
- `notifications` - 通知表
- `operation_logs` - 操作日志表

## 功能截图

### 用户端
- 首页：车辆搜索、热门车辆展示
- 车辆列表：按条件筛选车辆
- 车辆详情：车辆信息、预订功能
- 我的预订：预订记录管理
- 我的订单：订单记录查看
- 个人中心：个人信息、实名认证
- 消息中心：系统消息查看

### 管理端
- 仪表盘：营收、订单、利用率统计
- 用户管理：用户列表、认证审核
- 车辆管理：车辆CRUD、图片上传
- 门店管理：门店信息管理
- 城市管理：城市信息管理
- 预订管理：预订确认、取消
- 订单管理：订单状态、导出
- 评价管理：评价审核、隐藏
- 维护管理：维护计划管理
- 优惠码管理：优惠码CRUD
- 定价规则：定价规则配置

## License

MIT
