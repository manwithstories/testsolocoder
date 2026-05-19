# Book Library - 个人图书收藏与阅读管理系统

## 项目简介

一个功能完整的个人图书收藏与阅读管理Web应用，使用Go + Gin后端和Vue3 + TypeScript前端构建。

## 核心功能

### 1. 图书管理
- 图书的增删改查
- ISBN自动获取书籍信息（调用公开接口）
- 支持封面图片上传
- 阅读状态管理（想读/在读/已读/放弃）

### 2. 阅读进度追踪
- 记录当前页码，自动计算阅读百分比
- 支持多本书同时标记为在读
- 读书笔记功能

### 3. 分类标签系统
- 自定义标签（支持颜色标记）
- 多级分类管理
- 按标签和分类筛选图书

### 4. 借阅管理
- 记录借阅人信息和联系方式
- 借出/归还日期管理
- 逾期提醒

### 5. 统计分析
- 年度阅读数量趋势图
- 阅读热力图（GitHub风格）
- 阅读时长分布
- 分类和标签统计

### 6. 阅读目标
- 年度读书计划设置
- 每月阅读目标
- 完成进度追踪

## 技术栈

### 后端
- Go 1.21+
- Gin 1.9 - Web框架
- GORM 1.25 - ORM
- SQLite - 数据库
- Logrus - 日志
- Godotenv - 配置管理

### 前端
- Vue 3.4
- TypeScript
- Pinia - 状态管理
- Vue Router 4 - 路由
- Element Plus - UI组件库
- ECharts - 图表
- Axios - HTTP客户端
- Vite - 构建工具

## 快速开始

### 后端启动

```bash
cd backend
go mod download
go run cmd/server/main.go
```

后端服务将在 http://localhost:8080 启动

### 前端启动

```bash
cd frontend
npm install
npm run dev
```

前端开发服务器将在 http://localhost:5173 启动

## 项目结构

```
.
├── backend/                    # Go后端
│   ├── cmd/
│   │   └── server/
│   │       └── main.go        # 应用入口
│   ├── internal/
│   │   ├── config/            # 配置管理
│   │   ├── logger/            # 日志
│   │   ├── database/          # 数据库连接
│   │   ├── models/            # 数据模型
│   │   ├── handlers/          # API处理器
│   │   ├── middleware/        # 中间件
│   │   ├── errors/            # 错误处理
│   │   └── utils/             # 工具函数
│   ├── .env                   # 环境变量
│   └── go.mod
└── frontend/                   # Vue前端
    ├── src/
    │   ├── api/               # API服务
    │   ├── router/            # 路由配置
    │   ├── store/             # 状态管理
    │   ├── views/             # 页面组件
    │   ├── components/        # 公共组件
    │   ├── types/             # TypeScript类型
    │   ├── styles/            # 全局样式
    │   ├── App.vue
    │   └── main.ts
    ├── package.json
    └── vite.config.ts
```

## API接口

### 图书管理
- `GET /api/books` - 获取图书列表
- `GET /api/books/:id` - 获取图书详情
- `POST /api/books` - 创建图书
- `PUT /api/books/:id` - 更新图书
- `DELETE /api/books/:id` - 删除图书
- `POST /api/books/:id/cover` - 上传封面
- `GET /api/books/isbn/:isbn` - 通过ISBN获取书籍信息
- `PATCH /api/books/:id/progress` - 更新阅读进度
- `PATCH /api/books/:id/status` - 更新阅读状态

### 标签管理
- `GET /api/tags` - 获取标签列表
- `POST /api/tags` - 创建标签
- `PUT /api/tags/:id` - 更新标签
- `DELETE /api/tags/:id` - 删除标签

### 分类管理
- `GET /api/categories` - 获取分类列表
- `POST /api/categories` - 创建分类
- `PUT /api/categories/:id` - 更新分类
- `DELETE /api/categories/:id` - 删除分类

### 借阅管理
- `GET /api/borrows` - 获取借阅记录
- `POST /api/borrows` - 创建借阅记录
- `POST /api/borrows/:id/return` - 归还图书
- `DELETE /api/borrows/:id` - 删除借阅记录

### 统计分析
- `GET /api/stats/overview` - 概览统计
- `GET /api/stats/yearly-trend` - 年度趋势
- `GET /api/stats/heatmap` - 阅读热力图
- `GET /api/stats/duration` - 阅读时长分布
- `GET /api/stats/categories` - 分类统计
- `GET /api/stats/tags` - 标签统计

### 阅读目标
- `GET /api/goals` - 获取目标列表
- `GET /api/goals/yearly-progress` - 年度目标进度
- `POST /api/goals` - 创建目标
- `PUT /api/goals/:id` - 更新目标
- `DELETE /api/goals/:id` - 删除目标

## 配置说明

后端配置通过 `.env` 文件进行管理：

```env
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
GIN_MODE=debug

DB_PATH=data/booklibrary.db

UPLOAD_MAX_SIZE=5
UPLOAD_PATH=uploads
UPLOAD_ACCESS_URL=/uploads

ISBN_TIMEOUT=10

LOG_LEVEL=debug
LOG_FORMAT=text
LOG_OUTPUT=stdout
```

## 边界情况处理

- **ISBN格式校验**：支持ISBN-10和ISBN-13格式自动检测和校验
- **封面图片上传**：限制5MB以内，支持jpg/png/gif/webp格式
- **并发借阅冲突**：同一本书只能有一个有效的借阅记录
- **页码有效性检查**：当前页码不能超过总页数
- **完善的错误处理**：统一的错误响应格式，详细的错误信息

## License

MIT
