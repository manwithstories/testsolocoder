# API Tester - API性能压测与诊断工具

一个功能强大的API性能压测与诊断命令行工具，基于Go Cobra框架开发。

## 功能特性

### 1. 测试场景管理
- 支持多种HTTP请求类型（GET, POST, PUT, DELETE, PATCH等）
- 自定义请求头配置
- 链式接口调用和依赖参数传递
- 支持从响应中提取变量用于后续请求
- 场景权重配置，支持多场景混合压测

### 2. 并发控制
- 设置并发数（Workers）
- 设置总请求数
- 设置持续测试时间
- **阶梯式加压模式**（Steps Mode）
- 支持按权重分配不同场景的执行比例

### 3. 结果分析
- 实时显示QPS、响应时间、错误率等指标
- 响应时间分布统计（P50, P90, P95, P99）
- 状态码分布统计
- 详细的错误分类统计
- 测试完成后生成详细的文本报告

### 4. 断言校验
- 响应状态码校验
- 响应体JSON字段校验
- 响应时间阈值校验
- 响应头校验
- 响应体内容包含校验

### 5. 配置管理
- YAML格式配置文件
- 支持配置继承（inherits）
- 支持环境变量覆盖
- 命令行参数覆盖配置

### 6. 日志记录
- 详细记录每个请求的发送时间、响应时间
- 错误信息和断言失败日志
- 支持日志级别控制（debug, info, warn, error）
- 支持日志输出到文件

### 7. 异常处理
- 网络超时自动重试
- 连接失败、服务不可用等错误捕获
- 可配置的重试次数
- Panic恢复机制，防止程序崩溃
- 优雅的信号处理（Ctrl+C停止）

## 目录结构

```
apitester/
├── cmd/
│   └── apitester/           # 命令行入口
│       ├── main.go          # 主入口
│       ├── run.go           # run命令
│       └── config.go        # validate和init命令
├── internal/
│   ├── config/              # 配置管理模块
│   │   └── config.go
│   ├── logger/              # 日志管理模块
│   │   └── logger.go
│   ├── scenario/            # 场景执行模块
│   │   └── scenario.go
│   ├── assertion/           # 断言校验模块
│   │   └── assertion.go
│   ├── executor/            # 并发执行模块
│   │   └── executor.go
│   ├── httpclient/          # HTTP客户端模块
│   │   └── client.go
│   └── stats/               # 结果统计模块
│       └── stats.go
├── pkg/
│   └── models/              # 数据模型
│       └── models.go
├── examples/                # 示例配置文件
│   ├── example.yaml
│   ├── stress-test.yaml
│   ├── base-config.yaml
│   └── prod-config.yaml
├── go.mod
├── go.sum
└── README.md
```

## 快速开始

### 编译

```bash
go build -o apitester ./cmd/apitester
```

### 查看帮助

```bash
./apitester --help
```

### 生成示例配置

```bash
./apitester init
```

### 验证配置文件

```bash
./apitester validate -c examples/example.yaml
```

### 运行压测

```bash
# 基本运行
./apitester run -c examples/example.yaml

# 指定并发数和持续时间
./apitester run -c examples/example.yaml -w 50 -d 120

# 输出详细日志
./apitester run -c examples/example.yaml -l debug

# 输出日志到文件
./apitester run -c examples/example.yaml --log-file test.log

# 生成报告文件
./apitester run -c examples/example.yaml -r report.txt
```

## 配置文件详解

### 基本结构

```yaml
name: "测试名称"
description: "测试描述"
base_url: "https://api.example.com"
timeout: 30           # 请求超时时间（秒）
retries: 2            # 重试次数

variables:            # 全局变量
  env: "production"
  token: "${API_TOKEN}"  # 从环境变量读取

headers:              # 全局请求头
  Content-Type: "application/json"
  Authorization: "Bearer {{token}}"

concurrency:          # 并发配置
  mode: "duration"    # workers | duration | steps
  workers: 50         # 并发数
  duration: 60        # 持续时间（秒）
  total_requests: 10000  # 总请求数（workers模式）
  steps:              # 阶梯模式配置
    - workers: 10
      duration: 30
    - workers: 30
      duration: 30
    - workers: 50
      duration: 30

scenarios:            # 测试场景
  - name: "场景名称"
    weight: 3         # 权重
    variables: {}     # 场景变量
    headers: {}       # 场景请求头
    requests:         # 请求列表
      - name: "请求名称"
        method: "GET"
        path: "/api/users"
        headers: {}   # 请求级请求头
        body: '{}'    # 请求体
        think_time: 100  # 请求后等待时间（毫秒）
        
        # 断言配置
        assertions:
          - type: "status_code"
            operator: "eq"
            value: 200
          - type: "response_time"
            operator: "lte"
            value: 500
          - type: "body_field"
            field: "code"
            operator: "eq"
            value: 0
          - type: "body_contains"
            value: "success"
          - type: "header"
            field: "Content-Type"
            operator: "contains"
            value: "json"
        
        # 变量提取
        extract:
          - name: "user_id"
            from: "body"    # body | header | cookie
            path: "data.id"
```

### 并发模式说明

1. **workers模式**：固定并发数，执行指定总请求数
2. **duration模式**：固定并发数，执行指定时间
3. **steps模式**：阶梯式增加并发数

### 断言类型

| 类型 | 说明 | 支持的操作符 |
|------|------|-------------|
| status_code | 状态码校验 | eq, ne, gt, lt, gte, lte |
| response_time | 响应时间校验（毫秒） | eq, ne, gt, lt, gte, lte |
| body_field | JSON字段校验 | eq, ne, gt, lt, gte, lte, contains, not_contains, starts_with, ends_with |
| body_contains | 响应体包含 | - |
| header | 响应头校验 | eq, ne, contains, not_contains, starts_with, ends_with |

### 变量提取

支持从以下位置提取变量：
- `body`：从JSON响应体中提取，如 `data.user.id`
- `header`：从响应头中提取
- `cookie`：从Set-Cookie中提取

提取的变量可以在后续请求中使用 `{{变量名}}` 引用。

### 配置继承

```yaml
# base-config.yaml
name: "基础配置"
timeout: 10
retries: 2
headers:
  Content-Type: "application/json"

# prod-config.yaml
inherits: base-config.yaml
name: "生产配置"
timeout: 30  # 覆盖基础配置
```

### 环境变量覆盖

配置文件中可以使用 `${ENV_VAR}` 语法引用环境变量：

```yaml
variables:
  api_token: "${API_TOKEN}"
  password: "${DB_PASSWORD}"
```

也可以通过环境变量 `APITESTER_BASE_URL` 覆盖base_url。

## 使用示例

### 1. 简单的GET请求测试

```yaml
name: "简单测试"
base_url: "https://jsonplaceholder.typicode.com"
concurrency:
  mode: "duration"
  workers: 10
  duration: 60
scenarios:
  - name: "获取用户"
    requests:
      - name: "用户列表"
        method: "GET"
        path: "/users"
        assertions:
          - type: "status_code"
            operator: "eq"
            value: 200
```

### 2. 链式调用和参数传递

```yaml
name: "链式调用"
base_url: "https://api.example.com"
scenarios:
  - name: "用户流程"
    requests:
      - name: "登录"
        method: "POST"
        path: "/auth/login"
        body: '{"username": "user", "password": "pass"}'
        assertions:
          - type: "status_code"
            operator: "eq"
            value: 200
        extract:
          - name: "token"
            from: "body"
            path: "data.token"

      - name: "获取用户信息"
        method: "GET"
        path: "/user/profile"
        headers:
          Authorization: "Bearer {{token}}"
        assertions:
          - type: "status_code"
            operator: "eq"
            value: 200
```

### 3. 阶梯式压力测试

```yaml
name: "压力测试"
base_url: "https://api.example.com"
concurrency:
  mode: "steps"
  steps:
    - workers: 10    # 10并发
      duration: 60   # 持续60秒
    - workers: 30    # 30并发
      duration: 60
    - workers: 50    # 50并发
      duration: 60
    - workers: 100   # 100并发
      duration: 120
```

## 命令行参数

### run 命令

```
-c, --config string     配置文件路径
-l, --log-level string  日志级别 (debug, info, warn, error) (default "info")
    --log-file string   日志文件路径
-w, --workers int       并发数 (覆盖配置文件)
-d, --duration int      测试持续时间(秒) (覆盖配置文件)
-n, --requests int      总请求数 (覆盖配置文件)
-r, --report string     输出报告文件路径
```

### validate 命令

```
-c, --config string     配置文件路径
```

## 输出报告示例

```
================================================================================
  TEST REPORT: API性能测试示例
================================================================================

  Start Time:    2024-01-01 12:00:00
  End Time:      2024-01-01 12:01:00
  Duration:      01:00

--------------------------------------------------------------------------------
  REQUEST SUMMARY
--------------------------------------------------------------------------------
  Total Requests:   10000
  Successful:       9850
  Failed:           150
  Error Rate:       1.50%
  QPS:              166.67

--------------------------------------------------------------------------------
  RESPONSE TIME DISTRIBUTION
--------------------------------------------------------------------------------
  Average:  45ms            Min:  12ms            Max:  2340ms
  P50:      38ms            P90:  89ms            P95:  156ms           P99:  567ms

--------------------------------------------------------------------------------
  STATUS CODE DISTRIBUTION
--------------------------------------------------------------------------------
  200: 9850 (98.5%)
  500: 120 (1.2%)
  404: 30 (0.3%)
```

## 许可证

MIT License
