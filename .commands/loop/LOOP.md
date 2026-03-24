---
description: 将复杂任务拆为可验证子任务，主窗口编排、子代理逐个执行。触发词：loop、use loop、loop mode、循环执行。
---

# Loop 工具

> 工具边界见 `<SYSTEM_SKILL_ROOT>/references/tool-boundaries.md` | 共享规则见 `<SYSTEM_SKILL_ROOT>/references/shared-rules.md`

将复杂工作拆成可执行原子任务，在 Task 系统中按顺序推进。主窗口只做编排，每次仅分派一个子任务给子代理执行。

## 核心原则

- **上下文隔离**：每个任务在子代理中运行，避免主窗口上下文爆炸
- **原子任务**：每个任务可独立执行、可独立验证
- **干净交接**：子代理执行一个任务后返回；主窗口继续分派下一个

## Tool Contract

### Use when

- 任务复杂，需要拆解为多个可验证子任务
- 需要长流程持续推进（主窗口按任务状态续跑）
- 需要隔离上下文，避免主窗口持续膨胀

### Failure fallback

- `init-loop.sh` 失败：停止推进并返回错误与修复建议
- `Task` 系统异常：停止推进并输出恢复建议（重试/缩小任务/手动收尾）
- 无法满足"单任务可执行"粒度：继续拆分或补充上下文

---

## Phase 0: 简单任务检查

**Goal**: 评估任务复杂度，避免对简单任务引入不必要的流程开销

**Actions**:
1. 满足以下全部条件时，**建议直接完成**：只涉及 1-2 个文件、范围清晰无需探索、一个任务即可完成
2. 提示用户：「这个任务看起来比较简单，直接完成会更快。要现在做还是用 loop？」
3. 用户选直接完成 → 不开 loop；用户坚持 → 继续 Phase 1

---

## Phase 1: 初始化

**Goal**: 在拆分任务前准备好 loop 目录

**Actions**:
1. 运行初始化脚本（不要用 `run_in_background: true`，后续步骤立刻需要 `LOOP_DIR`）：
```
bash <SYSTEM_SKILL_ROOT>/tools/loop/scripts/init-loop.sh <slug>
```
2. **slug**：简短英文标识（如 `snake-game`、`auth-module`）
3. 记住脚本输出的 `LOOP_DIR`：
```
LOOP_DIR=.claude/skills/pensieve/loop/2026-01-27-login
```

---

## Phase 2: 捕获上下文

**Goal**: 在生成任务前记录会话上下文，确保任务拆解基于共识而非假设

**Actions**:
1. 创建 `LOOP_DIR/_context.md`，包含以下摘要字段：
   - **交互历史**：轮次 / 模型尝试 / 用户反馈（表格）
   - **最终共识**：目标、范围、约束
   - **理解与假设**：预期涉及模块、实现方式、难点
   - **文档引用**：需求文档 / 设计文档 / 计划文档（路径或"无"）
   - **上下文链接**（可选）：基于 / 导致 / 相关（`[[链接]]` 格式）
2. 向用户展示上下文摘要并**确认后再继续**
3. 按需创建文档（仅在以下条件触发时）：

   | 条件 | 需要 | 模板 |
   |------|------|------|
   | 需求不清晰（目标/范围/约束未确认） | 需求文档 | `loop/REQUIREMENTS.template.md` |
   | 实现方式不明显 | 设计文档 | `loop/DESIGN.template.md` |

4. 如果 loop 可能产出 `decision` 或 `pipeline`，预填上下文链接（规则见 `<SYSTEM_SKILL_ROOT>/references/shared-rules.md` § 语义链接规则）

---

## Phase 3: 生成任务

**Goal**: 将工作拆解为原子化、可执行的任务（前置条件：Phase 2 用户已确认）

**Actions**:
1. 读取所有项目准则 `<USER_DATA_ROOT>/maxims/*.md`，用准则约束任务边界和验收标准
2. 读取 `<USER_DATA_ROOT>/pipelines/` 下所有 `*.md`，检查是否包含任务蓝图（`## Task Blueprint` 或有序 `### Task 1/2/3...` 标题）
   - 有蓝图 → 1:1 按序映射为运行时任务
   - 无蓝图 → 按粒度标准拆分
3. **粒度标准**：子代理能否不提问就执行完？每个任务需指定目标文件/组件，包含具体的构建/修改/测试操作
4. 确保验收标准与准则对齐
5. 在 Claude Task 系统中创建任务（markdown checklist 不等于任务创建，只有在 Task 系统中创建的才算数）
6. 展示简要快照（task id + 主题），然后创建/运行第一个任务

---

## Phase 4: 主窗口续跑

**Goal**: 持续推进直到所有任务完成

**Actions**:
1. 第一个任务创建后，主窗口获取下一个待处理任务，每次分派一个子代理
2. 不依赖 hooks 或后台绑定进程，主窗口主动轮询任务状态

---

## Phase 5: 执行任务

**Goal**: 通过隔离的子代理逐个执行任务

**Actions**:
1. 为第一个待处理任务启动通用子代理：
```
Task(
  subagent_type: "general-purpose",
  prompt: "读取 .claude/skills/pensieve/loop/{date}-{slug}/_agent-prompt.md 并执行 task_id={id}"
)
```
   `_agent-prompt.md`（由 init-loop.sh 生成）包含角色定义、loop 上下文路径、准则路径、执行约束。
2. 子代理：TaskGet → 执行 → 返回
3. 主窗口检查任务状态，分派下一个待处理任务

---

## Phase 6: 收尾

**Goal**: 完成沉淀与知识写入

**Actions**:
1. 所有任务完成 → 询问是否运行 self-improve（`tools/self-improve/_self-improve.md`）
2. 如果 loop 产出了 `decision` 或 `pipeline`，确保输出包含链接（见 `<SYSTEM_SKILL_ROOT>/references/shared-rules.md`）

---

## Phase 选择指南

| 任务特征 | Phase 组合 |
|----------|------------|
| 范围清晰、规模小 | 直接生成任务 |
| 需要理解代码 | 计划 → 任务 |
| 需要技术设计 | 计划 → 设计 → 任务 |
| 需求不清晰 | 计划 → 需求 → 设计 → 任务 |

---

## 相关文件

- `tools/loop/README.md` — 详细文档
- `<SYSTEM_SKILL_ROOT>/tools/loop/scripts/init-loop.sh` — 初始化 loop 目录