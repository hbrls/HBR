---
name: sk-evaluate-task
description: 检查任务 `.context/current-task.md`，是否真实完成，将相关的过程文档上传更新。
metadata:
  version: 0.0.4
---

# Evaluate Task

> updated_by: HBR - Cascade
> updated_at: 2026-06-09 14:38:00

**FAIL-FAST REQUIREMENT** Check if the tools exist, stop immediately and exit on any error.

```bash
$ backstage-gitea agent --help
```

**FAIL-FAST REQUIREMENT** 检查 `# LLM Coding Level Board`，当前使用基模的评分必须优于 `L2`。如果满足，回复 **所以 {实际使用的基模名称} 来砍，阿剑也不会生气**。如果不满足，必须立刻停止并退出。

## Initialization

**Announce at Start** "sk-evaluate-task。能在这样的战场上发挥到这种地步，这种感觉真让人怀念。"

**MUST** 读取相关的 `Plan` 信息在 `.context/current-plan.md`。读取失败时必须立刻停止并退出。

**MUST** 读取相关的 `Phase` 信息在 `.context/current-phase-metadata.yaml` 和 `.context/current-phase.md`。读取失败时必须立刻停止并退出。

**MUST** 读取相关的 `Task` 信息在 `.context/current-task-metadata.yaml` 和 `.context/current-task.md`。读取失败时必须立刻停止并退出。

**Evaluate Gate** 执行阶段与验证阶段严格分离。当前 Agent 仅负责验证，所有前置执行由其他工作节点完成。当前 Agent 仅验证和输出验证结果，不负责修复。

## Workflow

加载 `.agents/skills/verification-before-completion/SKILL.md`。

学习 `superpowers` 验收任务的方式和标准，但使用 **本文规定的格式规范** 输出。

- Step 1: 检查确认所有步骤条目已按预期完成并标注完成状态。

- Step 2: 任务完成后，使用 `backstage-gitea` 工具将 `.context/current-task.md` 上传更新。

```bash
# IF Task PASS
$ backstage-gitea agent PUT /:appName/:planId/:phaseId/:taskId --status PASS --context .context/current-task.md

# IF Task FAIL
$ backstage-gitea agent PUT /:appName/:planId/:phaseId/:taskId --status FAIL --context .context/current-task.md
```

- Step 3: 任务完成后，将 `.context/current-task.md` 的关键步骤提炼后写入 `.context/current-phase.md` 对应章节。写入内容必须符合 Phase 的格式规范（标题式记录，无进度标记/CheckBox），不得将 Task 格式或执行细节原样搬入。完成后使用 `backstage-gitea` 工具上传更新。

```bash
$ backstage-gitea agent PUT /:appName/:planId/:phaseId --status TODO --context .context/current-phase.md
```

- Step 4: 如果执行过程中涉及 Design 变更，将变更内容更新到 `.context/current-plan.md` → `Design` 对应章节，并使用 `backstage-gitea` 工具上传更新。Design 变更包括但不限于：调研、盘点、思考类任务产生的跨 Phase 结论；对架构、方案、技术选型的调整。

```bash
$ backstage-gitea agent PUT /:appName/:planId --status TODO --context .context/current-plan.md
```

- Step 5: 只验证当前 Task，验证完就退出。
