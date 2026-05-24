---
name: business-registration-platform-eval-43166
description: Go+Gin+Vue3+TS企业工商注册代办平台R1评测：FAIL(后端jwt库版本不存在+缺少go.sum+前端vue-tsc版本不兼容+TypeScript类型错误+密码加密函数未实际加密)
metadata:
  type: project
---

Go+Gin+Vue3+TypeScript企业工商注册代办与公司设立服务平台R1评测：FAIL

**严重问题：**
1. 后端go.mod指定jwt库版本v5.1.1不存在，正确版本应为v5.2.0+，且缺少go.sum文件
2. 前端vue-tsc@^1.8.25与Node.js v25不兼容，构建失败
3. 前端存在多处TypeScript类型错误（UserRole类型使用错误、API响应类型错误、axios拦截器类型不兼容）
4. agent_service.go中HashPassword函数直接返回原始密码，未实际加密，存在安全隐患
