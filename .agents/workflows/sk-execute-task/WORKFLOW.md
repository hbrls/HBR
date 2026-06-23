---
name: sk-execute-task
description: 执行当前任务 `.context/current-task.md`，更新本地过程文档，为后续验收准备交付产物。
metadata:
  version: 0.0.5
---

# Execute Task

**FAIL-FAST REQUIREMENT** 检查 `# LLM Coding Level Board`，当前使用基模的评分必须优于 `L2`。如果满足，回复 **所以 {实际使用的基模名称} 来砍，阿剑也不会生气**。如果不满足，必须立刻停止并退出。

## Initialization

**Announce at Start** "sk-execute-task。能在这样的战场上发挥到这种地步，这种感觉真让人怀念。"

**MUST** Read the relevant `Plan` information from `.context/current-plan-metadata.yaml` and `.context/current-plan.md`. On read failure, stop immediately and exit.

**MUST** Read the relevant `Phase` information from `.context/current-phase-metadata.yaml` and `.context/current-phase.md`. On read failure, stop immediately and exit.

**Source of Truth** `Current Task` 的身份、状态与执行资格以 `.context/current-task-metadata.yaml` 与 `.context/current-task.md` 为信息源，二者共同构成当前任务的单一事实源。未上传、未落库的信息不构成有效状态，也不能覆盖信息源的定义。

**Local Working Copies** `.context/current-phase.md`、`.context/current-plan.md` 仅用于读取上下文，其中的说明、人工验收文案都只是本地工作副本，可能是过期副本。副本与信息源冲突时，以信息源为准。

**Explore Gate** 探索阶段与执行阶段严格分离。当前 Agent 仅负责执行，所有前置探索由其他工作节点完成。**执行前须确认**：任务详情中存在探索已完成、可以开始执行的明确信号。若信号缺失，立刻停止并退出——这说明工作流程调用有误。

## Workflow

- Step 1: **理解并确认任务目标** 在开始执行前，必须先核对 `.context/current-task-metadata.yaml` 与 `.context/current-task.md`，确认当前正在执行的是后台定义的 Current Task；随后明确回答以下三个问题：
  1. 任务的真实目标是什么？（从 Task 名称和 Do 列表推断）
  2. 为什么要执行这个任务？（从 Phase 上下文理解）
  3. 完成标准是什么？（从 Check 列表确认）

  **MUST** 输出固定格式的确认信息：

  |     | 输出内容    | 自信度 |
  |-----|------------|-------|
  | 目标 | {推断的目标} | 1-9分 |
  | 原因 | {推断的原因} | 1-9分 |
  | 标准 | {推断的标准} | 1-9分 |

  如果无法明确回答这三个问题，或发现本地 Task 内容与 Metadata 指向的当前任务不一致，**立刻停止并询问用户**。

- Step 2: 开始执行任务 `.context/current-task.md`。按顺序执行任务，从第一个 `- [ ]` 开始。一次只执行一个步骤条目。完成后将对应条目从 `- [ ]` 更新为 `- [x]`。

  - **输出**：代码实现、测试结果或其他交付产物。

  **Execute Rules**

  - **MUST** 执行阶段只修改工作区与本地过程文档，不上传后端状态，不宣告 Task 最终完成；执行收尾后只应进入待验收状态。
  - **MUST** `sk-execute-task` 必须具备可重复执行的特性：同一 Task、同一步骤在重复执行时，执行方案本身必须能够收敛到同一目标状态，而不是依赖“该操作只会被执行一次”的前提。
  - **MUST** 应优先采用天然支持重复执行的操作方式，例如替换、覆盖、同步、重建、对齐、幂等写入；禁止把“先检查是否做过、如果做过就跳过”作为主要保障。
  - **MUST** 如果某个步骤天然带有不可重复的副作用，必须先调整执行方案使其具备可重复执行特性；若无法做到，立刻停止并询问用户。
  - **MUST** 对本地过程文档的更新应采用“覆盖/同步到对应章节”的方式，确保重复执行时文档内容仍然收敛，避免产生多份并列结论或累计式日志污染。
  - **MUST** 执行过程中如识别到关键变更，必须立刻停止并退出。这种场景说明 Design 需要更新。
  - **MUST** 同一个问题发生三次错误时必须立即停止执行，禁止在同一条回复中完成分析并执行修改。必须先分析，等待用户确认后再执行修改。
  - **MUST** 现场必须收敛：执行、验证、可见性验收过程中，由本次执行或被调用 `skill` 引入的临时运行时资源，不得作为待验收状态的隐式依赖残留。`sk-execute-task` 不规定具体启停方式，但必须要求启动方完成生命周期闭环；若无法确认现场已经收敛，视为执行未完成，不得宣称已具备验收条件。
  - **违反后果** 立即触发 `.agents/commands/andon`，进入 Andon 状态

  **Thinking 显化规则**

  - **MUST** Thinking 过程必须在对话中显化打印，必须在工具调用前打印。
  - **MUST** 关键决策点（遇到错误、需要判断、用户提问、改变执行方向、放弃某个方案）- 必须打印分析过程，必须打印决策依据。

  格式：
  ```
  【Thinking】
  - 当前状态: {描述当前进度和上下文}
  - 下一步动作: {描述将要执行的操作，明确写出将要使用的工具名称，明确写出将要执行的动作摘要}
  - 预期结果: {描述操作成功后的预期状态}
  - 风险评估: {描述潜在风险和应对方案}
  - 回滚策略: {如失败如何回滚}
  ```

  示例：
  ```
  【Thinking】
  - 当前状态: 已完成步骤 1-3，准备修改 minimax_provider.py
  - 下一步动作: 使用 replace_in_file 修改 __init__ 方法签名
  - 预期结果: 移除 api_key 参数，使用 Config() 获取配置
  - 风险评估: replace_in_file 的 SEARCH 块必须精确匹配原文件内容
  - 回滚策略: 如修改失败，恢复本次修改前的文件内容，或停止并等待用户确认

  执行 ...
  ```

  **执行质量检查**：

  1. 每个步骤执行前，确认理解该步骤的具体要求
  2. 关注"可执行、可验收、可追溯"，验收与检查必须落到可测试/可检查标准（可复现、可度量、可对照）。
  3. 如步骤涉及外部资源（API、文档、代码），必须明确打印出资源位置和使用理由
  4. 每一个步骤条目完成后，立即检查产出是否符合预期，如不符合应立即修正

- Step 3: 任务产生的所有过程产出，写入 `.context/current-task.md` 对应章节。**不允许在 `.context/` 下新建文件**。

  **过程质量检查**：

  1. 过程日志必须结构清晰、格式规范
  2. 内容必须基于实际发生，禁止虚构
  3. 产出物应包含当时的场景，必要的示例、说明，便于后续使用

- Step 4: **Completion Gate** 检查确认现场已经收敛：所有步骤条目已按预期完成并标注完成状态；验收证据已写入 `.context/current-task.md`；本次执行或被调用 `skill` 引入的临时运行时资源不存在未说明、未移交、未处理的依赖。只有通过该门禁，才可声明“本地执行已收敛，已具备验收条件”；该状态不表示后端状态已更新。

- Step 5: **完成后必须退出** 只执行当前 Task；通过 `Completion Gate` 后必须退出，不继续调试、不等待观察、不进入下一个 Task、不维持任何临时运行时资源作为交付状态。
