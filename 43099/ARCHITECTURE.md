# 场地预约与设备租赁管理系统 - 架构设计

## 技术栈

### 后端
- **语言**: Go 1.21+
- **框架**: Gin v1.9.0
- **数据库**: MySQL 8.0+
- **ORM**: GORM v1.26.1
- **认证**: JWT (golang-jwt/jwt/v5)
- **密码加密**: bcrypt
- **邮件**: gomail
- **Excel导出**: Excelize

### 前端
- **框架**: React 18+
- **语言**: TypeScript 5+
- **构建工具**: Vite 5+
- **状态管理**: Zustand
- **路由**: React Router v6
- **UI组件**: Ant Design 5+
- **HTTP请求**: Axios
- **日历**: FullCalendar
- **图表**: ECharts

## 目录结构

```
/
├── backend/                    # 后端Go项目
│   ├── cmd/
│   │   └── server/             # 应用入口
│   ├── internal/
│   │   ├── api/                # API路由和处理器
│   │   │   ├── handler/        # 请求处理器
│   │   │   └── middleware/     # 中间件
│   │   ├── service/            # 业务逻辑层
│   │   ├── repository/         # 数据访问层
│   │   ├── model/              # 数据模型
│   │   ├── dto/                # 数据传输对象
│   │   └── config/             # 配置
│   ├── pkg/                    # 公共包
│   │   ├── database/           # 数据库
│   │   ├── auth/               # 认证
│   │   ├── email/              # 邮件
│   │   ├── excel/              # Excel处理
│   │   └── logger/             # 日志
│   ├── uploads/                # 上传文件目录
│   ├── go.mod
│   ├── go.sum
│   └── .env.example
│
├── frontend/                   # 前端React项目
│   ├── src/
│   │   ├── api/                # API接口
│   │   ├── assets/             # 静态资源
│   │   ├── components/         # 公共组件
│   │   ├── pages/              # 页面组件
│   │   ├── store/              # 状态管理
│   │   ├── types/              # TypeScript类型
│   │   ├── utils/              # 工具函数
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── public/
│   ├── package.json
│   ├── tsconfig.json
│   └── vite.config.ts
│
└── docs/                       # 文档
    └── api.md
```

## 数据库设计

### 用户表 (users)
- id, username, email, password_hash, real_name, phone, avatar
- role: user/admin/super_admin
- email_verified: bool
- status: active/inactive
- created_at, updated_at

### 场地表 (venues)
- id, name, location, capacity, facilities, description, cover_image
- status: online/offline
- created_by, created_at, updated_at

### 场地价格表 (venue_prices)
- id, venue_id, day_of_week (0-6)
- time_slots: JSON [{start, end, price}]

### 设备分类表 (device_categories)
- id, name, description, sort_order

### 设备表 (devices)
- id, category_id, name, description, specification
- stock_quantity, available_quantity
- rental_price (per hour), deposit_amount
- status: online/offline
- images: JSON
- created_at, updated_at

### 预约订单表 (orders)
- id, order_no (唯一), user_id
- type: venue/device
- item_id (venue_id或device_id)
- start_time, end_time, total_hours
- total_amount, deposit_amount
- status: pending/confirmed/paid/completed/cancelled
- purpose, contact_name, contact_phone
- cancel_reason, cancelled_at
- reviewed_by, review_note, reviewed_at
- created_at, updated_at

### 支付记录表 (payments)
- id, order_id, transaction_no
- amount, payment_method: wechat/alipay/cash
- status: pending/success/failed/refunded
- paid_at, created_at

### 评价表 (reviews)
- id, order_id, user_id, rating (1-5)
- content, images: JSON
- status: pending/approved/rejected
- review_note, reviewed_by, reviewed_at
- created_at

### 操作日志表 (operation_logs)
- id, user_id, action, module, ip_address, user_agent
- details: JSON
- created_at

## API接口设计

### 认证相关
- POST /api/auth/register - 用户注册
- POST /api/auth/login - 用户登录
- POST /api/auth/logout - 登出
- POST /api/auth/verify-email - 邮箱验证
- POST /api/auth/forgot-password - 忘记密码

### 用户管理
- GET /api/users/me - 获取当前用户信息
- PUT /api/users/me - 更新个人信息
- GET /api/users - 获取用户列表 (admin)
- PUT /api/users/:id/role - 修改用户角色 (super_admin)

### 场地管理
- GET /api/venues - 获取场地列表
- POST /api/venues - 创建场地 (admin)
- GET /api/venues/:id - 获取场地详情
- PUT /api/venues/:id - 更新场地 (admin)
- DELETE /api/venues/:id - 删除场地 (admin)
- PATCH /api/venues/:id/status - 上下架 (admin)
- POST /api/venues/:id/prices - 设置价格 (admin)
- GET /api/venues/:id/availability - 查询可用时间段

### 设备管理
- GET /api/devices/categories - 获取设备分类
- POST /api/devices/categories - 创建分类 (admin)
- GET /api/devices - 获取设备列表
- POST /api/devices - 创建设备 (admin)
- PUT /api/devices/:id - 更新设备 (admin)
- POST /api/devices/batch-import - 批量导入 (admin)
- GET /api/devices/:id/availability - 查询设备可用库存

### 预约系统
- POST /api/bookings/check-availability - 检查可用性
- POST /api/bookings - 创建预约
- GET /api/bookings/calendar - 日历视图数据

### 订单管理
- GET /api/orders - 获取订单列表
- GET /api/orders/:id - 获取订单详情
- PUT /api/orders/:id/cancel - 取消预约
- PUT /api/orders/:id/confirm - 确认订单 (admin)
- PUT /api/orders/:id/complete - 标记完成 (admin)

### 支付管理
- GET /api/payments - 获取支付记录
- POST /api/payments/:id/confirm - 确认支付 (admin)
- GET /api/payments/export - 导出对账单 (admin)

### 评价管理
- POST /api/reviews - 创建评价
- GET /api/reviews - 获取评价列表
- PUT /api/reviews/:id/approve - 审核通过 (admin)
- PUT /api/reviews/:id/reject - 审核拒绝 (admin)

### 统计管理
- GET /api/stats/overview - 统计概览
- GET /api/stats/bookings - 预约统计
- GET /api/stats/revenue - 收入统计
- GET /api/stats/popular-venues - 热门场地

## 业务规则

1. **取消预约**: 预约开始前24小时可取消，超时不可取消
2. **时间冲突**: 同一场地同一时间段只能有一个有效预约
3. **库存检查**: 设备预约时检查可用库存，确认后扣减库存
4. **角色权限**:
   - 普通用户: 浏览、预约、取消、评价
   - 管理员: 场地/设备管理、订单审核、确认支付
   - 超级管理员: 用户管理、角色分配、所有管理员权限
5. **评价审核**: 用户提交的评价需要管理员审核后展示
