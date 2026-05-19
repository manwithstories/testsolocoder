# 活动组织与报名管理系统

一个全栈活动管理应用，后端使用 Go + Gin 框架，前端使用 React + TypeScript + Vite。

## 功能特性

### 用户模块
- 用户注册（邮箱验证）
- 用户登录（JWT认证）
- 密码使用 bcrypt 加密存储
- 邮箱验证码验证（模拟发送）

### 活动管理
- 活动创建、查看、编辑、删除
- 设置活动标题、描述、地点、时间、人数上限、报名截止时间
- 已截止的活动不允许报名
- 组织者可以查看报名名单

### 报名管理
- 用户可以报名活动
- 用户可以取消报名
- 活动人数满了自动关闭报名
- 取消后名额释放
- 同一用户对同一活动取消报名不能超过3次

### 导出功能
- 组织者可以导出报名名单为 CSV 文件

## 技术栈

### 后端
- Go 1.21+
- Gin Web 框架
- GORM ORM
- SQLite 数据库
- JWT 认证
- bcrypt 密码加密

### 前端
- React 18
- TypeScript
- Vite
- React Router
- Axios
- Tailwind CSS

## 快速开始

### 后端启动

```bash
cd backend

# 安装依赖
go mod tidy

# 构建并运行
go run main.go
# 或者
go build -o server . && ./server
```

后端服务将在 `http://localhost:8080` 启动

### 前端启动

```bash
cd frontend

# 安装依赖
npm install

# 开发模式
npm run dev

# 构建生产版本
npm run build
```

前端服务将在 `http://localhost:5173` 启动

## API 接口

### 认证接口
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/verify` - 邮箱验证

### 用户接口
- `GET /api/user/profile` - 获取用户信息
- `GET /api/user/registrations` - 获取我的报名记录

### 活动接口
- `GET /api/events` - 获取活动列表
- `GET /api/events/:id` - 获取活动详情
- `POST /api/events` - 创建活动
- `PUT /api/events/:id` - 更新活动
- `DELETE /api/events/:id` - 删除活动
- `POST /api/events/:id/register` - 报名活动
- `POST /api/events/:id/cancel` - 取消报名
- `GET /api/events/:id/registrations` - 获取活动报名名单（组织者）
- `GET /api/events/:id/export` - 导出报名名单CSV（组织者）

## 项目结构

```
.
├── backend/                    # 后端代码
│   ├── config/                # 配置
│   ├── controllers/           # 控制器
│   ├── middleware/            # 中间件
│   ├── models/                # 数据模型
│   ├── routes/                # 路由
│   ├── utils/                 # 工具函数
│   ├── main.go                # 入口文件
│   └── go.mod
└── frontend/                  # 前端代码
    ├── src/
    │   ├── components/        # 组件
    │   ├── context/           # React Context
    │   ├── pages/             # 页面组件
    │   ├── services/          # API 服务
    │   ├── types/             # TypeScript 类型
    │   ├── App.tsx
    │   └── main.tsx
    └── package.json
```

## 使用说明

1. 注册账户：输入邮箱、用户名、密码，系统会生成验证码（在后端控制台输出）
2. 验证邮箱：输入收到的验证码完成验证
3. 登录系统：使用邮箱和密码登录
4. 浏览活动：查看所有活动，可以报名感兴趣的活动
5. 发布活动：点击"发布活动"创建新活动
6. 管理活动：作为组织者可以查看报名名单、导出CSV、删除活动
7. 我的报名：查看自己的报名记录和状态

## 注意事项

- 邮箱验证是模拟实现的，验证码会打印在后端控制台
- JWT Token 存储在 localStorage 中
- 活动报名截止时间必须早于活动开始时间
- 取消报名超过3次后将无法再报名该活动
