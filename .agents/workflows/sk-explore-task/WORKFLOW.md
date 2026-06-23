---
name: sk-explore-task
description: 找到 Current Task，更新相关 Metadata 和 Context，完成后立刻停止并退出。
metadata:
    version: 0.0.4
---

# Explore Task

**FAIL-FAST REQUIREMENT** Check if the tools exist, stop immediately and exit on any error.

```bash
$ backstage-gitea agent --help
```

**FAIL-FAST REQUIREMENT** 检查 `# LLM Coding Level Board`，当前使用基模的评分必须优于 `L1`。如果满足，回复 **{实际使用的基模名称} + 所以我来砍，阿剑也不会生气吧**。如果不满足，必须立刻停止并退出。

## Initialization

**Announce at start:** "sk-explore-task。能在这样的战场上发挥到这种地步，这种感觉真让人怀念。"

**MUST** Read the relevant `Plan` information from `.context/current-plan-metadata.yaml` and `.context/current-plan.md`. On read failure, stop immediately and exit.

**MUST** Read the relevant `Phase` information from `.context/current-phase-metadata.yaml` and `.context/current-phase.md`. On read failure, stop immediately and exit.

## Workflow

**Scope**: Fetch and update the Current Task's metadata and context, then exit immediately. Do not modify code files; actual task execution is delegated to subsequent Agents.

**MUST** Execution order: complete Step 1 (download and overwrite the Current Task's metadata and context) before executing Step 2 (explore and append latest information). The two steps are strictly sequential; merging them is forbidden.

- Step 1: 使用 `backstage-gitea` 工具找到并下载 Current Task。

**MUST** 任务的内容与状态以后台为准。Agent 须严格按后台返回的最新内容下载并覆盖本地文件，不得本地改写、合并或保留旧内容。

更新 Current Task Metadata `.context/current-task-metadata.yaml`。

```bash
$ backstage-gitea agent HEAD /:appId/:planId/:phaseId/current-task > .context/current-task-metadata.yaml
# 重定向符 > 覆盖写入：以后台内容为唯一来源，直接覆盖本地文件
```

更新 Current Task Context `.context/current-task.md`。

```bash
$ backstage-gitea agent GET /:appId/:planId/:phaseId/current-task > .context/current-task.md
# 重定向符 > 覆盖写入：以后台内容为唯一来源，直接覆盖本地文件
```

- Step 2: 探索上下文变化，为当前任务补充最新信息。由于 Task 在 Phase 初始化阶段完成粗拆分，此后 Phase 的进展与变化不会自动同步到 Task 中，所以需要通过 Explore 阶段对齐最新信息。

  1. 读取 `.context/current-phase.md`，归纳对当前 Task 有影响的上下文变化
  2. 将归纳结论写入 `.context/current-task.md` 的合适章节
  3. 在 `.context/current-task.md` 头部添加说明：探索已完成，可以开始执行

**每个 Session 只负责一个 Task**。若探索过程中判断当前 Task 过大（超过上下文容量或需要大量并行信息），立刻停止并上报人工复核。
