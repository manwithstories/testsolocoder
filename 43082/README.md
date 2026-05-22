# 多商家在线商城平台

一个基于 Go + Gin + React + TypeScript + PostgreSQL 构建的多商家在线商城平台，主要面向小商家和个人卖家。

## 技术栈

### 后端
- **语言**: Go 1.21+
- **框架**: Gin v1.9.1
- **ORM**: GORM v1.25.5
- **数据库**: PostgreSQL
- **认证**: JWT
- **密码加密**: bcrypt
- **Excel导出**: excelize

### 前端
- **框架**: React 18 + TypeScript 5
- **构建工具**: Vite 5
- **路由**: React Router v6
- **状态管理**: Zustand
- **UI组件**: Ant Design 5
- **HTTP客户端**: Axios
- **图表**: ECharts

## 项目结构

```
├── backend/                # 后端项目
│   ├── cmd/server/         # 服务入口
│   ├── internal/
│   │   ├── config/         # 配置管理
│   │   ├── database/       # 数据库连接
│   │   ├── models/         # 数据模型
│   │   ├── middleware/     # 中间件
│   │   ├── handlers/       # API处理器
│   │   ├── services/       # 业务逻辑
│   │   ├── dto/            # 数据传输对象
│   │   ├── routes/         # 路由配置
│   │   └── utils/          # 工具函数
│   └── pkg/                # 公共包
│       ├── auth/           # JWT认证
│       └── errors/         # 错误处理
└── frontend/               # 前端项目
    ├── src/
    │   ├── api/            # API封装
    │   ├── components/     # 通用组件
    │   ├── pages/          # 页面组件
    │   ├── store/          # 状态管理
    │   ├── types/          # TypeScript类型
    │   └── styles/         # 样式文件
    └── public/             # 静态资源
```

## 核心功能

### 1. 商家入驻与管理
- 商家申请入驻
- 管理员审核
- 商家发布和管理商品
- 独立店铺主页

### 2. 商品管理
- 多规格商品（颜色、尺码等）
- 库存管理
- 商品分类上架
- 图片上传

### 3. 购物车与订单系统
- 购物车功能
- 下单、支付、退款
- 订单状态流转：待付款 → 待发货 → 已发货 → 已完成 → 已退款
- 库存并发控制（防止超卖）

### 4. 评价系统
- 商品评分和评价
- 商家回复评价

### 5. 店铺收藏与搜索
- 店铺收藏
- 按店铺名、商品名、分类、价格区间搜索筛选
- 排序和分页

### 6. 通知模块
- 商家新订单通知
- 买家发货和退款通知

### 7. 管理员功能
- 审核商家入驻
- 管理平台分类
- 处理纠纷投诉

### 8. 数据统计与导出
- 商家店铺销售额、订单量趋势
- 管理员平台整体统计
- 支持导出Excel报表

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- PostgreSQL 13+

### 后端启动

1. 进入后端目录
```bash
cd backend
```

2. 复制环境变量配置
```bash
cp .env.example .env
```

3. 修改 `.env` 文件中的数据库配置

4. 安装依赖
```bash
go mod download
```

5. 启动服务
```bash
go run cmd/server/main.go
```

后端服务将在 `http://localhost:8080` 启动

### 前端启动

1. 进入前端目录
```bash
cd frontend
```

2. 安装依赖
```bash
npm install
```

3. 启动开发服务器
```bash
npm run dev
```

前端服务将在 `http://localhost:5173` 启动

## API 文档

### 认证接口
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `GET /api/auth/me` - 获取当前用户信息

### 店铺接口
- `POST /api/shops/apply` - 申请入驻
- `GET /api/shops/:id` - 获取店铺详情
- `GET /api/shops` - 获取店铺列表
- `PUT /api/shops/:id` - 更新店铺信息

### 商品接口
- `POST /api/products` - 创建商品
- `GET /api/products` - 获取商品列表（支持搜索、筛选、排序、分页）
- `GET /api/products/:id` - 获取商品详情
- `PUT /api/products/:id` - 更新商品
- `DELETE /api/products/:id` - 删除商品

### 订单接口
- `POST /api/orders` - 创建订单
- `GET /api/orders` - 获取订单列表
- `GET /api/orders/:id` - 获取订单详情
- `POST /api/orders/:id/pay` - 支付订单
- `POST /api/orders/:id/ship` - 发货
- `POST /api/orders/:id/confirm` - 确认收货
- `POST /api/orders/:id/cancel` - 取消订单
- `POST /api/orders/:id/refund` - 申请退款

### 管理员接口
- `POST /api/admin/shops/:id/approve` - 审核通过商家
- `POST /api/admin/shops/:id/reject` - 审核拒绝商家
- `GET /api/admin/statistics` - 平台统计
- `GET /api/admin/orders/export` - 导出订单报表

## 关键技术实现

### 并发库存控制
使用数据库事务 + 行级锁确保库存扣减的原子性：
```go
// 在订单创建时
result := tx.Model(&sku).Where("id = ? AND stock >= ?", skuID, quantity).Update("stock", gorm.Expr("stock - ?", quantity))
if result.RowsAffected == 0 {
    return errors.New("库存不足")
}
```

### 多商家订单拆分
购物车下单时按店铺分组，为每个店铺创建独立订单。

### 退款流程完整性
退款审核通过时自动恢复库存、更新订单状态、发送通知。

### 角色权限控制
通过中间件实现三级权限（买家/商家/管理员）API访问控制。

## 开发说明

### 数据库迁移
系统启动时会自动执行数据库迁移，创建所有表结构。

### 默认账号
系统初始化时会创建以下测试账号：
- 管理员: admin / admin123
- 测试买家: buyer / buyer123
- 测试商家: seller / seller123

## License

MIT
