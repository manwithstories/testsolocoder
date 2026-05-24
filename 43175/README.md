# 智能家居能源管理平台

基于 Go (Gin) + Vue3 + TypeScript + Redis 构建的智能家居设备管理与家庭能源消耗监控平台。

## 技术栈

### 后端
- **Go 1.21** - 编程语言
- **Gin** - Web 框架
- **GORM** - ORM 框架
- **SQLite** - 数据库（可切换为 MySQL/PostgreSQL）
- **Redis** - 缓存和实时数据同步
- **JWT** - 身份认证
- **robfig/cron** - 定时任务调度
- **tealeg/xlsx** - Excel 报表导出

### 前端
- **Vue 3** - 前端框架
- **TypeScript** - 类型安全
- **Vite** - 构建工具
- **Pinia** - 状态管理
- **Vue Router** - 路由管理
- **Element Plus** - UI 组件库
- **ECharts** - 图表可视化
- **Axios** - HTTP 客户端
- **Day.js** - 日期处理

## 项目结构

```
.
├── backend/                    # Go 后端
│   ├── main.go                # 入口文件
│   ├── config/                # 配置管理
│   ├── models/                # 数据模型
│   ├── middleware/            # 中间件（JWT、CORS、日志）
│   ├── handlers/              # 请求处理器
│   ├── services/              # 业务服务层
│   ├── utils/                 # 工具函数
│   └── .env                   # 环境变量配置
├── frontend/                   # Vue3 前端
│   ├── src/
│   │   ├── api/               # API 接口封装
│   │   ├── layouts/           # 布局组件
│   │   ├── views/             # 页面组件
│   │   ├── stores/            # Pinia 状态管理
│   │   ├── router/            # 路由配置
│   │   └── styles/            # 全局样式
│   ├── index.html
│   ├── vite.config.ts
│   └── package.json
└── README.md
```

## 核心功能

### 1. 用户认证与家庭管理
- 用户注册/登录（JWT Token 认证）
- 家庭创建与管理
- 成员邀请（邮箱邀请）
- 角色权限（管理员、普通用户、访客）

### 2. 设备管理
- 设备 CRUD（添加、编辑、删除）
- 设备参数（类型、厂商、位置、功率、通信协议）
- 设备状态实时同步到 Redis
- 支持设备类型：照明、空调、取暖、窗帘、摄像头、传感器、开关、插座等

### 3. 设备分组
- 按房间、楼层、功能分组
- 批量控制（一键开启/关闭）
- 分组能耗统计

### 4. 能耗监控
- 实时用电功率监控
- 按设备/房间/时间段统计能耗
- 日/周/月能耗趋势分析
- 异常高耗电自动预警

### 5. 场景联动
- 创建智能场景（离家模式、睡眠模式等）
- 支持时间、设备状态、传感器多维度触发条件
- 一键执行场景动作

### 6. 定时任务
- 设置定时开关、调温等自动化任务
- Cron 表达式灵活配置
- 任务执行日志记录
- 能耗影响追踪

### 7. 报表导出
- Excel 格式能耗报表导出
- 包含能耗明细、趋势数据
- 节能建议生成

### 8. 通知提醒
- 设备离线告警
- 能耗异常预警
- 定时任务执行失败通知
- 站内消息 + 邮件推送

## 快速开始

### 后端启动

```bash
cd backend

# 安装依赖
go mod tidy

# 配置 .env 文件（根据需要修改）
cp .env.example .env

# 启动服务
go run main.go
```

后端服务默认运行在 `http://localhost:8080`

### 前端启动

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端服务默认运行在 `http://localhost:3000`

### Redis 配置

确保 Redis 服务已启动，默认配置：
- 地址：127.0.0.1:6379
- 密码：空
- 数据库：0

修改 `backend/.env` 可自定义 Redis 配置。

## API 接口

所有接口前缀：`/api/v1`

### 认证接口
- POST `/auth/register` - 用户注册
- POST `/auth/login` - 用户登录
- POST `/auth/refresh` - 刷新 Token

### 用户接口
- GET `/user/profile` - 获取个人信息
- PUT `/user/profile` - 更新个人信息
- PUT `/user/password` - 修改密码

### 家庭接口
- GET/POST `/families` - 家庭列表/创建
- GET/PUT/DELETE `/families/:id` - 家庭详情/更新/删除
- POST `/families/:id/invite` - 邀请成员
- DELETE `/families/:id/members/:memberId` - 移除成员

### 设备接口
- GET/POST `/devices` - 设备列表/创建
- GET/PUT/DELETE `/devices/:id` - 设备详情/更新/删除
- PUT `/devices/:id/status` - 更新设备状态
- GET `/devices/:id/energy` - 获取设备能耗

### 分组接口
- GET/POST `/groups` - 分组列表/创建
- GET/PUT/DELETE `/groups/:id` - 分组详情/更新/删除
- POST `/groups/:id/devices` - 添加设备到分组
- PUT `/groups/:id/control` - 批量控制分组

### 能耗接口
- GET `/energy/realtime` - 实时能耗
- GET `/energy/statistics` - 能耗统计
- GET `/energy/trend` - 能耗趋势
- GET `/energy/alerts` - 能耗告警
- GET `/energy/export` - 导出报表

### 场景接口
- GET/POST `/scenes` - 场景列表/创建
- GET/PUT/DELETE `/scenes/:id` - 场景详情/更新/删除
- POST `/scenes/:id/execute` - 执行场景

### 定时任务接口
- GET/POST `/schedules` - 任务列表/创建
- GET/PUT/DELETE `/schedules/:id` - 任务详情/更新/删除
- GET `/schedules/:id/logs` - 任务执行日志

### 通知接口
- GET `/notifications` - 通知列表
- PUT `/notifications/:id/read` - 标记已读
- PUT `/notifications/read-all` - 全部已读
- DELETE `/notifications/:id` - 删除通知

## 数据校验

平台实现了完善的数据校验：

- **设备参数校验**：功率范围 0-10000W，必填字段验证
- **能耗数据验证**：数值范围检查，异常值过滤
- **定时任务校验**：Cron 表达式格式验证，时间冲突检测
- **用户输入验证**：邮箱格式、密码强度、字段长度等
- **权限校验**：基于角色的访问控制（RBAC）

## 异常处理

- 全局异常捕获，详细日志记录
- 友好的错误提示信息
- Redis 连接失败自动降级
- 数据库操作失败回滚

## 开发说明

### 添加新设备类型

在 `frontend/src/api/device.ts` 的 `deviceTypeOptions` 数组中添加新类型。

### 添加新的通知类型

在 `frontend/src/api/notification.ts` 的 `notificationTypeMap` 中添加新类型配置。

### 自定义场景图标

在 `frontend/src/api/scene.ts` 的 `sceneIcons` 数组中添加新图标。

## License

MIT
