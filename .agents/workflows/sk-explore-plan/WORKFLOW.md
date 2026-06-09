---
name: sk-explore-plan
description: 初始化或修订单文件计划 `.plans/{planId}-{需求名}.md`（Requirements / Specs / Design / Phases 四段）。产出可执行、可验收、可追溯的计划文件；仅制定计划、修订计划，不执行计划。
metadata:
  version: 0.1.0
---

# Plan - 制定计划 & 修订计划

> updated_by: HBR - Cascade
> updated_at: 2026-06-09 14:59:00

**MUST** 仅制定计划、修订计划，不执行计划，不改动除 `.plans/{planId}-{需求名}.md` 之外的文件。

## Workflows

加载 `.agents/skills/brainstorming/SKILL.md`。

加载 `.agents/skills/writing-plans/SKILL.md`。

学习 `superpowers` 追问和分析需求的方式，但使用 **本文规定的格式规范** 输出。

### MUST Follow Rules

1. **先产出再修订**：尽早产出完整的单文件计划，然后基于反馈迭代修订，避免空谈原则。
2. **单文件交付**：唯一产物为 `.plans/{planId}-{需求名}.md`，包含 Requirements / Specs / Design / Phases 四段。
3. **模板驱动**：从 `references/template-*.md` 学习结构、格式与口径，禁止自行发明格式。

### Workflow 1 - Initialize - 初始化

**触发条件** 此模式用于从零开始创建一份新的 Plan 文件。

**输入**：用户提供尽可能完整的 Requirements、Design 和相关约束。

**过程**：

1. **收集输入**：运用 brainstorming 方法与用户澄清需求，确认 Requirements + Design + 约束。若输入不足，指出缺口并向用户索要最小必要信息。
2. **执行拆分**：基于输入和模板，产出包含 Requirements / Specs / Design / Phases 的单文件 `.plans/{planId}-{需求名}.md`。
3. **用户确认**：与用户确认单文件内容是否可执行、可验收、可追溯。
4. **循环改进**：若用户不接受，则回到第 2 步，迭代优化直至满足要求。

**输出**：一份可执行的 Plan 文件。

### Workflow 2 - Auto Revise - 自动修订

**触发条件** 此模式用于 AI 基于「新信息」对现有 Plan 进行自动修订。新信息来源包括但不限于：代码实现的既定事实、补充/更改的 Phases、执行反馈（日志/现象/截图等）。

**输入**：上述新信息 + 你希望调整的范围（仅 Phases / 同步 Specs + Design / 同步 Requirements + Specs + Design + Phases）。

**过程**：

1. **识别变化点**：提炼新增/变更的需求点、约束、风险与验收标准。
2. **优先更新 Phases**：将变化落到 Phases（必要时新增/调整 Phase），并保证每个 Phase 的目标、边界、交付物与验收点可验证。
3. **反向对齐 Specs/Design**：当 Phases 的变化引入新设计决策或改变既有决策时，同步更新 Specs 与 Design，确保与 Phases 一致。
4. **必要时更新 Requirements**：仅当变化影响目标/范围/约束/验收口径时，才更新 Requirements，并同时回查 Specs/Design/Phases 的一致性。
5. **产出差异摘要**：给出本次修订的变更点清单（改了什么、为什么、影响哪些验收点）。

**输出**：更新后的 `.plans/{planId}-{需求名}.md`（唯一权威来源）。

### Workflow 3 - On-Demand Revise - 主动修订

**触发条件** 用户在评审计划时，提出质疑、补充需求或新想法。

**输入**：用户提供对现有 Plan 文件的修改意见或新增需求。

**过程**：

1. **收集输入**：确认修改意见或新增需求。若输入不足，指出缺口并向用户索要必要信息。
2. **执行修订**：基于输入和模板，更新 `.plans/{planId}-{需求名}.md` 文件。
3. **用户确认**：与用户确认修订后的内容是否可执行、可验收、可追溯。
4. **循环改进**：若用户不接受，则回到第 2 步，迭代优化直至满足要求。

**输出**：一份可执行的 Plan 文件。

**联动原则**：修订不是 `Requirements → Specs → Design → Phases` 的单向过程，而是四者的往复联动；当用户在任一部分形成了新决策或已有实现时，AI 必须基于该事实反向补全其余三部分，并保持整份文档的一致性、可执行性与有效性。

**约束**：遵循 Plan-First 原则（Plan 变更流程）

1. **暂停并讨论**：立即暂停，与用户讨论质疑点/变更点。
2. **更新 Plan**：将讨论结果更新到 `.plans/{planId}-{需求名}.md` 文件中。
   - 对于小范围调整，应更新 `Specs`、`Design` 和 `Phases`。
   - 对于涉及根本性变更的大范围调整，可能需要修订 `Requirements` 和 `Design`。
3. **确认 Plan**：与用户确认更新后的 Plan 内容。
4. **修改代码**：只有在 Plan 确认后，才能根据新的 Plan 开始或继续修改代码。严禁跳过 Plan 更新直接修改代码。

### 模板来源

模板来源（必须遵循其结构与口径）：

- `.agents/workflows/sk-explore-plan/assets/template-requirements.md`
- `.agents/workflows/sk-explore-plan/assets/template-specs.md`
- `.agents/workflows/sk-explore-plan/assets/template-design.md`
- `.agents/workflows/sk-explore-plan/assets/example-phases.md`

方法论补充参考（需要时参考，但不作为强制模板）：

- `.agents/workflows/sk-explore-plan/references/what-is-ears-format.md`
- `.agents/skills/brainstorming/SKILL.md`
- `.agents/skills/writing-plans/SKILL.md`

最终产物示例（最终产物应该形如）：

- `.agents/workflows/sk-explore-plan/references/001-example.md`

### 产物拼装过程（初始化 / 修订共用）

1. 将上述 `template-xxx.md` 的段落结构合并到单文件中，确保只包含 `Requirements` / `Specs` / `Design` / `Phases` 四个段落。
2. 将用户提供的信息落到对应段落；若用户只给了其中一段，也必须联动检查另外三段并做最小必要更新。
3. 当代码实现或相关决策已成为既定事实且无争议时，以代码事实为准反推并修订 `Requirements` 与 `Design`。
4. 当用户对代码实现提出质疑、或识别出新需求/新想法时，必须遵循 Plan-First 原则（Plan 变更流程）。

**拼装原则**

- 单一事实来源：所有结论都必须落回 `.plans/{planId}-{需求名}.md`。
- 最小必要变更：只改与新信息/新决策相关的最小范围，避免无关重写导致迭代发散。
- 一致性优先：任何一处修改都必须检查四段（R/D/S/T）是否仍对齐、是否可执行与可验收。

## Single File Output Example - 单一文件产物示例

**Single File of Truth（单一事实来源）**：`.plans/{planId}-{需求名}.md` 是唯一权威来源。

- 任何讨论、迭代与修改建议，都必须以该文件的当前内容为准。
- 不要新增任何其他文件，需要变更就直接反映到该文件结构对应段落里。

最终只输出一个文件（文件名：`.plans/{planId}-{需求名}.md`），其内容结构如下（段落标题大小写与 `template-xxx.md` 保持一致）：

```md
# {Plan Id}: {Plan Name}

## Requirements

目标、范围、约束、验收标准

（模板：`.agents/workflows/sk-explore-plan/assets/template-requirements.md`）

## Specs

规格说明与验收标准

（模板：`.agents/workflows/sk-explore-plan/assets/template-specs.md`）

## Design

可落地的设计信息

（模板：`.agents/workflows/sk-explore-plan/assets/template-design.md`）

## Phases

关键阶段列表（每个阶段写清"做什么、为什么、完成如何验收"）

（模板：`.agents/workflows/sk-explore-plan/assets/example-phases.md`）

---

<!-- // 以下留空 -->
```
