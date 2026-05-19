# Fit Tracker - 命令行健身训练记录工具

一个基于 Python + Click 的命令行健身训练记录工具，专为有健身习惯的用户设计。

## 功能特性

### 1. 动作库管理
- ✅ 添加、删除、更新动作
- ✅ 支持分类：力量、有氧、柔韧性、其他
- ✅ 每个动作可设置默认组数、次数、时长
- ✅ 智能删除：有关联记录的动作标记为归档而非物理删除

### 2. 训练记录
- ✅ 开始/结束训练会话
- ✅ 逐条记录动作的组数、次数、重量、时长
- ✅ 自动计算总训练容量
- ✅ 训练中途可随时取消，丢弃本次记录
- ✅ 查看训练历史

### 3. 训练计划模板
- ✅ 创建可复用的训练计划
- ✅ 向计划中添加/移除动作
- ✅ 安排计划到每周具体日期（0=周一, 6=周日）
- ✅ 计划中的动作被删除时自动提示已归档

### 4. 统计进度
- ✅ 总体统计概览（总次数、总容量、连续天数）
- ✅ 每周训练频率和容量趋势（带可视化柱状图）
- ✅ 每月训练频率和容量趋势（带可视化柱状图）
- ✅ 每个动作的个人最佳记录（最大重量、最大容量）
- ✅ 连续训练天数统计

### 5. 数据导出
- ✅ 导出训练历史为 CSV
- ✅ 导出动作库为 CSV
- ✅ 导出训练计划为 CSV
- ✅ 一键导出所有数据

## 安装

```bash
pip install -e .
```

## 使用方法

### 查看帮助
```bash
fit --help
fit exercise --help
fit train --help
fit plan --help
fit stats --help
fit export --help
```

### 动作库管理
```bash
# 查看可用分类
fit exercise categories

# 添加动作
fit exercise add -n "卧推" -c 力量 -s 4 -r 10
fit exercise add -n "跑步" -c 有氧 -d 1800

# 列出所有动作
fit exercise list
fit exercise list -c 力量  # 按分类筛选
fit exercise list -a       # 包含已归档的动作

# 更新动作
fit exercise update 1 --name "平板卧推" --sets 5

# 删除动作（有关联记录则归档）
fit exercise delete 1
```

### 训练记录
```bash
# 开始训练
fit train start -n "胸部训练日"

# 记录动作组
fit train add-set -e 卧推 -r 10 -w 60
fit train add-set -e 深蹲 -r 8 -w 80

# 查看当前训练状态
fit train status

# 删除某一组
fit train remove-set 1

# 结束训练
fit train finish

# 取消训练（丢弃所有记录）
fit train cancel

# 查看训练历史
fit train history
fit train history -l 20
```

### 训练计划模板
```bash
# 创建计划
fit plan create -n "推拉腿分化" -d "经典力量训练分化"

# 向计划添加动作
fit plan add-exercise -p 1 -e 卧推 -s 4 -r 10
fit plan add-exercise -p 1 -e 深蹲

# 查看计划详情
fit plan show 1

# 从计划移除动作
fit plan remove-exercise -p 1 -e 1

# 安排计划到周几（0=周一, 6=周日）
fit plan schedule -p 1 -d 0  # 安排到周一
fit plan schedule -p 1 -d 2  # 安排到周三

# 查看本周计划
fit plan week

# 取消安排
fit plan unschedule 1

# 列出所有计划
fit plan list

# 删除计划
fit plan delete 1
```

### 统计进度
```bash
# 总体概览
fit stats overview

# 每周统计（带可视化图表）
fit stats weekly
fit stats weekly -w 8

# 每月统计（带可视化图表）
fit stats monthly
fit stats monthly -m 12

# 个人最佳记录
fit stats pr              # 显示所有动作PR
fit stats pr -e 卧推      # 查看特定动作PR

# 连续训练天数
fit stats streak
```

### 数据导出
```bash
# 导出训练历史
fit export training -o training.csv

# 导出动作库
fit export exercises -o exercises.csv

# 导出训练计划
fit export plans -o plans.csv

# 一键导出所有数据
fit export all -o ./export_dir
```

## 数据存储

数据存储在 SQLite 数据库中，位置：
- macOS/Linux: `~/.fit_tracker/fit_tracker.db`

## 技术栈

- **Click** - 命令行界面框架
- **SQLAlchemy** - ORM 数据库操作
- **Rich** - 终端格式化输出
- **SQLite** - 本地数据库
