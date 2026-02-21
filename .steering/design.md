# Design Principles

> updated_by: Windsurf - Cascade
> updated_at: 2026-02-21 14:36:00

## Design & Architecture

RE10 是一个 AI 原生工作流框架，其设计分为定义层和执行层。

### Languages

| Purpose | Language | Version | Notes |
|---|---|---|---|
| 定义层 | Markdown | - | 用于编写所有 `commands` 和 `skills`。 |
| 执行层 (假设) | Python | 3.9+ | 用于未来可能开发的、解析和执行工作流的脚本。 |

### Frameworks

本项目作为框架本身，不依赖外部应用框架。

| Purpose | Framework | Version | Notes |
|---|---|---|---|
| Testing (假设) | pytest | - | 用于未来可能开发的执行层逻辑测试。 |

### Infrastructure

本项目在本地开发环境中运行，不依赖特定的基础设施组件。

## Performance Targets

不适用。本项目的性能取决于执行其工作流的 AI 智能体的能力。

## Security Requirements

- **执行安全**: 所有由 AI 提议的 shell 命令都必须经过用户审批才能执行。
- **数据保护**: 不处理或存储用户个人身份信息 (PII)。

## Dependencies

### External Services

| Service | Purpose | Criticality | Fallback |
|---|---|---|---|
| AI Agent (如 Cascade) | 解析和执行本框架中定义的工作流。 | High | 无 |
| IDE (如 VS Code) | 提供用户交互界面和工具链。 | High | 无 |

## Constraints

### Must Use

- **Markdown**: 用于定义所有 `commands` 和 `skills`，以确保机器可读性和人类可读性。

### Must Not Use

- **二进制格式**: 避免使用任何非文本、非结构化的格式来定义工作流组件。
- **硬编码逻辑**: 应将逻辑封装在可复用的 `skills` 中，而不是在 `commands` 中硬编码。

### Prefer

- **声明式定义**: 优先使用声明式的方式描述工作流和规则，而不是命令式的脚本。

## Changelog

<!-- // 这是一个 Living Document，如无必要，无需维护变更历史。 -->
