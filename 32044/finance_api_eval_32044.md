---
name: finance-api-eval-32044
description: 个人财务管理API评测：Go+Gin，R1 FAIL(账户更新余额损坏+交易更新类型分类不匹配)
metadata:
  type: project
---

个人财务管理API服务评测：Go + Gin + GORM + MySQL + JWT，R1 FAIL。

**R1问题：**
- 【严重】账户更新接口InitialBalance float64零值导致余额数据损坏：只更新备注/名称时initial_balance被误重置为0
- 【中等】交易更新接口修改type但不提供category_id时跳过分类类型校验，导致类型/分类不一致

**需求覆盖：** 账户管理、交易记录、分类管理、月度统计、预算功能、数据导出CSV、JWT认证、用户数据隔离 — 全部实现。
