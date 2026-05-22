# 在线赛事报名与成绩管理平台

一个基于 **Go (Gin + GORM + Redis) + Vue3 (TypeScript + Vite + Element Plus)** 的全栈赛事管理系统。

## 功能模块

| 模块 | 说明 |
|------|------|
| 赛事管理 | 创建、编辑、发布/取消发布、设置项目、分组、年龄组、名额、候补、费用 |
| 参赛者管理 | 注册、登录、实名认证、个人资料、历史参赛记录 |
| 报名系统 | 单人/团体报名、Redis 原子计数防超卖、满额候补、截止时间自动关闭 |
| 成绩录入 | 管理员按项目录入、批量 Excel 导入、自动排名与积分 |
| 成绩查询与证书 | 参赛者查询成绩排名、达到标准可下载 PDF 电子证书（含失败重试） |
| 数据统计 | 按赛事/项目/分组统计、Redis 缓存、导出 Excel 报表 |
| 消息通知 | 报名成功、候补确认、成绩发布等通过 Redis Pub/Sub 推送站内消息 |
| 日志审计 | 所有操作通过中间件记录日志 |

## 技术要点

- **并发防超卖**：`INCR` 原子操作 + 报名前先校验数据库计数，失败回滚 `DECR`
- **成绩重复录入校验**：查询 `user_id + event_item_id` 唯一
- **证书生成失败重试**：`pkg/retry` 指数退避，最多 3 次
- **JWT 鉴权**：所有 `/api/*` 接口统一 `Authorization: Bearer xxx`
- **参数校验**：`binding` 标签 + `ShouldBindJSON`
- **异常捕获**：`middleware.Recovery()` + 全局统一响应 `pkg/response`
- **配置文件管理**：`configs/config.yaml`（MySQL、Redis、JWT、上传路径）

## 目录结构

```
backend/
├── cmd/server/main.go          # 程序入口
├── configs/config.yaml         # 配置文件
├── internal/
│   ├── config/                 # 配置加载
│   ├── logger/                 # 日志
│   ├── database/               # MySQL & Redis 连接
│   ├── middleware/             # JWT 鉴权、Admin、恢复
│   ├── model/                  # GORM 数据模型
│   ├── dto/                    # 请求 DTO
│   ├── repository/             # 数据访问层
│   ├── service/                # 业务逻辑层
│   ├── handler/                # HTTP 接口层
│   ├── router/                 # 路由注册
│   ├── queue/                  # Redis Pub/Sub 消息队列
│   └── util/                   # 工具函数
└── pkg/
    ├── response/ jwt/ crypto/ excel/ pdf/ retry/

frontend/
├── src/
│   ├── api/                    # axios 封装与接口
│   ├── router/                 # vue-router
│   ├── stores/                 # pinia 状态
│   ├── views/                  # 页面（普通用户 + 管理员）
│   ├── components/             # 组件
│   ├── types/                  # TypeScript 类型
│   ├── styles/                 # 样式
│   ├── App.vue, main.ts
└── vite.config.ts, tsconfig.json
```

## 数据库表结构

所有表通过 `AutoMigrate` 自动创建，核心表：

- `users`：用户表（username 唯一，role 区分 admin/user/judge）
- `events`：赛事表（is_published 控制发布）
- `event_items`：赛事项目表（含性别、年龄组、名额、候补、费用）
- `registrations`：报名表（状态 pending/confirmed/waitlist/rejected/cancelled）
- `scores`：成绩表（自动 rank、points）
- `certificates`：证书表（状态 generating/generated/failed，带 retry_count）
- `messages`：站内消息表
- `operation_logs`：操作日志表

## 运行方式

### 1. 启动依赖

```bash
# 启动 MySQL 并创建数据库
mysql -uroot -p -e "CREATE DATABASE event_platform CHARACTER SET utf8mb4;"
# 启动 Redis
redis-server
```

### 2. 启动后端

```bash
cd backend
go mod tidy
go run ./cmd/server -config configs/config.yaml
```

后端监听 `:8080`。

### 3. 启动前端

```bash
cd frontend
npm install
npm run dev
```

前端监听 `:5173`，API 请求代理到后端。

## 默认管理员

后端未内置种子数据。首次启动后，可通过 MySQL 直接执行：

```sql
UPDATE users SET role='admin' WHERE username='your_user';
```

或通过 API 注册普通账号后再修改角色。

## 核心 API 一览

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/auth/register` | 注册 |
| POST | `/api/v1/auth/login` | 登录，返回 token |
| GET | `/api/v1/events` | 赛事列表（已发布） |
| GET | `/api/v1/events/:id` | 赛事详情 |
| POST | `/api/v1/registrations` | 提交报名 |
| GET | `/api/v1/registrations/me` | 我的报名 |
| GET | `/api/v1/scores/me` | 我的成绩 |
| POST | `/api/v1/certificates/:score_id/generate` | 生成证书 |
| GET | `/api/v1/certificates/me` | 我的证书 |
| GET | `/api/v1/messages` | 消息列表 |
| POST | `/api/v1/admin/events` | 创建赛事（admin） |
| PUT | `/api/v1/admin/events/:id/publish` | 发布赛事 |
| POST | `/api/v1/admin/scores` | 录入成绩 |
| POST | `/api/v1/admin/scores/import` | 批量导入成绩 Excel |
| GET | `/api/v1/admin/stats/export` | 导出统计报表 |
