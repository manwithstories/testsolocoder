---
name: drone-rental-eval-43171-r4
description: Go+Gin+Vue3+TS无人机租赁与航拍服务预约管理平台R4评测
metadata:
  type: project
---

Go+Gin+Vue3+TS无人机租赁与航拍服务预约管理平台R4评测：R3三个问题已全部修复(批量导入+可用时段设置+日期搜索过滤)，但新发现1个严重问题：数据统计导出功能因未携带认证token致401失败，用户完全无法导出Excel报表。

**Why：** 前端使用window.open()下载导出文件但未在URL中包含JWT token，后端中间件虽支持查询参数传token但前端未使用。

**How to apply：** 下一轮需修复导出功能，在前端导出链接中添加token查询参数。
