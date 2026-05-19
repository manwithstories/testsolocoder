# 个人财务记账工具 (Finance Tracker)

一个基于Go和Cobra开发的命令行个人财务管理工具，帮助你追踪日常收支、管理预算、生成财务报告。

## 功能特性

### 1. 账户管理
- 支持创建多个账户（现金、银行卡、支付宝、微信、信用卡等）
- 每个账户可设置初始余额和币种
- 支持账户间转账，自动检查余额是否充足
- 完整的账户余额变动记录

### 2. 交易记录
- 支持添加收入和支出记录
- 每笔交易包含：金额、分类、日期、备注
- 支持按时间范围、分类、账户筛选查询
- 支持批量导入CSV格式的交易记录
- 交易记录写入前自动验证日期格式和分类是否存在

### 3. 分类管理
- 支持自定义收入和支出分类
- 支持多级分类（如：餐饮 -> 早餐、午餐、晚餐）
- 树形结构展示分类层级关系

### 4. 预算管理
- 为每个分类设置月度预算
- 添加交易时自动检查是否超支并发出警告
- 预算执行情况实时统计

### 5. 统计分析
- 按时间范围生成收支汇总报告
- 按分类统计支出/收入占比（可视化图表）
- 按账户统计资金流动情况
- 支持导出TXT和CSV格式报告

### 6. 数据备份
- 支持手动创建数据库备份
- 支持自动定时备份（小时/日/周）
- 自动清理过期备份，保留最近N个备份文件
- 支持从备份恢复数据库

### 7. 数据安全
- 所有金额操作进行数据校验，防止负数或非法输入
- 每个操作都记录日志到文件，方便排查问题
- 配置文件使用YAML格式，支持个性化设置

## 安装

### 从源码编译

```bash
# 克隆仓库
git clone <repository-url>
cd finance-tracker

# 编译
go build -o finance

# 移动到PATH目录
sudo mv finance /usr/local/bin/
```

## 快速开始

### 1. 创建账户

```bash
# 创建现金账户，初始余额1000元
finance account create --name "现金" --type cash --balance 1000

# 创建银行卡账户
finance account create --name "工商银行" --type bank --balance 5000

# 创建支付宝账户
finance account create --name "支付宝" --type alipay --balance 2000

# 查看所有账户
finance account list
```

### 2. 创建分类

```bash
# 创建支出分类
finance category create --name "餐饮" --type expense
finance category create --name "交通" --type expense
finance category create --name "购物" --type expense
finance category create --name "房租" --type expense

# 创建子分类
finance category create --name "早餐" --type expense --parent 1
finance category create --name "午餐" --type expense --parent 1
finance category create --name "晚餐" --type expense --parent 1

# 创建收入分类
finance category create --name "工资" --type income
finance category create --name "奖金" --type income

# 以树形结构查看分类
finance category list --tree
```

### 3. 记录交易

```bash
# 添加一笔支出（早餐）
finance transaction add --type expense --amount 25.50 --category 8 --account 1 --desc "早餐"

# 添加一笔收入（工资）
finance transaction add --type income --amount 15000.00 --category 2 --account 2 --desc "5月工资"

# 查看交易记录
finance transaction list

# 按时间范围筛选
finance transaction list --start 2026-05-01 --end 2026-05-31
```

### 4. 转账

```bash
# 从银行卡转账500元到现金
finance account transfer --from 2 --to 1 --amount 500 --desc "取钱"
```

### 5. 设置预算

```bash
# 为餐饮分类设置5月预算1500元
finance budget set --category 1 --amount 1500 --month 05 --year 2026

# 查看预算执行情况
finance budget list --month 05 --year 2026

# 检查超支情况
finance budget check --month 05 --year 2026
```

### 6. 统计分析

```bash
# 查看收支汇总（默认最近30天）
finance stats summary

# 指定时间范围
finance stats summary --start 2026-05-01 --end 2026-05-31

# 按分类查看支出占比
finance stats category --type expense --start 2026-05-01 --end 2026-05-31

# 生成详细报告
finance stats report --start 2026-05-01 --end 2026-05-31 --output report.txt

# 导出CSV格式
finance stats report --start 2026-05-01 --end 2026-05-31 --output report.csv
```

### 7. 数据备份

```bash
# 创建手动备份
finance backup create

# 查看备份列表
finance backup list

# 从备份恢复
finance backup restore --path /path/to/backup.db
```

## 配置文件

配置文件位于 `~/.finance-tracker/config.yaml`，支持以下配置：

```yaml
default_account: cash
default_currency: CNY
budget_alert: true

backup:
  enabled: false           # 是否启用自动备份
  schedule: daily          # 备份频率: hourly, daily, weekly
  local_path: ~/.finance-tracker/backups
  retention_days: 30       # 备份保留天数
  max_backups: 10          # 最大保留备份数

database_path: ~/.finance-tracker/finance.db
log_path: ~/.finance-tracker/logs
```

## CSV导入格式

批量导入交易记录的CSV文件格式：

```csv
date,type,amount,category,account,description
2026-05-01,expense,25.50,餐饮,现金,早餐
2026-05-01,expense,45.00,餐饮,银行卡,午餐
2026-05-02,income,15000.00,工资,银行卡,5月工资
```

字段说明：
- `date`: 交易日期，格式 YYYY-MM-DD
- `type`: 交易类型，`income` 或 `expense`
- `amount`: 交易金额，正数
- `category`: 分类名称（必须已存在）
- `account`: 账户名称（必须已存在）
- `description`: 备注信息（可选）

导入命令：
```bash
finance transaction import --file transactions.csv
```

## 数据存储

- 数据库：SQLite
- 默认位置：`~/.finance-tracker/finance.db`
- 日志位置：`~/.finance-tracker/logs/`
- 备份位置：`~/.finance-tracker/backups/`

## 项目结构

```
finance-tracker/
├── cmd/                    # 命令行入口
│   ├── account/           # 账户管理命令
│   ├── category/          # 分类管理命令
│   ├── transaction/       # 交易记录命令
│   ├── budget/            # 预算管理命令
│   ├── stats/             # 统计分析命令
│   ├── backup/            # 备份管理命令
│   └── root.go            # 根命令
├── internal/
│   ├── config/            # 配置管理
│   ├── models/            # 数据模型
│   ├── storage/           # 数据库操作
│   ├── logger/            # 日志管理
│   ├── account/           # 账户业务逻辑
│   ├── category/          # 分类业务逻辑
│   ├── transaction/       # 交易业务逻辑
│   ├── budget/            # 预算业务逻辑
│   ├── stats/             # 统计分析逻辑
│   └── backup/            # 备份业务逻辑
├── configs/               # 配置和示例文件
├── main.go                # 程序入口
├── go.mod
└── go.sum
```

## 安全特性

1. **数据校验**：所有金额操作都会检查是否为正数，防止负数输入
2. **余额检查**：转账和支出操作会检查账户余额是否充足
3. **分类验证**：添加交易时验证分类是否存在且类型匹配
4. **日期验证**：自动验证日期格式，支持多种格式解析
5. **操作日志**：所有操作都记录日志，包含时间、操作类型和详细信息
6. **事务保护**：涉及多个账户的操作使用数据库事务，确保数据一致性

## 命令速查

| 命令 | 说明 |
|------|------|
| `finance account create` | 创建账户 |
| `finance account list` | 列出账户 |
| `finance account transfer` | 账户转账 |
| `finance category create` | 创建分类 |
| `finance category list --tree` | 树形查看分类 |
| `finance transaction add` | 添加交易 |
| `finance transaction list` | 列出交易 |
| `finance transaction import` | 导入CSV |
| `finance budget set` | 设置预算 |
| `finance budget list` | 查看预算 |
| `finance budget check` | 检查超支 |
| `finance stats summary` | 收支汇总 |
| `finance stats category` | 分类统计 |
| `finance stats report` | 生成报告 |
| `finance backup create` | 创建备份 |
| `finance backup list` | 备份列表 |

## License

MIT License
