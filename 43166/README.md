# 企业工商注册代办平台

一个完整的企业工商注册代办与公司设立服务平台，基于 Go + Gin + Vue3 + TypeScript 技术栈开发。

## 项目结构

```
.
├── server/                 # 后端服务 (Go + Gin)
│   ├── config/             # 配置文件
│   ├── controllers/        # 控制器层
│   ├── database/           # 数据库初始化
│   ├── middleware/          # 中间件
│   ├── models/             # 数据模型
│   ├── routes/             # 路由配置
│   ├── services/           # 业务逻辑层
│   ├── utils/              # 工具函数
│   └── main.go             # 入口文件
├── client/                 # 前端应用 (Vue3 + TypeScript)
│   ├── src/
│   │   ├── api/            # API 接口
│   │   ├── components/     # 公共组件
│   │   ├── layouts/        # 布局组件
│   │   ├── router/         # 路由配置
│   │   ├── store/          # 状态管理 (Pinia)
│   │   ├── styles/         # 全局样式
│   │   ├── types/          # TypeScript 类型定义
│   │   ├── utils/          # 工具函数
│   │   └── views/          # 页面视图
│   └── package.json
└── README.md
```

## 功能特性

### 用户管理
- 支持创业者、代办专员、平台管理员三种角色
- 基于 JWT 的认证授权
- 不同角色有不同的权限和数据访问范围

### 公司注册申请管理
- 创业者提交公司注册申请
- 填写公司名称、注册资本、股东信息等
- 支持上传身份证、营业执照等材料
- 申请状态跟踪和审核意见查看

### 代办流程管理
- 核名、工商登记、税务登记、银行开户等环节
- 多环节并行处理
- 流程状态实时更新
- 证明文件上传

### 代办专员管理
- 专员信息管理（工号、专业领域、工作时间等）
- 自动分配申请
- 工作统计和绩效评估

### 申请费用管理
- 自动计算代办费用
- 多种支付方式
- 费用标准和优惠策略管理
- 费用统计报表

### 进度通知系统
- 站内消息、邮件、短信通知
- 消息中心查看历史通知
- 通知模板管理

### 统计分析
- 申请总量、处理时长分布
- 专员工作效率、费用收入统计
- 多维度数据筛选
- 图表展示

### 数据导出
- 申请列表、费用明细、专员记录导出
- 异步处理避免长时间等待
- 下载链接有效期管理

## 技术栈

### 后端
- **框架**: Gin v1.9
- **ORM**: GORM v1.25
- **数据库**: MySQL
- **认证**: JWT (golang-jwt/jwt/v5)
- **密码加密**: bcrypt
- **Excel导出**: excelize

### 前端
- **框架**: Vue 3.4
- **语言**: TypeScript
- **状态管理**: Pinia
- **UI组件库**: Element Plus
- **路由**: Vue Router 4
- **HTTP客户端**: Axios
- **图表**: ECharts

## 快速开始

### 前置要求
- Go 1.21+
- Node.js 18+
- MySQL 5.7+

### 后端启动

1. 创建数据库并执行初始化脚本
```bash
mysql -u root -p < server/database/init.sql
```

2. 修改配置文件 `server/config/config.yaml`，设置数据库连接信息

3. 安装依赖
```bash
cd server
go mod download
```

4. 启动服务
```bash
go run main.go
```

服务默认运行在 `http://localhost:8080`

### 前端启动

1. 安装依赖
```bash
cd client
npm install
```

2. 启动开发服务器
```bash
npm run dev
```

访问 `http://localhost:3000`

### 默认账户
- 管理员: admin / admin123

## API 文档

所有 API 接口统一返回格式：

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

### 认证相关
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/logout` - 退出登录

### 申请管理
- `GET /api/v1/applications` - 获取申请列表
- `POST /api/v1/applications` - 创建申请
- `GET /api/v1/applications/:id` - 获取申请详情
- `PUT /api/v1/applications/:id` - 更新申请
- `POST /api/v1/applications/:id/submit` - 提交申请
- `POST /api/v1/applications/:id/cancel` - 取消申请

### 代办流程
- `GET /api/v1/applications/:applicationId/steps` - 获取流程步骤
- `PUT /api/v1/applications/:applicationId/steps/:stepId` - 更新步骤
- `POST /api/v1/applications/:applicationId/steps/:stepId/complete` - 完成步骤

### 费用管理
- `POST /api/v1/fees/calculate` - 计算费用
- `POST /api/v1/fees/pay` - 支付费用
- `GET /api/v1/fees` - 获取费用列表

### 通知管理
- `GET /api/v1/notifications` - 获取通知列表
- `PUT /api/v1/notifications/:id/read` - 标记已读
- `PUT /api/v1/notifications/read-all` - 全部标为已读

### 管理员功能
- `GET /api/v1/admin/agents` - 获取专员列表
- `POST /api/v1/admin/agents` - 创建专员
- `GET /api/v1/admin/statistics/overview` - 获取统计概览
- `POST /api/v1/exports` - 创建导出任务

## 开发计划

- [ ] 单元测试覆盖
- [ ] Docker 容器化部署
- [ ] CI/CD 流水线
- [ ] 性能优化
- [ ] 移动端适配

## 许可证

MIT License
