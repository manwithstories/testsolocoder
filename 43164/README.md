# 在线一对一家教预约与实时授课平台

一个完整的在线家教平台，使用Go (Gin) 作为后端，React + TypeScript 作为前端。

## 项目结构

```
tutoring-platform/
├── backend/                 # Go 后端
│   ├── config/              # 配置管理
│   ├── database/            # 数据库连接和迁移
│   ├── handlers/            # API 处理函数
│   ├── middleware/           # 中间件 (JWT认证, CORS)
│   ├── models/              # 数据模型
│   ├── routes/              # 路由定义
│   ├── main.go              # 应用入口
│   ├── go.mod               # Go 依赖
│   └── .env                 # 环境变量
└── frontend/                # React 前端
    ├── src/
    │   ├── components/      # UI 组件
    │   ├── layouts/         # 布局组件
    │   ├── pages/           # 页面组件
    │   ├── services/        # API 服务
    │   ├── store/           # 状态管理
    │   ├── types/           # TypeScript 类型
    │   ├── App.tsx          # 应用组件
    │   └── main.tsx         # 应用入口
    ├── index.html
    ├── package.json
    ├── tailwind.config.js
    ├── tsconfig.json
    └── vite.config.ts
```

## 核心功能

### 1. 家教管理
- 老师创建详细档案（教授科目、时薪、可预约时段、学历证书）
- 管理员审核后才能上架
- 老师评分和评价管理

### 2. 学生管理
- 学生注册和个人资料管理
- 设置学习目标
- 学科水平评估问卷
- 智能匹配老师

### 3. 课程预约系统
- 按科目和时间筛选老师
- 预约时间段
- 改期和取消功能（24小时通知期限）
- 预约时间冲突校验

### 4. 实时视频授课
- 集成第三方视频API
- 自动记录课程时长
- 监控连接质量
- 视频异常降级机制
- 中断时自动重连

### 5. 钱包与支付管理
- 多种支付方式
- 老师收入统计
- 平台佣金处理（默认10%）
- 提现申请
- 交易记录查询

### 6. 学习进度追踪
- 课程笔记
- 课后作业布置和提交
- 学生反馈收集
- 学习里程碑展示

### 7. 评价系统
- 每节课结束后双方互评
- 老师可以回复评价
- 综合评分展示
- 评价统计

### 8. 站内消息系统
- 师生实时沟通
- 文件共享
- 消息通知

## 技术特性

### 后端
- **框架**: Gin (Go Web Framework)
- **数据库**: PostgreSQL
- **ORM**: GORM
- **认证**: JWT (JSON Web Token)
- **权限**: 基于角色的访问控制 (RBAC)
- **时区处理**: 支持全球用户时区
- **日志**: 系统日志和交易流水

### 前端
- **框架**: React 18
- **语言**: TypeScript
- **构建工具**: Vite
- **样式**: Tailwind CSS
- **状态管理**: Zustand
- **路由**: React Router v6
- **图表**: Recharts
- **表单**: React Hook Form

## 快速开始

### 后端设置

1. 安装依赖:
```bash
cd backend
go mod download
```

2. 配置环境变量 (.env):
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=tutoring_db
SERVER_PORT=8080
JWT_SECRET=your-secret-key
```

3. 运行服务器:
```bash
go run main.go
```

### 前端设置

1. 安装依赖:
```bash
cd frontend
npm install
```

2. 启动开发服务器:
```bash
npm run dev
```

3. 打开浏览器访问: http://localhost:3000

## API 端点

### 认证
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录

### 用户
- `GET /api/users/me` - 获取当前用户资料
- `PUT /api/users/me` - 更新用户资料
- `GET /api/users/:id` - 获取用户详情

### 老师
- `GET /api/teachers` - 老师列表
- `GET /api/teachers/:id` - 老师详情
- `GET /api/teachers/:id/reviews` - 老师评价
- `GET /api/teacher/profile` - 老师个人资料
- `PUT /api/teacher/profile` - 更新老师资料
- `POST /api/teacher/subjects` - 添加科目
- `POST /api/teacher/availability` - 添加可预约时段

### 学生
- `GET /api/student/profile` - 学生个人资料
- `PUT /api/student/profile` - 更新学生资料
- `POST /api/student/goals` - 添加学习目标
- `GET /api/student/match-teachers` - 智能匹配老师
- `POST /api/student/assessment/submit` - 提交评估

### 课程预约
- `GET /api/bookings` - 获取预约列表
- `POST /api/bookings` - 创建预约
- `POST /api/bookings/:id/confirm` - 确认预约
- `POST /api/bookings/reschedule` - 改期
- `POST /api/bookings/cancel` - 取消预约
- `POST /api/bookings/:id/complete` - 完成课程

### 视频
- `POST /api/video/sessions` - 创建视频会话
- `POST /api/video/sessions/start` - 开始会话
- `POST /api/video/sessions/end` - 结束会话
- `GET /api/video/sessions/:id/quality` - 获取质量报告

### 钱包
- `GET /api/wallet` - 获取钱包
- `POST /api/wallet/deposit` - 充值
- `POST /api/wallet/withdraw` - 提现
- `GET /api/wallet/transactions` - 交易记录

### 评价
- `GET /api/reviews` - 获取评价列表
- `POST /api/reviews` - 创建评价
- `POST /api/reviews/:id/reply` - 回复评价

### 消息
- `GET /api/messages` - 获取消息
- `POST /api/messages` - 发送消息
- `GET /api/messages/conversations` - 获取对话列表

### 管理员
- `GET /api/admin/stats` - 统计数据
- `GET /api/admin/pending-approvals` - 待审核列表
- `POST /api/admin/teachers/:id/approve` - 审核通过
- `POST /api/admin/teachers/:id/reject` - 审核拒绝
- `POST /api/admin/withdraw-requests/:id/process` - 处理提现

## 数据库模型

主要数据模型包括:
- **User**: 用户信息
- **TeacherProfile**: 老师资料
- **StudentProfile**: 学生资料
- **Subject**: 科目
- **AvailabilitySlot**: 可预约时段
- **Booking**: 预约记录
- **VideoSession**: 视频会话
- **Wallet**: 钱包
- **Transaction**: 交易记录
- **Review**: 评价
- **Message**: 消息
- **Notification**: 通知
- **LessonNote**: 课程笔记
- **Homework**: 作业
- **LearningGoal**: 学习目标
- **Milestone**: 里程碑

## 安全特性

- JWT 认证
- 密码加密 (bcrypt)
- 基于角色的权限控制
- CORS 处理
- 输入验证
- SQL 注入防护 (GORM)

## 时区处理

- 所有时间存储为 UTC
- 用户可设置个人时区
- 显示时自动转换为用户时区

## 视频异常降级

- 自动检测连接质量
- 质量评分系统
- 重连机制
- 详细的连接日志

## 交易日志

- 完整的交易流水
- 平台佣金记录
- 提现申请跟踪
- 支付状态更新
