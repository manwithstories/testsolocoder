---
name: 43069-eval-round1
description: Go Cobra API压测工具第1轮评测发现ramp_up未实现、配置继承合并缺陷
metadata:
  type: project
---

第1轮评测发现2个中等问题：

1. **ramp_up功能未实现**：配置模型和示例文件中定义了ramp_up字段，但executor.go中完全没有使用。设置ramp_up:5时并发数直接启动而非渐进式增加。

2. **配置继承Concurrency合并缺陷**：mergeConfigs函数要求子配置必须指定mode才能覆盖Concurrency。若子配置只设置workers/duration但未指定mode，整个Concurrency不会被覆盖，导致继承的workers值不生效。

**Why**: 需求要求支持"阶梯式加压模式"和配置继承，但这两个边界情况处理不当。

**How to apply**: 下一轮需修复executor.go实现ramp_up逻辑，以及config.go改进Concurrency合并逻辑（不依赖mode字段判断是否覆盖）。