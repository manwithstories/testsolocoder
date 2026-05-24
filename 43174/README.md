# 校园二手教材交易与学习笔记共享平台

一个基于Go Gin + React + TypeScript + PostgreSQL的全栈校园二手教材交易平台。

## 技术栈

### 后端
- **Go 1.21** - 编程语言
- **Gin** - Web框架
- **GORM** - ORM库
- **PostgreSQL** - 数据库
- **JWT** - 身份认证
- **Excelize** - Excel报表导出

### 前端
- **React 18** - UI框架
- **TypeScript** - 类型安全
- **Ant Design** - UI组件库
- **Tailwind CSS** - 样式框架
- **Zustand** - 状态管理
- **Axios** - HTTP客户端

## 功能特性

### 1. 教材管理
- 教材录入、编辑、删除
- ISBN扫描快速检索
- 按课程分类筛选
- 封面图片上传
- 浏览量统计

### 2. 用户管理
- 学生/书商/管理员三种角色
- JWT身份认证
- 学生学校验证
- 书商资质审核
- 用户评分系统

### 3. 交易管理
- 发布教材出售/交换
- 协商定价
- 交易状态追踪
- 自动更新库存

### 4. 笔记共享
- 上传扫描/PDF笔记
- 按科目分类整理
- 笔记评价系统
- 精选笔记推荐

### 5. 订单管理
- 订单生成
- 支付模拟
- 状态流转记录
- 物流追踪
- 历史订单查询

### 6. 站内消息
- 买卖双方私信沟通
- 议价功能
- 纠纷申诉
- 消息推送通知

### 7. 评价系统
- 教材/笔记评分
- 高评分用户推荐
- 恶意评价管理

### 8. 数据统计
- 教材流通量统计
- 用户交易量排行
- 热门教材分析
- 月度报表导出(Excel)

## 项目结构

```
.
├── backend/                 # 后端服务
│   ├── cmd/server/         # 入口文件
│   ├── config/             # 配置管理
│   ├── internal/
│   │   ├── handlers/       # HTTP处理器
│   │   ├── middleware/     # 中间件
│   │   ├── models/         # 数据模型
│   │   ├── repository/     # 数据访问层
│   │   ├── services/       # 业务逻辑层
│   │   └── utils/          # 工具函数
│   ├── pkg/
│   │   ├── database/       # 数据库连接
│   │   ├── jwt/            # JWT认证
│   │   └── upload/         # 文件上传
│   ├── migrations/         # 数据库迁移
│   └── Dockerfile
├── frontend/               # 前端应用
│   ├── src/
│   │   ├── components/     # 组件
│   │   ├── pages/          # 页面
│   │   ├── services/       # API服务
│   │   ├── context/        # 状态管理
│   │   └── types/          # 类型定义
│   └── Dockerfile
└── docker-compose.yml      # Docker编排配置
```

## 快速开始

### 方式一：Docker Compose (推荐)

```bash
# 克隆项目
git clone <repository-url>
cd 43174

# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

访问:
- 前端: http://localhost:3000
- 后端API: http://localhost:8080/api/v1

### 方式二：本地开发

#### 后端

```bash
cd backend

# 安装依赖
go mod tidy

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件配置数据库连接等

# 运行服务
go run cmd/server/main.go
```

#### 前端

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

## 默认账户

| 角色 | 用户名 | 密码 |
|------|--------|------|
| 管理员 | admin | admin123 |

## API文档

### 认证
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录

### 用户
- `GET /api/v1/users/profile` - 获取个人资料
- `PUT /api/v1/users/profile` - 更新个人资料
- `PUT /api/v1/users/password` - 修改密码
- `POST /api/v1/users/avatar` - 上传头像

### 教材
- `GET /api/v1/textbooks` - 获取教材列表
- `GET /api/v1/textbooks/:id` - 获取教材详情
- `POST /api/v1/textbooks` - 发布教材
- `PUT /api/v1/textbooks/:id` - 更新教材
- `DELETE /api/v1/textbooks/:id` - 删除教材
- `GET /api/v1/textbooks/search/isbn?isbn=xxx` - ISBN搜索

### 笔记
- `GET /api/v1/notes` - 获取笔记列表
- `GET /api/v1/notes/:id` - 获取笔记详情
- `POST /api/v1/notes` - 上传笔记
- `GET /api/v1/notes/featured` - 获取精选笔记

### 订单
- `GET /api/v1/orders` - 获取订单列表
- `POST /api/v1/orders` - 创建订单
- `PUT /api/v1/orders/:id/pay` - 支付订单
- `PUT /api/v1/orders/:id/ship` - 发货
- `PUT /api/v1/orders/:id/deliver` - 确认收货
- `PUT /api/v1/orders/:id/complete` - 完成订单

### 消息
- `GET /api/v1/messages/conversation` - 获取对话记录
- `POST /api/v1/messages` - 发送消息
- `GET /api/v1/messages/unread-count` - 获取未读消息数

### 管理后台
- `GET /api/v1/statistics/textbooks` - 教材统计
- `GET /api/v1/statistics/users` - 用户统计
- `GET /api/v1/statistics/orders` - 订单统计
- `GET /api/v1/statistics/monthly` - 月度统计
- `GET /api/v1/statistics/export` - 导出报表

## 数据库配置

默认数据库配置:
- 数据库名: campus_trade
- 用户名: postgres
- 密码: postgres
- 端口: 5432

初始化数据库:
```bash
psql -U postgres -f backend/migrations/init.sql
```

## 开发注意事项

1. 修改后端代码后需要重新编译
2. 前端开发模式支持热更新
3. 上传的文件存储在 `uploads/` 目录
4. 导出的报表存储在 `exports/` 目录
5. JWT密钥请在生产环境中修改为安全的随机字符串

## 许可证

MIT License
