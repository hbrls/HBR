# Project Constitution

> updated_by: Windsurf - Cascade
> updated_at: 2026-02-21 14:36:00

## Purpose

本文档为 RE10 项目定义了核心的治理原则、决策流程和质量标准，旨在确保 AI 与人类协作者之间的高效、可预测和高质量的互动。

## Core Values

1.  **AI 原生 (AI-Native)**: 所有流程、文档和工具链的设计都优先考虑 AI 智能体的参与和自动化执行，追求机器可读性和无歧义的指令。
2.  **明确高于默契 (Explicit over Implicit)**: 依赖结构化的模板、元数据和清晰的文字记录来沟通意图，减少因上下文模糊或假设不一致导致的错误。
3.  **流程即代码 (Process as Code)**: 将开发和协作的最佳实践固化为可执行的“命令”和“技能”，通过代码化的方式实现流程的持续改进和治理。

## Decision-Making Principles

### When to Write a Spec

在以下情况需要编写 Spec 文档：

- 引入新的 `command` 或 `skill`
- 对现有 `command` 或 `skill` 的核心行为进行重大修改
- 变更会影响多个 `command` 之间的交互
- 预计工作量超过 **3 人日**

在以下情况无需编写 Spec 文档：

- 修复 Bug 且范围明确
- 更新文档
- 不改变核心行为的重构

### Approval Process

| 变更类型 | 审批人 | 审阅时间 |
|---|---|---|
| 次要 | 1 位团队成员 | 1 天 |
| 标准 | 2 位团队成员 | 3 天 |
| 重大 | 项目负责人 + 1 位核心成员 | 5 天 |
| 破坏性 | 所有利益相关者 | 1 周 |

### Conflict Resolution

当出现分歧时，按以下顺序解决：

1.  **数据优先**: 以可度量的数据、用户反馈或基准测试结果为依据。
2.  **原型验证**: 如果不确定，构建最小化的原型来验证方案。
3.  **限时决策**: 为决策设置一个明确的截止日期。
4.  **升级路径**: 若仍无法解决，由 **项目负责人** 做出最终决定。

## Quality Standards

### Code Quality

- 所有代码提交前必须通过 Lint 检查。
- 核心库的单元测试覆盖率不低于 **80%**。
- 不引入任何已知的严重安全漏洞。
- 所有对外暴露的 API 或接口都必须有文档说明。

### Review Criteria

- 满足 Spec 中定义的验收标准。
- 不破坏现有功能。
- 性能可接受 (参见 `design.md` 中的目标)。
- 遵循 `structure.md` 中定义的代码风格和结构约定。

## Communication Guidelines

### Spec Updates

- 当 Spec 状态变更时，通知所有利益相关者。
- 在开发过程中，每日更新 `.plans/` 中对应的计划文件状态。

### Blocking Issues

- 在发现阻塞性问题后 24 小时内提出。
- 在 **GitHub Issues** 中创建 issue 并标记相关人员。
- 在对应的 Spec 或 `.plans/` 目录下的计划文件中注明阻塞状态。

## Changelog

<!-- // 这是一个 Living Document，如无必要，无需维护变更历史。 -->
