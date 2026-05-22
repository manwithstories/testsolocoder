---
name: housekeeping-platform-eval-43115
description: Go+Gin+React+TS家政服务预约与管理平台R1评测：前端存在大量TypeScript编译错误导致无法构建
metadata:
  type: project
---

## 项目信息
- **项目名称**：家政服务预约与管理平台
- **技术栈**：Go+Gin+React+TypeScript+PostgreSQL
- **评测轮次**：R1

## R1评测结果：FAIL

### 核心问题摘要

**前端编译错误（高严重度）**：
1. 路由文件重复导入WithdrawList，引用不存在的WithdrawReview组件
2. Dashboard页面API方法名与实际服务不匹配（getOverview/getChartData/exportExcel不存在）
3. BillList页面缺少Form导入，调用不存在的API方法，字段名错误
4. MessageCenter页面API方法名不匹配
5. Certification页面引用不存在的类型和API方法
6. ComplaintHandle页面缺少Select导入，字段名错误
7. MyServices页面API返回值处理错误
8. UserManage页面工具函数名错误

**后端编译错误（高严重度）**：
1. go.mod中jwt版本v5.1.1不存在
2. order.go/review.go/bill.go缺少gorm.io/gorm导入

### 下轮修复方向
需要系统性检查前后端代码编译问题，特别是：
- API方法名与服务定义的一致性
- 组件导入的完整性
- 类型字段名称的正确性
- go.mod依赖版本的有效性

[[housekeeping-platform-eval-43115]]
