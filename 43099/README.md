# 场地预约与设备租赁管理系统

一个完整的场地预约与设备租赁管理的全栈应用。

## 技术栈

### 后端
- Go 1.21+
- Gin v1.9.1
- GORM v1.25.5
- MySQL 8.0+
- JWT 认证
- Excel 导出

### 前端
- React 18+
- TypeScript 5+
- Vite 5+
- Ant Design 5+
- Zustand (状态管理)
- React Router v6
- ECharts (图表)
- FullCalendar (日历)

## 功能特性

### 1. 用户体系
- 用户注册/登录
- JWT 认证
- 邮箱验证
- 角色权限（普通用户、管理员、超级管理员）
- 忘记密码/重置密码

### 2. 场地管理
- 场地 CRUD
- 场地价格设置（按星期、分时段定价）
- 场地上下架
- 场地可用时段查询

### 3. 设备管理
- 设备分类管理
- 设备 CRUD
- 库存管理
- 租赁单价/押金设置
- 批量导入

### 4. 预约系统
- 日历视图展示预约情况
- 按日期时段查询可用性
- 提交预约（时间冲突检测、库存检查
- 生成订单号

### 5. 订单管理
- 订单列表/详情
- 取消预约（24小时前可取消
- 管理员审核确认订单
- 标记完成

### 6. 支付管理
- 支付记录管理
- 确认支付
- 导出对账单为 Excel

### 7. 评价系统
- 订单完成后可评价
- 打分评价
- 管理员审核评价

### 8. 数据统计
- 预约量/收入统计
- 热门场地排行
- 按时间段筛选

### 9. 操作日志
- 所有操作完整日志记录

## 项目结构

```
/
├── backend/                    # 后端 Go 项目
│   ├── cmd/server/            # 应用入口
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handler/    # 请求处理器
│   │   │   └── middleware/ # 中间件
│   │   ├── service/          # 业务逻辑层
│   │   ├── repository/     # 数据访问层
│   │   ├── model/            # 数据模型
│   │   ├── dto/              # 数据传输对象
│   │   └── config/         # 配置
│   ├── pkg/                    # 公共包
│   │   ├── database/       # 数据库
│   │   ├── auth/           # 认证
│   │   ├── email/           # 邮件
│   │   ├── excel/          # Excel 处理
│   │   └── logger/          # 日志
│   └── uploads/             # 上传文件目录
│
├── frontend/                  # 前端 React 项目
│   ├── src/
│   │   ├── api/              # API 接口
│   │   ├── components/    # 公共组件
│   │   ├── pages/         # 页面组件
│   │   ├── store/         # 状态管理
│   │   ├── types/         # TypeScript 类型
│   │   └── utils/          # 工具函数
│   └── ...
│
└── docs/                      # 文档
```

## 快速开始

### 后端启动

1. 配置数据库

```bash
cd backend
cp .env.example .env
# 修改 .env 中的数据库配置
```

2. 安装依赖

```bash
go mod download
```

3. 启动服务

```bash
go run cmd/server/main.go
```

服务默认运行在 http://localhost:8080

### 前端启动

```bash
cd frontend
npm install
npm run dev
```

前端默认运行在 http://localhost:3000

## API 文档

### 认证接口

- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/logout` - 登出
- `POST /api/auth/verify-email` - 邮箱验证
- `POST /api/auth/forgot-password` - 忘记密码
- `POST /api/auth/reset-password` - 重置密码

### 用户接口

- `GET /api/users/me` - 获取当前用户信息
- `PUT /api/users/me` - 更新个人信息
- `GET /api/users` - 获取用户列表 (admin)
- `PUT /api/users/:id/role` - 修改用户角色 (super_admin)

### 场地接口

- `GET /api/venues` - 获取场地列表
- `POST /api/venues` - 创建场地 (admin)
- `GET /api/venues/:id` - 获取场地详情
- `PUT /api/venues/:id` - 更新场地 (admin)
- `DELETE /api/venues/:id` - 删除场地 (admin)
- `PATCH /api/venues/:id/status` - 上下架 (admin)
- `POST /api/venues/:id/prices` - 设置价格 (admin)
- `GET /api/venues/:id/availability` - 查询可用时间段

### 设备接口

- `GET /api/devices/categories` - 获取设备分类
- `POST /api/devices/categories` - 创建分类 (admin)
- `GET /api/devices` - 获取设备列表
- `POST /api/devices` - 创建设备 (admin)
- `PUT /api/devices/:id` - 更新设备 (admin)
- `POST /api/devices/batch-import` - 批量导入 (admin)

### 预约接口

- `POST /api/bookings` - 创建预约
- `GET /api/bookings/calendar` - 日历视图数据

### 订单接口

- `GET /api/orders` - 获取订单列表
- `GET /api/orders/:id` - 获取订单详情
- `PUT /api/orders/:id/cancel` - 取消预约
- `PUT /api/orders/:id/confirm` - 确认订单 (admin)
- `PUT /api/orders/:id/complete` - 标记完成 (admin)

### 支付接口

- `GET /api/payments` - 获取支付记录 (admin)
- `POST /api/payments/:id/confirm` - 确认支付 (admin)
- `GET /api/payments/export` - 导出对账单 (admin)

### 评价接口

- `POST /api/reviews` - 创建评价
- `GET /api/reviews` - 获取评价列表
- `PUT /api/reviews/:id/approve` - 审核通过 (admin)
- `PUT /api/reviews/:id/reject` - 审核拒绝 (admin)

### 统计接口

- `GET /api/stats/overview` - 统计概览 (admin)
- `GET /api/stats/bookings` - 预约统计 (admin)
- `GET /api/stats/revenue` - 收入统计 (admin)
- `GET /api/stats/popular-venues` - 热门场地 (admin)

## 业务规则

1. **取消预约**: 预约开始前 24 小时可取消，超时不可取消
2. **时间冲突**: 同一场地同一时间段只能有一个有效预约
3. **库存检查**: 设备预约时检查可用库存，确认后扣减库存
4. **角色权限**:
   - 普通用户: 浏览、预约、取消、评价
   - 管理员: 场地/设备管理、订单审核、确认支付
   - 超级管理员: 用户管理、角色分配、所有管理员权限
5. **评价审核**: 用户提交的评价需要管理员审核后展示

## 数据库表

- users - 用户表
- venues - 场地表
- venue_prices - 场地价格表
- device_categories - 设备分类表
- devices - 设备表
- orders - 预约订单表
- payments - 支付记录表
- reviews - 评价表
- operation_logs - 操作日志表

## License

MIT
