---
name: snippetbox-eval-43067
description: Go+Cobra代码片段管理CLI工具R1评测：编辑加密片段加密被移除+删除默认vault后配置失效
metadata:
  type: project
---

## SnippetBox 代码片段管理CLI评测 R1

**技术栈**: Go + Cobra + JSON文件存储

**评测结果**: FAIL

### 发现的问题

1. **编辑加密片段时加密状态被静默移除**：编辑已加密片段但不传 `--encrypt` 参数时，片段加密标志被设为false，导致内容以明文形式保存。编辑逻辑应该保留原加密状态，除非用户明确指定改变。

2. **删除默认vault后配置指向不存在的vault**：删除被设为默认的vault后，配置中DefaultVault被重置为"default"字符串，但该vault可能不存在，导致后续所有片段操作失败。

### 通过的功能

- 片段CRUD操作完整
- 多保险库支持完整
- 标签系统完整（多标签、筛选、列出）
- 搜索功能完整（模糊匹配、多字段、标签/语言筛选）
- 模板变量交互式填充
- JSON导入导出
- AES-GCM加密实现
- 配置文件和日志记录
