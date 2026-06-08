---
name: steering
description: 从当前项目中总结和提炼现状，形成项目公约（Steering）文档，为 Agent 提供关于本项目的核心指导规则。
metadata:
    version: 0.0.2
---

# Steering - 项目公约

本能力用于从当前项目中总结和提炼现状，形成项目公约文档。这些文档为 Agent 提供关于本项目的背景知识、项目公约、产品方向、技术栈约束和代码结构约定，作为 Agent 在后续任务中进行推理和决策的核心指导。

这些文档是项目的核心“规则”，比易变的“记忆”或临时的“日志”更具指导性和稳定性。Agent 在执行任务时，会首先加载这些项目公约文档来校准其行为，确保与项目目标和约束保持一致。

## Workflow

### Step 1: 全仓库分析

Agent 将首先对整个仓库进行全面分析，阅读和理解包括代码、现有公约、命令和技能在内的所有内容，以形成对项目当前状态的整体认知。

### Step 2: 生成更新提案

基于分析结果，Agent 将生成四份核心公约文档。在生成提案时，AI 会严格遵守各项文档的既定稳定性：

-   **`VISION.md` (产品方向)**: 极稳定。除非项目发生重大业务转向，否则不应变更。
-   **`.steering/constitution.md` (项目宪法)**: 非常稳定。通常只在团队遇到新问题并需要固化新规约时进行增补。
-   **`.steering/structure.md` (代码结构)**: 相对稳定。根据项目的实际演化进行调整。
-   **`.steering/design.md` (设计原则)**: 变化较多。反映跨项目的关键设计决策，并与 `docs/**` 中的细节设计保持同步，但不重复其细节。

### Step 3: 用户确认与写入

**重要约定：** 所有 Steering 文档均为 **动态文档 (Living Documents)**。在生成或更新这些文件时，末尾的 `Changelog` 部分必须遵循以下格式，不记录详细的变更历史：

```markdown
## Changelog

<!-- // 这是一个 Living Document，如无必要，无需维护变更历史。 -->
```

### 输出与模板

所有公约文档均基于参考模板生成，并输出到对应的目录下。在 `.agents/commands/steering/references/` 目录查看和修改这些模板：

| 文档              | 输出路径                     | 模板路径                            |
|-------------------|-----------------------------|-------------------------------------|
| `VISION.md`       | `VISION.md`                 | `references/template-VISION.md`     |
| `constitution.md` | `.steering/constitution.md` | `references/template-constitution.md` |
| `structure.md`    | `.steering/structure.md`    | `references/template-structure.md`  |
| `design.md`       | `.steering/design.md`       | `references/template-design.md`     |
