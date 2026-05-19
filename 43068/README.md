# Freelancer Management System

一个完整的自由职业者项目与发票管理全栈应用，包含用户认证、项目管理、工时记录、发票生成、客户管理和统计仪表盘功能。

## 技术栈

### 后端
- **Go 1.21+
- **Gin** - Web 框架
- **GORM** - ORM
- **SQLite** - 数据库
- **JWT** - 认证
- **gofpdf** - PDF 生成

### 前端
- **React 18**
- **TypeScript**
- **Vite**
- **Tailwind CSS**
- **React Router**
- **Zustand** - 状态管理
- **React Hook Form**
- **Recharts** - 图表库
- **Axios** - HTTP 客户端

## 核心功能

### 1. 用户认证
- 邮箱注册/登录
- JWT 令牌刷新
- 密码强度校验（8位以上，大小写字母+数字+特殊字符）

### 2. 项目管理
- 创建/编辑/删除项目
- 里程碑管理
- 截止日期设置
- 状态流转严格校验（草稿→进行中→已完成→已归档，不可跳步）

### 3. 工时记录
- 手动录入工时
- 计时器模式
- 同项目每日工时上限24小时

### 4. 发票管理
- 根据工时和费率自动计算金额
- PDF 导出
- 发票编号按年自增且唯一
- 并发控制

### 5. 客户管理
- 客户信息维护
- 合同信息
- 邮箱格式校验
- 有关联项目的客户不能删除

### 6. 统计仪表盘
- 月度收入统计
- 项目进度
- 逾期提醒
- 图表展示

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- npm 或 yarn

### 启动方式

#### 方式一：使用启动脚本 (推荐)

```bash
./scripts/start.sh
```

#### 方式二：手动启动

**启动后端：**

```bash
cd backend
go mod download
go run cmd/server/main.go
```

后端服务将在 http://localhost:8080 启动

**启动前端：**

```bash
cd frontend
npm install
npm run dev
```

前端服务将在 http://localhost:3000 启动

### 配置文件

后端配置文件位于 `backend/config.yaml`

```yaml
server:
  port: 8080
  mode: debug

database:
  driver: sqlite
  dsn: ./freelancer.db

jwt:
  secret: your-secret-key-change-in-production
  access_token_ttl: 3600
  refresh_token_ttl: 604800

email:
  smtp_host: smtp.example.com
  smtp_port: 587
  username: your-email@example.com
  password: your-password
  from_name: Freelancer Management
  from_address: no-reply@example.com

app:
  max_hours_per_day: 24
  invoice_prefix: INV
```

## API 文档

### 认证接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | /api/auth/register | 注册 | 否 |
| POST | /api/auth/login | 登录 | 否 |
| POST | /api/auth/refresh | 刷新令牌 | 否 |
| GET | /api/auth/me | 获取当前用户 | 是 |

### 客户管理

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/clients | 获取客户列表 | 是 |
| POST | /api/clients | 创建客户 | 是 |
| GET | /api/clients/:id | 获取客户详情 | 是 |
| PUT | /api/clients/:id | 更新客户 | 是 |
| DELETE | /api/clients/:id | 删除客户 | 是 |

### 项目管理

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/projects | 获取项目列表 | 是 |
| POST | /api/projects | 创建项目 | 是 |
| GET | /api/projects/:id | 获取项目详情 | 是 |
| PUT | /api/projects/:id | 更新项目 | 是 |
| DELETE | /api/projects/:id | 删除项目 | 是 |
| POST | /api/projects/:id/milestones | 添加里程碑 | 是 |
| PUT | /api/projects/milestones/:milestone_id | 更新里程碑 | 是 |
| DELETE | /api/projects/milestones/:milestone_id | 删除里程碑 | 是 |

### 工时记录

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/time-entries | 获取工时列表 | 是 |
| POST | /api/time-entries | 创建工时记录 | 是 |
| GET | /api/time-entries/active-timer | 获取活动计时器 | 是 |
| POST | /api/time-entries/timer/start | 开始计时 | 是 |
| POST | /api/time-entries/timer/:id/stop | 停止计时 | 是 |
| PUT | /api/time-entries/:id | 更新工时记录 | 是 |
| DELETE | /api/time-entries/:id | 删除工时记录 | 是 |

### 发票管理

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/invoices | 获取发票列表 | 是 |
| POST | /api/invoices | 创建发票 | 是 |
| GET | /api/invoices/:id | 获取发票详情 | 是 |
| PUT | /api/invoices/:id/status | 更新发票状态 | 是 |
| GET | /api/invoices/:id/download | 下载发票PDF | 是 |
| DELETE | /api/invoices/:id | 删除发票 | 是 |

### 统计仪表盘

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/dashboard | 获取统计数据 | 是 |

## 项目结构

```
.
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go          # 入口文件
│   ├── internal/
│   │   ├── auth/             # 认证相关
│   │   ├── config/           # 配置管理
│   │   ├── database/         # 数据库连接
│   │   ├── handlers/       # API 处理器
│   │   ├── logger/         # 日志
│   │   ├── middleware/     # 中间件
│   │   ├── models/         # 数据模型
│   │   ├── router/         # 路由
│   │   └── utils/          # 工具函数
│   ├── pkg/
│   │   ├── email/          # 邮件服务
│   │   ├── jwt/            # JWT 服务
│   │   └── pdf/            # PDF 生成
│   ├── config.yaml           # 配置文件
│   └── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── api/             # API 客户端
│   │   ├── components/      # 组件
│   │   ├── pages/           # 页面
│   │   ├── store/           # 状态管理
│   │   └── types/           # 类型定义
│   │   └── utils/          # 工具函数
│   ├── index.html
│   ├── package.json
│   ├── tsconfig.json
│   └── vite.config.ts
│   └── tailwind.config.js
└── scripts/
    └── start.sh            # 启动脚本
```

## 安全特性

- JWT 认证与刷新
- 密码哈希存储 (bcrypt)
- CORS 中间件
- 请求参数校验
- 统一错误处理
- 操作日志记录
- 数据库事务
- 并发控制（发票编号生成）

## 日志

系统会自动创建日志文件：

- `backend/logs/app.log` - 应用日志
- `backend/logs/access.log` - 访问日志

## License

MIT
