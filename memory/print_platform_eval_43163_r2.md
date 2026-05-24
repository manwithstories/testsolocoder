---
name: print-platform-eval-43163-r2
description: Go+Gin+React+TS在线印刷定制与订单管理平台R2评测：FAIL(R1全部12问题已修复，新发现信用额度未校验+拆单数据一致性问题)
metadata:
  type: project
---

## 评测结果：FAIL

### R1问题修复情况（全部已修复）

1. **项目结构完整** - 后端8个Go文件（main.go, handlers, models, service, middleware, auth, config, database），前端12个TSX页面组件
2. **后端编译成功** - go.sum存在，所有依赖正确
3. **前端构建成功** - TypeScript编译通过，Vite构建完成
4. **核心功能全部实现**：商品模板管理、在线定制下单、订单流程管理、智能报价引擎、生产排程、客户财务管理、数据统计看板、用户认证权限

### R2新发现问题（2个中等问题）

**【问题1】信用额度控制未实现**
- Customer模型有credit_limit字段
- CreateOrder获取客户信息但未校验订单金额是否超过信用额度
- **影响**：客户可无限下单，信用额度功能形同虚设

**【问题2】拆单数据一致性问题**
- SplitOrder创建新订单复制商品项，但未从原订单删除对应商品项
- 仅更新原订单金额，商品明细仍保留全部商品
- **影响**：拆单后原订单商品总数与金额不匹配，数据错乱

### 技术栈验证

- Go 1.22 + Gin + GORM + SQLite ✓
- React 18 + TypeScript + Vite + Axios ✓
- JWT认证 + 操作日志审计 ✓