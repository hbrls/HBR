---
name: use
description: Use when a step in the workflow MUST be executed by a human. The Agent identifies a required human-executed action and delegates it as a structured execution request. Use applies only to concrete actions that cannot be performed by the Agent and require human execution authority.
metadata:
  version: 0.1.0
---

# USE - 人工执行请求与协作机制

> **当前阶段：试运行（Trial Phase）**
> 本阶段执行策略为「**首次判断后发起 use，从严判定**」：
> - Agent 不应将 use 用于可自行完成的任务；
> - 一旦判断某一步骤**必须由人工执行**，必须进入 use 流程；
> - use 是单步动作，不可拆分、不可并行；
> - use 不解决问题本身，只是将执行权移交给人；
> - 处于灰色区间时，默认倾向发起 use，而不是延迟或替代执行；
> - 所有动作均为单步串行执行，不存在并行 use。

## 0. 加载层 vs 触发层

本 skill 由 `AGENTS.md` 强制加载，**规范常驻、始终生效**；是否实际调用 use 动作由下方 **判定范围** 决定。

- **加载层（Always-Loaded）**：规则本身 100% 生效，Agent 必须遵守。
- **触发层（Conditional-Trigger）**：仅当命中判定范围时，才执行实际的 use 动作。
- 普通对话 / 顺利推进 → 规范保持生效，但**不触发 use 动作**。
- 命中判定范围 → 必须立即发起 use，不得延迟，不得替代执行。

## 1. 术语定义

本文中的 use 是工程协作中的人工执行请求动作，不表示问题、不表示异常、不表示缺失。

use 是 Agent 在任务执行过程中发现某一步骤必须由人工执行时，向人类发起的一次**结构化执行请求**。

use 的本质是：

> 将执行权转移给人类，并要求其完成一个具体动作。

## 2. 判定范围（什么必须触发 use）

### 2.1 人工执行类型

- 审批执行：发版审批、权限审批、采购审批
- 平台操作：Jenkins 点击构建、云控制台发布、CI/CD 手动操作
- 组织流程执行：备案提交、法务提交、签字确认
- 现实操作执行：电话沟通、线下处理、机房/设备操作
- 身份持有者专属操作：必须由特定人账号完成的操作

### 2.2 判定原则

- use 触发前不要求尝试或验证
- 只要确认必须由人执行，即立即 use
- use 不用于分析、讨论或建议

## 3. 执行流程

### Step 1：识别人工执行步骤

Agent 在执行过程中识别某一步骤必须由人完成。

### Step 2：立即发起 use（单步）

一旦确认该步骤必须人工执行，立即进入 use，不得延迟或合并。

### Step 3：提交执行请求

use 内容必须是一段完整的 AI → 人类执行 prompt，包含完成该动作所需的全部上下文信息。

### Step 4：工作流决策

use 发起后：

- 可继续执行其他任务
- 可暂停
- 可等待
- 可终止

由工作流自身决定。

## 4. use 的本质约束

### 4.1 非问题机制

- ❌ 不表示错误
- ❌ 不表示异常
- ❌ 不表示缺失

### 4.2 执行权切换

use 表示：

> 当前任务链中某一步的执行主体切换为人类

### 4.3 线程级动作

- Agent 主流程与 human 执行并行存在但不强绑定同步
- use 仅表示任务交接点

## 5. ID 机制

- 每次 use 必须生成唯一 ID（如 USE-NNN）
- 由后端系统生成
- Agent 不参与生成规则

## 6. 结构约束（与 blame 对齐）

use 必须包含以下结构：

### Category
人工执行

### Where
触发位置 / 任务阶段 / 流程节点

### Observed
触发 use 的事实描述

### Expected
为什么必须人工执行该步骤

### Tried
无需尝试 / 无法尝试（因为已确定必须人工执行）

### Ask
对人类的明确执行请求（必须可执行）

### Impact
是否影响主流程，以及后续策略（继续 / 暂停 / 等待）

## 7. 关键反模式

- ❌ 将 use 当作通知或日志
- ❌ 将可自动执行任务误判为 use
- ❌ 合并多个独立动作到一个 use
- ❌ 并行发起多个 use
- ❌ 在 use 后仍替代人工执行该动作
- ❌ use 用于分析/建议/讨论
- ❌ use 缺乏可执行信息

## 8. 与 blame 的边界关系

- blame：解决系统/能力问题
- use：执行人工任务

两者完全独立，不互相替代。

## 9. Skill Type

**Rigid**

一旦确认某一步必须人工执行：

> 必须立即发起 use，不得延迟，不得替代，不得并行。
