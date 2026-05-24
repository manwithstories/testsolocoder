# 在线体育联赛与赛事管理平台

基于 **Go (Gin) + React (TypeScript) + PostgreSQL** 的全栈体育联赛管理系统。

## 功能特性

1. **联赛赛季管理** - 创建联赛、设置赛季时间、分组规则和积分制度，支持自定义比赛规则和奖项设置
2. **球队报名与阵容管理** - 球队在线报名、添加球员名单和位置信息，支持球员转会和阵容调整
3. **比赛排期与场地分配** - 系统自动生成赛程表或手动调整，每场比赛分配场地，支持场地冲突检测
4. **比分上报与积分榜** - 比赛结束后上报比分自动更新积分榜，支持加时赛、点球大战等特殊结果处理
5. **淘汰赛生成** - 支持单败淘汰、双败淘汰、小组赛+淘汰赛多种赛制，自动生成对阵表和晋级路径
6. **球员数据统计** - 记录进球、助攻、犯规等数据，生成球员排行榜
7. **裁判指派管理** - 管理员指派裁判执法比赛，裁判可确认或拒绝任务
8. **联赛费用管理** - 球队支付报名费和场地费，支持发票生成和费用报表
9. **消息通知** - 比赛变更、晋级通知等通过站内信推送
10. **数据导出** - 赛程表、积分榜、统计数据可导出为 Excel 和 PDF
11. **用户认证与权限控制** - JWT + RBAC，支持管理员/球队队长/球员/裁判四种角色

## 技术栈

### 后端
- **Go 1.21** + **Gin** 框架
- **PostgreSQL** + **GORM** ORM
- **JWT** 认证
- **excelize** Excel 导出

### 前端
- **React 18** + **TypeScript**
- **Vite** 构建工具
- **Ant Design 5** UI 组件库
- **React Router** 路由
- **Axios** HTTP 客户端

## 项目结构

```
.
├── backend/                  # Go 后端
│   ├── main.go              # 入口
│   ├── go.mod
│   ├── .env.example          # 环境变量样例
│   ├── config/              # 配置
│   ├── models/              # 数据库模型
│   ├── pkg/
│   │   ├── auth/            # JWT 认证
│   │   ├── database/        # 数据库连接
│   │   ├── logger/          # 日志
│   │   └── middleware/      # 中间件
│   ├── handlers/            # API 处理器
│   ├── routes/              # 路由配置
│   └── seeds/               # 种子数据
└── frontend/                # React 前端
    ├── src/
    │   ├── pages/           # 页面组件
    │   ├── api/             # API 调用
    │   ├── types/           # TypeScript 类型
    │   ├── App.tsx          # 主应用
    │   └── main.tsx         # 入口
    └── package.json
```

## 快速开始

### 前置条件
- Go 1.21+
- Node.js 18+
- PostgreSQL 12+

### 1. 配置数据库

```sql
CREATE DATABASE sports_league;
```

### 2. 启动后端

```bash
cd backend
cp .env.example .env
# 修改 .env 中的数据库连接信息
go mod tidy
go run main.go
```

后端将在 `http://localhost:8080` 启动，首次启动会自动迁移数据库并创建管理员账户。

默认管理员: `admin@example.com` / `admin123`

### 3. 启动前端

```bash
cd frontend
npm install
npm run dev
```

前端将在 `http://localhost:3000` 启动，通过 Vite 代理将 API 请求转发到后端。

## API 文档

### 认证
- `POST /api/auth/register` - 注册
- `POST /api/auth/login` - 登录
- `GET /api/auth/me` - 获取当前用户

### 联赛
- `GET /api/leagues` - 联赛列表
- `POST /api/leagues` - 创建联赛 (admin)
- `GET /api/leagues/:id` - 联赛详情
- `PUT /api/leagues/:id` - 更新联赛 (admin)
- `DELETE /api/leagues/:id` - 删除联赛 (admin)
- `POST /api/leagues/:id/seasons` - 创建赛季 (admin)
- `GET /api/leagues/:id/seasons/:season_id` - 赛季详情

### 球队
- `GET /api/teams` - 球队列表
- `POST /api/teams` - 创建球队 (admin)
- `GET /api/teams/:id` - 球队详情
- `PUT /api/teams/:id` - 更新球队 (admin)
- `DELETE /api/teams/:id` - 删除球队 (admin)
- `POST /api/teams/:id/players` - 添加球员 (admin)
- `PUT /api/teams/:id/players/:player_id` - 更新球员 (admin)
- `DELETE /api/teams/:id/players/:player_id` - 删除球员 (admin)
- `POST /api/teams/:id/register` - 球队报名 (admin/captain)

### 比赛
- `GET /api/matches` - 比赛列表
- `POST /api/seasons/:season_id/generate-schedule` - 生成赛程 (admin)
- `GET /api/matches/:id` - 比赛详情
- `PUT /api/matches/:id` - 更新比赛 (admin)
- `POST /api/matches/:id/report-score` - 上报比分 (admin)
- `GET /api/seasons/:season_id/standings` - 积分榜
- `POST /api/seasons/:season_id/generate-knockout` - 生成淘汰赛 (admin)
- `GET /api/venues` - 场地列表
- `POST /api/venues` - 添加场地 (admin)
- `GET /api/venues/check-conflict` - 场地冲突检测

### 裁判
- `GET /api/referees` - 裁判列表
- `POST /api/referees/assign` - 指派裁判 (admin)
- `GET /api/referees/assignments` - 裁判任务列表
- `PUT /api/referees/assignments/:id` - 接受/拒绝任务

### 费用
- `GET /api/fees` - 费用列表
- `POST /api/fees` - 创建费用 (admin)
- `PUT /api/fees/:id/paid` - 标记已付 (admin)
- `GET /api/fees/seasons/:season_id/report` - 费用报表

### 消息通知
- `GET /api/notifications` - 通知列表
- `PUT /api/notifications/:id/read` - 标记已读
- `PUT /api/notifications/read-all` - 全部标记已读
- `GET /api/notifications/unread-count` - 未读数量

### 数据统计
- `GET /api/stats/rankings` - 球员排行榜
- `POST /api/stats` - 录入数据 (admin)
- `GET /api/stats/players/:id` - 球员统计详情

### 数据导出
- `GET /api/export/schedule/:season_id` - 导出赛程 Excel
- `GET /api/export/standings/:season_id` - 导出积分榜 Excel
- `GET /api/export/stats` - 导出球员数据 Excel
- `GET /api/export/pdf/:season_id` - 导出 PDF

## 角色权限

| 功能 | admin | captain | player | referee |
|------|:-----:|:-------:|:------:|:-------:|
| 创建/编辑联赛 | ✓ | ✗ | ✗ | ✗ |
| 创建/编辑球队 | ✓ | ✓ | ✗ | ✗ |
| 生成赛程 | ✓ | ✗ | ✗ | ✗ |
| 上报比分 | ✓ | ✗ | ✗ | ✗ |
| 指派裁判 | ✓ | ✗ | ✗ | ✗ |
| 管理费用 | ✓ | ✗ | ✗ | ✗ |
| 查看数据 | ✓ | ✓ | ✓ | ✓ |
| 球队报名 | ✓ | ✓ | ✗ | ✗ |
| 接受/拒绝任务 | - | - | - | ✓ |

## 数据库模型

- **User** - 用户 (email, password, role, team_id)
- **League** - 联赛 (name, sport, status)
- **Season** - 赛季 (league_id, dates, format, points rules)
- **Team** - 球队 (name, captain_id)
- **Player** - 球员 (team_id, number, position)
- **Registration** - 报名 (season_id, team_id, group_name, status)
- **Match** - 比赛 (season_id, teams, venue, scores, status)
- **Venue** - 场地 (name, address, capacity)
- **RefereeAssignment** - 裁判指派 (match_id, referee_id, status)
- **Fee** - 费用 (season_id, team_id, type, amount, invoice_no)
- **Notification** - 通知 (user_id, title, content, type)
- **PlayerStat** - 球员统计 (player_id, match_id, goals, assists, etc.)

## License

MIT
