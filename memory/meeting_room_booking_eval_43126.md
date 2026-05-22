---
name: meeting-room-booking-eval-43126
description: Go+Gin+Vue3+TypeScript+Redis在线会议室预订与会议管理Web应用R1评测
metadata:
  type: project
---

会议室预订管理系统R1评测：FAIL（后端编译失败+缺少go.sum+周期预订父关联逻辑错误+空间管理员权限未实现+会议室可用时段硬编码）

**Why:** 评测发现5个严重/中等问题导致系统无法正常运行
**How to apply:** 下一轮需修复编译错误、生成go.sum、修复周期预订关联逻辑、启用空间管理员权限、实现会议室可用时段配置