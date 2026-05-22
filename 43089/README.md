# 旅行规划与行程管理系统

一个功能完整的个人旅行规划与行程管理Web应用，使用Go + Gin后端和React + TypeScript前端构建。

## ✨ 功能特性

### 核心功能

1. **旅行计划管理**
   - 创建、编辑、删除旅行计划
   - 设置目的地、日期范围、预算上限
   - 多人协作编辑，支持角色权限控制
   - 计划状态管理（草稿、进行中、已完成）

2. **行程安排**
   - 按天规划具体行程
   - 支持多种活动类型：景点、交通、住宿、餐饮、其他
   - 活动包含时间段、地点、费用、备注
   - 活动排序和时间冲突检查

3. **预算追踪**
   - 实时统计各类别支出
   - 对比计划预算和实际花费
   - 预算使用率可视化
   - 超支预警提醒

4. **文件附件管理**
   - 上传机票、酒店订单等PDF文件
   - 支持图片上传
   - 文件分类管理（机票、酒店、保险、收据等）
   - 在线预览和下载

5. **清单管理**
   - 出发前准备清单
   - 行李清单
   - 支持勾选完成状态
   - 进度可视化

6. **数据导出**
   - 导出为JSON格式
   - 导出为PDF格式
   - 方便分享和备份

7. **地图集成**
   - 行程地点标记在地图上
   - 可视化路线
   - 地点列表和详情

8. **提醒功能**
   - 重要活动提前通知
   - 支持邮件和应用内通知
   - 可关联特定活动

### 技术特性

- **后端**：Go 1.21 + Gin v1.9.1 + GORM
- **数据库**：PostgreSQL 15+
- **认证**：JWT Token认证
- **权限**：RBAC角色权限控制
- **并发控制**：乐观锁防止数据覆盖
- **审计日志**：所有操作可追溯
- **前端**：React 18 + TypeScript + Vite 5
- **状态管理**：Redux Toolkit
- **样式**：Tailwind CSS 3
- **表单验证**：React Hook Form + Zod
- **响应式设计**：支持PC和移动端

## 📁 项目结构

```
travel-planner/
├── backend/                    # 后端Go项目
│   ├── cmd/
│   │   └── main.go            # 应用入口
│   ├── configs/
│   │   └── config.yaml        # 配置文件
│   ├── internal/
│   │   ├── config/            # 配置加载
│   │   ├── database/          # 数据库连接
│   │   ├── logger/            # 日志配置
│   │   ├── middleware/        # 中间件（JWT、RBAC、CORS、审计）
│   │   ├── models/            # 数据模型
│   │   ├── handlers/          # API处理器
│   │   ├── routes/            # 路由配置
│   │   └── utils/             # 工具函数
│   ├── migrations/            # 数据库迁移脚本
│   ├── uploads/               # 文件上传目录
│   └── go.mod
├── frontend/                   # 前端React项目
│   ├── src/
│   │   ├── components/        # 公共组件
│   │   ├── pages/             # 页面组件
│   │   ├── routes/            # 路由配置
│   │   ├── services/          # API服务
│   │   ├── store/             # Redux状态管理
│   │   ├── types/             # TypeScript类型定义
│   │   └── main.tsx           # 应用入口
│   ├── index.html
│   ├── package.json
│   ├── vite.config.ts
│   └── tailwind.config.js
├── docker-compose.yml          # Docker编排
└── README.md
```

## 🚀 快速开始

### 前置要求

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+ 或 Docker

### 方式一：使用Docker（推荐）

1. 启动PostgreSQL服务：
```bash
docker-compose up -d
```

2. 初始化数据库：
```bash
cd backend
psql -h localhost -U postgres -d travel_planner -f migrations/001_init_schema.sql
```

### 方式二：手动安装

1. 创建PostgreSQL数据库：
```sql
CREATE DATABASE travel_planner;
CREATE USER travel_user WITH PASSWORD 'travel_password';
GRANT ALL PRIVILEGES ON DATABASE travel_planner TO travel_user;
```

2. 执行数据库迁移脚本。

### 启动后端服务

```bash
cd backend
go mod download
go run cmd/main.go
```

后端服务将在 `http://localhost:8080` 启动。

### 启动前端服务

```bash
cd frontend
npm install
npm run dev
```

前端服务将在 `http://localhost:5173` 启动。

## 🔧 配置说明

### 后端配置 (backend/configs/config.yaml)

```yaml
server:
  port: 8080
  mode: debug

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: travel_planner
  sslmode: disable

jwt:
  secret: your-secret-key-change-in-production
  expires_in: 24h

upload:
  path: ./uploads
  max_size: 10485760  # 10MB
  allowed_types:
    - application/pdf
    - image/jpeg
    - image/png
    - image/gif

cors:
  allowed_origins:
    - http://localhost:5173
```

### 前端配置

前端API代理配置在 `vite.config.ts` 中，默认代理 `/api` 到 `http://localhost:8080`。

## 📡 API文档

### 认证接口

- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `GET /api/auth/me` - 获取当前用户信息
- `PUT /api/auth/profile` - 更新个人资料

### 旅行计划接口

- `GET /api/plans` - 获取计划列表
- `GET /api/plans/:id` - 获取计划详情
- `POST /api/plans` - 创建计划
- `PUT /api/plans/:id` - 更新计划
- `DELETE /api/plans/:id` - 删除计划
- `POST /api/plans/:id/participants` - 添加参与者
- `DELETE /api/plans/:id/participants/:pid` - 移除参与者

### 行程活动接口

- `GET /api/plans/:id/activities` - 获取活动列表
- `POST /api/plans/:id/activities` - 创建活动
- `PUT /api/plans/:id/activities/:aid` - 更新活动
- `DELETE /api/plans/:id/activities/:aid` - 删除活动

### 预算接口

- `GET /api/plans/:id/budget` - 获取预算汇总

### 文件接口

- `GET /api/plans/:id/files` - 获取文件列表
- `POST /api/plans/:id/files` - 上传文件
- `DELETE /api/plans/:id/files/:fid` - 删除文件
- `GET /api/files/download/:fid` - 下载文件

### 清单接口

- `GET /api/plans/:id/checklists` - 获取清单列表
- `POST /api/plans/:id/checklists` - 创建清单
- `PUT /api/plans/:id/checklists/:cid` - 更新清单
- `DELETE /api/plans/:id/checklists/:cid` - 删除清单
- `POST /api/plans/:id/checklists/:cid/items` - 添加清单项
- `PUT /api/plans/:id/checklists/:cid/items/:iid` - 更新清单项
- `DELETE /api/plans/:id/checklists/:cid/items/:iid` - 删除清单项

### 提醒接口

- `GET /api/reminders` - 获取提醒列表
- `POST /api/plans/:id/reminders` - 创建提醒
- `PUT /api/reminders/:rid` - 更新提醒
- `DELETE /api/reminders/:rid` - 删除提醒

### 导出接口

- `GET /api/plans/:id/export/json` - 导出JSON
- `GET /api/plans/:id/export/pdf` - 导出PDF

### 地图接口

- `GET /api/plans/:id/map` - 获取地图数据

## 🔐 权限说明

### 角色定义

- **所有者 (Owner)**：计划创建者，拥有所有权限
- **编辑者 (Editor)**：可以编辑计划内容和活动
- **查看者 (Viewer)**：只能查看计划内容

### 权限矩阵

| 操作 | 所有者 | 编辑者 | 查看者 |
|------|--------|--------|--------|
| 查看计划 | ✓ | ✓ | ✓ |
| 编辑计划 | ✓ | ✓ | ✗ |
| 删除计划 | ✓ | ✗ | ✗ |
| 管理参与者 | ✓ | ✗ | ✗ |
| 创建/编辑活动 | ✓ | ✓ | ✗ |
| 删除活动 | ✓ | ✓ | ✗ |
| 上传文件 | ✓ | ✓ | ✗ |
| 删除文件 | ✓ | ✓ | ✗ |
| 管理清单 | ✓ | ✓ | ✗ |
| 导出数据 | ✓ | ✓ | ✓ |

## 🛡️ 安全特性

1. **JWT认证**：无状态Token认证，支持Token过期刷新
2. **密码加密**：使用bcrypt加密存储用户密码
3. **输入验证**：所有API输入参数进行严格验证
4. **SQL注入防护**：使用GORM预编译SQL
5. **CORS保护**：配置跨域白名单
6. **文件上传安全**：文件类型白名单、大小限制、随机文件名
7. **审计日志**：记录所有用户操作，包含IP、用户ID、操作内容
8. **乐观锁**：防止并发编辑时数据覆盖

## 📊 数据库设计

主要数据表：

- `users` - 用户表
- `roles` - 角色表
- `permissions` - 权限表
- `travel_plans` - 旅行计划表
- `plan_participants` - 计划参与者表
- `activities` - 行程活动表
- `files` - 文件附件表
- `checklists` - 清单表
- `checklist_items` - 清单项表
- `reminders` - 提醒表
- `audit_logs` - 审计日志表

## 🎯 前端页面

- `/login` - 登录页
- `/register` - 注册页
- `/dashboard` - 仪表盘
- `/plans` - 计划列表
- `/plans/:id` - 计划详情
- `/plans/:id/activities` - 行程安排
- `/plans/:id/budget` - 预算追踪
- `/plans/:id/files` - 文件管理
- `/plans/:id/checklist` - 清单管理
- `/plans/:id/map` - 地图视图
- `/reminders` - 提醒管理

## 🤝 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📝 许可证

本项目采用 MIT 许可证。

## 📧 联系方式

如有问题或建议，请通过Issue反馈。
