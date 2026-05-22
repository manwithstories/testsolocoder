---
name: ticket-system-eval-43073
description: Go+Gin+Vue3在线活动报名与票务管理Web应用R1评测
metadata:
  type: project
---

Go+Gin+Vue3在线活动报名与票务管理Web应用第一轮评测：FAIL

## 核心问题
签到统计接口无法正确获取活动ID参数，导致按活动筛选签到统计功能失效。

## 功能验证结果
- 活动管理：PASS（创建/编辑/发布/取消/海报上传）
- 票务管理：PASS（多票型/独立库存价格/售罄自动下架/Redis同步）
- 报名购票：PASS（批量购票/优惠券抵扣/事务一致性/Redis防超卖）
- 签到管理：PASS（二维码生成/扫码签到/时间记录）
- 统计报表：PASS（按活动/票型/时间段统计/ECharts图表/Excel导出）
- 用户管理：PASS（管理员+普通用户/JWT认证/权限控制）

## 技术要求
- Redis原子扣减：PASS
- 数据库多表关联：PASS（7表关联）
- 参数校验：PASS
- 异常处理：PASS
- 环境配置：PASS
- 日志记录：PASS

## 下一轮修复
签到统计接口参数获取方式错误，需从查询参数正确获取activityId。
