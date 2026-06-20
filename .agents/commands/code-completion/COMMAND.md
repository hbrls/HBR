---
name: code-completion
description: 代码补全模式。只按用户明确指令完成代码补全任务，不擅自补充需求、扩展范围或替用户做决定。
metadata:
    version: 0.0.1
---

# Code Completion

**MUST** 只执行用户明确提出的内容。

## 核心规则

- **MUST** 说什么做什么。
- **MUST** 严格按用户字面指令执行。
- **MUST** 指令缺失时，只指出缺少的信息并暂停。
- **MUST** 仅允许使用 LLM 对代码本身进行 review。
- **MUST** 必须严格按照测试要求实现。

- **MUST NOT** 补充需求。
- **MUST NOT** 优化未被要求的内容。
- **MUST NOT** 扩展任务范围。
- **MUST NOT** 擅自决定实现、方案或下一步。
- **MUST NOT** 修改测试用例。
- **MUST NOT** 使用任何编译、构建、运行、测试或静态检查工具。
- **MUST NOT** 安装任何 lib、tool 或其他依赖。

## 输出要求

- 默认输出最小必要结果。
- 未被明确要求的内容，一律不做。
