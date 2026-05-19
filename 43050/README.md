# 群组账单分摊 Web 应用

一个功能完整的群组账单分摊应用，使用 Go Gin 框架后端和 React + TypeScript 前端构建。

## 功能特性

### 1. 群组管理
- 创建群组，每个群组有独立的账本
- 邀请码机制，成员通过邀请码加入群组
- 成员可随时退出群组，历史记录保留
- 查看群组成员列表

### 2. 消费记录
- 添加群组消费，选择付款人和参与成员
- 支持三种分摊方式：
  - **均摊**：所有参与成员平均分摊
  - **按比例**：按设定的比例分摊
  - **自定义金额**：手动输入每人承担的金额
- 编辑和删除账单记录
- 分摊预览功能

### 3. 智能结算
- 自动计算最优转账路径
- 使用贪心算法简化债务网络
- 最少转账笔数，避免循环欠款
- 清晰展示谁应该给谁转账

### 4. 历史统计
- 个人收支总览（总支出、待还金额、净余额）
- 月度收支趋势图（柱状图）
- 消费分类占比（饼图）
- 收支明细列表
- 支持按时间段和群组筛选

### 5. 技术特性
- JWT 认证机制
- SQLite 数据库存储
- 乐观锁确保数据一致性
- 数据库事务保证操作原子性
- 响应式前端界面
- RESTful API 设计

## 技术栈

### 后端
- **语言**: Go 1.21+
- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: SQLite
- **认证**: JWT (golang-jwt)
- **密码加密**: bcrypt

### 前端
- **语言**: TypeScript
- **框架**: React 18
- **构建工具**: Vite
- **路由**: React Router DOM
- **HTTP 客户端**: Axios
- **图表**: Recharts
- **样式**: Tailwind CSS
- **日期处理**: Day.js

## 项目结构

```
.
├── backend/                    # 后端 Go 项目
│   ├── cmd/
│   │   └── server/
│   │       └── main.go        # 应用入口
│   ├── internal/
│   │   ├── database/          # 数据库初始化
│   │   ├── models/            # 数据模型
│   │   ├── handlers/          # API 处理器
│   │   ├── middleware/        # 中间件 (JWT认证)
│   │   ├── utils/             # 工具函数
│   │   └── settlement/        # 智能结算算法
│   ├── .env                   # 环境变量
│   ├── go.mod
│   └── go.sum
└── frontend/                   # 前端 React 项目
    ├── src/
    │   ├── api/               # API 接口定义
    │   ├── components/        # 公共组件
    │   ├── context/           # React Context
    │   ├── pages/             # 页面组件
    │   ├── types/             # TypeScript 类型定义
    │   ├── App.tsx
    │   └── main.tsx
    ├── index.html
    ├── package.json
    ├── tailwind.config.js
    ├── tsconfig.json
    └── vite.config.ts
```

## 快速开始

### 前置要求
- Go 1.21 或更高版本
- Node.js 18 或更高版本
- npm 或 yarn

### 后端启动

1. 进入后端目录：
```bash
cd backend
```

2. 安装依赖：
```bash
go mod tidy
```

3. 配置环境变量（可选，默认值已设置）：
```bash
# 复制 .env 文件并修改
cp .env.example .env
```

4. 启动后端服务：
```bash
go run cmd/server/main.go
```

后端服务将在 `http://localhost:8080` 启动。

### 前端启动

1. 进入前端目录：
```bash
cd frontend
```

2. 安装依赖：
```bash
npm install
```

3. 启动开发服务器：
```bash
npm run dev
```

前端服务将在 `http://localhost:5173` 启动。

4. 构建生产版本：
```bash
npm run build
```

## API 文档

### 认证接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | /api/auth/register | 用户注册 | 否 |
| POST | /api/auth/login | 用户登录 | 否 |
| GET | /api/auth/me | 获取当前用户信息 | 是 |

### 群组接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | /api/groups | 创建群组 | 是 |
| GET | /api/groups | 获取用户群组列表 | 是 |
| GET | /api/groups/:id | 获取群组详情 | 是 |
| POST | /api/groups/join | 通过邀请码加入群组 | 是 |
| POST | /api/groups/:id/leave | 退出群组 | 是 |
| GET | /api/groups/:id/members | 获取群组成员 | 是 |

### 账单接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | /api/groups/:groupId/expenses | 添加账单 | 是 |
| GET | /api/groups/:groupId/expenses | 获取群组账单列表 | 是 |
| GET | /api/expenses/:id | 获取账单详情 | 是 |
| PUT | /api/expenses/:id | 更新账单 | 是 |
| DELETE | /api/expenses/:id | 删除账单 | 是 |

### 结算接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/groups/:groupId/balances | 获取群组余额 | 是 |
| GET | /api/groups/:groupId/transfers | 获取最优转账方案 | 是 |
| GET | /api/groups/:groupId/stats | 获取群组成员统计 | 是 |

### 统计接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/stats/summary | 获取用户统计概览 | 是 |
| GET | /api/stats/monthly | 获取月度统计 | 是 |
| GET | /api/stats/history | 获取历史账单 | 是 |

## 智能结算算法说明

智能结算使用贪心算法来最小化转账笔数：

1. 计算每个用户的净余额（已支付 - 应承担）
2. 将用户分为债权人（余额 > 0）和债务人（余额 < 0）
3. 按金额从大到小排序
4. 每次匹配最大的债务人和最大的债权人，进行转账
5. 重复直到所有债务结清

这种算法可以确保用最少的转账次数完成结算，避免了循环转账的问题。

## 数据一致性保障

- **数据库事务**: 所有写操作使用事务确保原子性
- **乐观锁**: Expense 表使用 Version 字段实现乐观锁
- **输入验证**: 服务端对所有输入进行严格验证
- **权限检查**: 每个操作都验证用户是否有权限访问资源

## 注意事项

1. **JWT 密钥**: 生产环境请务必修改 `JWT_SECRET` 环境变量
2. **数据库**: SQLite 适合小型应用，大规模使用建议迁移到 PostgreSQL 或 MySQL
3. **CORS**: 当前配置允许所有来源，生产环境应限制为特定域名
4. **邀请码**: 邀请码是全局唯一的，长度为 8 位大写字母和数字组合

## 开发建议

- 后端使用 `air` 工具进行热重载开发
- 前端 Vite 已配置代理到后端 API
- 数据库文件 `splitwise.db` 会在首次运行时自动创建
- 可以使用 Postman 或类似工具测试 API

## License

MIT
