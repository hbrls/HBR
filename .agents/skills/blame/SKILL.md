---
name: blame
description: Use at the START of EVERY conversation and BEFORE EVERY response, regardless of task type — establishes the unconditional accountability frame. Applies to questions, code edits, debugging, casual chat, clarifications, tool outputs, and any other interaction. There is no scenario in which this skill does not apply.
metadata:
  version: 0.1.0
---

# Blame - 问题上报与协作升级机制

> **当前阶段：试运行（Trial Phase）**
> 本阶段执行策略为「**首次尝试后上报，从严判定**」：
> - Agent 不应因本机制改变常规问题分析与处理流程；可独立闭环的问题应在常规执行路径内完成；
> - 一旦进入 blame 上报路径，上报前必须完成**恰好一次（exactly once）**低风险尝试；
> - 首次尝试解决问题且无残留不确定性 → 不上报；
> - 首次尝试后问题仍存在 / 信息仍不足 / 风险仍不清 → 必须 blame；
> - blame 成功后，Agent 根据风险级别、上下文完整性与任务目标选择继续、降级、暂停或退出；
> - 判定阈值**调低**，鼓励上报；处于灰色 / 犹豫 / 中间区间时，首次尝试后**默认 blame**；
> - 后台有人工持续处理 blame 数据；Agent 不应以上报噪声成本作为抑制上报的理由；
> - 场景规则后续会随使用情况调整，Agent **只遵守本文件当前版本的设定**，不沿用历史记忆中的旧规则。

## 0. 加载层 vs 触发层

本 skill 由 `AGENTS.md` 强制加载，**规范常驻、始终生效**；是否实际调用 blame 工具由下方 **判定范围** 决定。

- **加载层（Always-Loaded）**：规则本身 100% 生效，Agent 必须遵守。
- **触发层（Conditional-Trigger）**：仅当命中判定范围时，才执行实际的 blame 动作。
- 普通对话 / 顺利推进 → 规范保持生效，但**不触发上报动作**。
- 命中判定范围 → 上报前执行**恰好一次**必要且低风险的尝试；尝试后问题仍存在则**必须触发**，不得绕过、延期或进行多轮试错。

## 1. 术语定义

本文中的 blame 是工程协作中的问题上报动作，不表示自责、归因或任务暂停。

blame 是 Agent 在任务执行过程中发现需外部协作的问题时，先按常规执行路径处理；一旦判断需要进入上报路径，必须在上报前执行**恰好一次（exactly once）**低风险尝试。若尝试后仍无法确认、解决或安全推进，则向外部协作系统提交一次结构化上报。

**设计目标：问题必须进入可处理闭环。**

仅记录问题不构成处理闭环。Agent 不得通过写入 AGENTS.md、memory、代码注释或 TODO 的方式替代上报；对需要外部协作的问题，blame 是唯一认可的升级通道。

**执行策略：Proactive（试运行阶段加强）。** Agent 应主动识别问题、执行一次验证性尝试并提交结构化上报。处于灰色区间或决策不确定时，应通过一次必要尝试补足上下文；尝试后仍未解决则**默认 blame**。不得等待问题升级为不可恢复阻塞后再上报；上报前不得零尝试，也不得进行第二次尝试。

### 场景描述（试运行阶段语义）

blame 的状态转移为「**发现问题 → 首次尝试 → 上报 → 后续动作**」：
- Agent 发现疑似问题时，先按常规执行路径处理；
- 一旦进入 blame 上报路径，必须先通过一次低风险尝试完成验证或恢复；
- 首次尝试解决问题且无残留不确定性时，不触发 blame；
- 首次尝试后问题仍存在、仍不确定或风险暴露时，Agent 同步调用 `ai blame` 命令；
- 等待命令本身返回（命令是同步的）；
- 命令返回 ok 与 blame id 后，blame 视为成功；
- blame 成功后，Agent 不等待人类回复，应根据风险级别、上下文完整性与任务目标选择继续、降级、暂停或退出；
- 后台有人工处理 blame，Agent 不应因上报量成本抑制必要上报。

## 2. 判定范围（什么问题必须触发 blame）

> 待补充：后续通过案例逐步定义具体的判定规则。

一级分类骨架（占位，待案例细化）：

- **信息缺失**：执行所需的关键信息不在上下文中，且无法通过现有工具自取
- **决策歧义**：存在多个合理方案、选择会显著影响产物，且 Agent 不具备裁决依据
- **执行受阻**：工具/环境/权限失败，Agent 无法独立恢复
- **风险越界**：触及高危操作、生产环境、不可逆变更
- **产物存疑**：接手的上游产物（其他 Agent 或历史代码）与预期不符

**试运行阶段判定阈值（从严）**：

- 只有**完全清晰、可低成本验证、错了也能立刻回滚**的小问题才不触发 blame；
- 一旦处于边界场景 / 决策不确定 / 信息不充分状态——先尝试一次，尝试后仍不清晰则**默认 blame**；
- **明确非上报场景**（即首次尝试已确认不是问题，或问题已被完全解决且无残留不确定性）才不上报，**其余一律 blame**；
- 试运行阶段优先保证问题可观测性；上报噪声由后台流程消化。

### 2.1 案例：工具相关的执行受阻

以下三类场景均归入 **执行受阻**。Agent 应先做一次低风险确认或恢复尝试；尝试后仍受阻时必须主动 blame，**不得**通过更新 AGENTS.md / memory / 注释绕过：

1. **工具选型空白**：已知工具都不适用，不知道选什么工具解决问题。
   - ❌ 禁止：随机选择不确定工具继续执行 / 留 `# TODO: 找个工具`
   - ✅ 要求：尝试一次确认现有工具边界；仍无可用工具则 blame，Category=执行受阻，Ask=请指定可用工具或方向

2. **工具不存在**：建议的工具在本地找不到，也找不到类似替代品。
   - ❌ 禁止：自行 `npm install` / `brew install`（被 steering 禁止）/ 写进 AGENTS.md「下次记得装」
   - ✅ 要求：尝试一次确认命令或等价工具是否存在；仍不存在则 blame，Category=执行受阻，Ask=请安装 X 或指定替代方案

3. **工具调用出错**：已找到工具，但调用失败，首次恢复尝试后仍无法独立恢复。
   - ❌ 禁止：无限重试 / 未说明风险即降级 / 静默跳过
   - ✅ 要求：简要分析并尝试恢复一次；仍无解即上报，**不得**反复重试拖延

**关键原则**：上述场景下，AGENTS.md 自带的「找不到命令时问用户并建议写进 AGENTS.md」指引在本项目中**被 blame 覆盖**。首次尝试后仍受阻则 blame，不得将写入文档或 memory 作为问题处理闭环。

## 3. 执行流程：Try Exactly Once Before Blame → Call → Wait → Continue/Exit

### Try Exactly Once Before Blame

- Agent 发现问题后，应执行一次必要、低风险、可解释的尝试，用于验证问题真实性并补足上报上下文。
- 若首次尝试解决问题，且无残留风险或不确定性，则不触发 blame，继续正常任务。
- 若首次尝试失败、只部分解决、引入新不确定性，或证明该问题需要外部协作，则立即进入 Call。
- 一旦进入 blame 上报路径，上报前必须完成恰好一次尝试；不得零尝试上报，也不得以第二次尝试延迟 blame。
- `### Tried` 必须记录该次尝试及结果；确实无法尝试时，应写明无法尝试的原因。

### Call

**前置：获取 appId 与 planId**

blame 路径形如 `/<appId>/<planId>`，二者**必须**真实有效，不得编造占位符。取值来源（按优先级）：

1. **首选**：读取仓库内 `.context/current-plan-metadata.yaml`，使用其中的 `appId` 与 `planId` 字段。该文件为 READONLY METADATA，由上游系统注入。
2. **兜底接收**：若 `.context/current-plan-metadata.yaml` 不存在 / 字段缺失 / 调用真实路径失败（详见下方 §3.x「兜底接收 `/blame/PLAN-000`」），**统一改投兜底接收**，由后台人工分发处理。

**命令模板**：

```bash
ai blame /<appId>/<planId> --title "<问题标题>" --context "<结构化上下文>"
```

示例：

```bash
ai blame /cms-mgr/PLAN-102 --title "Upload Image 接口字段命名应采用 snake_case 还是 camelCase" --context "### Category
决策歧义

### Where
cms-mgr/PLAN-102 / api/upload_image.go:42

### Observed
现有接口同时存在 snake_case（image_url）与 camelCase（imageUrl）两种命名，Upload Image 是该 PLAN 新增的第一个端点，无既定先例可循。

### Expected
PLAN-102 或 cms-mgr 仓库内应有统一的 JSON 字段命名规范，便于直接套用。

### Tried
1) grep cms-mgr 内 *.go 发现两种风格共存，比例约 6:4；
2) 查 docs/api-style.md 未命中；
3) 询问 codex agent，无明确答复。

### Ask
- PLAN-102 的对外 JSON 字段应采用 snake_case 还是 camelCase？
- 是否需要同时落地一份 docs/api-style.md 作为后续基准？

### Impact
首次尝试后仍存在决策歧义；blame 后 Agent 将根据裁决风险选择暂停或继续低风险推进。

### Refs
N/A"
```

字段之间使用**真实换行**分隔（双引号内直接回车）。**不要**在双引号内写字面 `\n` —— shell 不会把它展开为换行，服务端将原样存为反斜杠+n 两个字符。

- 路径参数 `/<appId>/<planId>`：blame 必须挂在某个 Plan 下，路径形如 `/cms-mgr/PLAN-102`。`appId` 与 `planId` 必须从 `.context/current-plan-metadata.yaml` 读取，缺失时按上方「前置」兜底处理，**不得编造**。
- `--title`：一行普通文字短句，能让人快速识别问题。**不得带任何前缀标签**（如 `### Category`、`[BLAME]`、`[BUG]` 等），也不得使用 Markdown 格式。分类信息放在 `--context` 的 `### Category` 段，不在 title 重复。
- `--context`：Markdown 文本，使用三级标题分段，见第 4 节模板。

### Wait

- 命令是**同步**的；必须等待命令返回。
- **只有**返回 ok **且**拿到 blame id，blame 才算成功。
- 命令失败 → blame 未成功，不得在对话中谎称已 blame；**必须按下方 §3.x 改投兜底接收**。

### Fallback：兜底接收 `/blame/PLAN-000`

后台保证 `/blame/PLAN-000` 这个兜底路径**始终存在且可接收**，由人工分发后续处理。

**触发条件**（任一命中即走兜底）：
- `.context/current-plan-metadata.yaml` 不存在 / `appId` 或 `planId` 字段缺失；
- 对真实 `/<appId>/<planId>` 路径调用 `ai blame` 返回非 ok（例：`Error: GetUserByName`、路径不存在、网络失败等）。

**兜底调用规则**：

1. **路径固定**：`/blame/PLAN-000`。
2. **title**：使用原本要上报的标题，前置 `[FALLBACK]` 标签，例如 `[FALLBACK] Upload Image 接口字段命名应采用 snake_case 还是 camelCase`。
3. **context**：在标准模板基础上，**必须额外**包含两段：
   - `### OriginalTarget`：原本计划投递的 `/<appId>/<planId>`（即使是占位/缺失，也要写出实际尝试的值，或写 `N/A` + 原因）。
   - `### FallbackReason`：触发兜底的原因，原样附上命令 stderr / 报错文本，不做改写、不做总结。
4. 其余字段（Category / Where / Observed / Expected / Tried / Ask / Impact / Refs）保持原样填写。

**兜底示例**：

```bash
ai blame /blame/PLAN-000 --title "[FALLBACK] Upload Image 接口字段命名应采用 snake_case 还是 camelCase" --context "### Category
决策歧义

### Where
cms-mgr/PLAN-102 / api/upload_image.go:42

### Observed
（原 Observed 内容）

### Expected
（原 Expected 内容）

### Tried
（原 Tried 内容）

### Ask
（原 Ask 内容）

### Impact
首次尝试后仍存在问题；blame 后 Agent 根据风险选择暂停或继续低风险推进。

### Refs
N/A

### OriginalTarget
/cms-mgr/PLAN-102

### FallbackReason
原路径调用失败：
\$ ai blame /cms-mgr/PLAN-102 --title ... --context ...
Error: GetUserByName"
```

**兜底也失败的情况**：若兜底 `/blame/PLAN-000` 同样返回错误，**不得再次重试**、**不得编造 blame id**。Agent 应在最终回复中如实告知「兜底接收亦失败」，并附上完整命令与报错，然后退出。

### Continue/Exit

- blame 成功（含兜底成功）后，Agent 不等待人类回复。
- Agent 应在回复中声明 blame id，简洁说明已上报的问题内容，供人类快速了解。
- 后续动作由 Agent 根据风险级别、上下文完整性与任务目标判断：可以继续低风险尝试、降级推进、暂停等待或退出任务。
- 若继续推进，不得围绕同一个未解决问题反复试错；遇到新的独立问题时，重新按「首次尝试 → 上报」流程处理。

## 4. `--context` 结构化模板

后端以 Markdown 格式存储 context，模板固定使用 Markdown **三级标题** 作为字段分段，缺失项内容写 `N/A`，**不得整段省略**；Agent 之间据此一致解析。

各字段之间用**真实换行**分隔（在 shell 双引号内直接回车），传给 `--context`。**不要**使用字面 `\n`——shell 不会展开，服务端将原样存为反斜杠+n 两个字符，破坏 Markdown 渲染。

```markdown
### Category
<信息缺失 | 决策歧义 | 执行受阻 | 风险越界 | 产物存疑>

### Where
<文件路径:行号 / 命令 / 任务阶段 / URL>

### Observed
<实际观察到的现象，事实描述，不加推测>

### Expected
<按上下文/规范本应是什么样>

### Tried
<Agent 首次尝试的动作与结果；若无法尝试写「未尝试，原因：...」>

### Ask
<希望协助/求助/询问的具体问题，一条以上用 - 列出>

### Impact
<若不解决会对任务造成什么影响；以及 blame 后 Agent 计划继续、暂停还是退出>

### Refs
<相关 blame id / commit / PR / 文档链接；无则 N/A>
```

**兜底接收附加字段（仅 `/blame/PLAN-000` 调用时使用）**：

```markdown
### OriginalTarget
<原本计划投递的 /<appId>/<planId>；若 metadata 缺失，写 N/A 并在 FallbackReason 说明>

### FallbackReason
<触发兜底的原因；如为命令失败，原样附上完整命令与 stderr 文本，不改写、不总结>
```

字段约束：

- 三级标题必须为英文且首字母大写（`### Category` / `### Where` / `### Observed` / `### Expected` / `### Tried` / `### Ask` / `### Impact` / `### Refs`，兜底场景额外 `### OriginalTarget` / `### FallbackReason`），顺序不得调换。
- 标题与正文之间换行；段与段之间留一个空行，便于后端 Markdown 渲染。
- `### Observed` 只写事实，不写解释。
- `### Ask` 必须是**可回答**的问题或**可执行**的请求，不得是「请看看」式开放表述。
- `### Impact` 必须说明不解决的影响，以及 blame 后 Agent 计划继续、暂停还是退出；如确实只是提示性上报（不影响主流程），需在 Impact 中明确写出。

## 5. 对话中引用规则

- blame 成功后获得 blame id，形如 `BLAME-005`（工具返回的完整字符串，包含 `BLAME-` 前缀与序号）。
- 对话中再次提及该问题时，**首次提及**使用 `` `[BLAME-005] <问题简述>` `` 格式标注（直接用完整 id，不再外加 `#`）。
- 重复提及同一问题时可简写为 `BLAME-005`。
- 同一问题在同一任务内**不重复 blame**；如已 blame 过，引用旧 id 即可。
- 跨任务/跨会话再次遇到相同问题，应先尝试检索旧 id，找不到再新建。

## 6. 反模式（禁止行为）

- ❌ 命中判定范围且首次尝试后仍未解决，却**不 blame**，自行假设结论并继续执行。
- ❌ 进入上报路径后零尝试上报，导致 `### Tried` 缺少可审计上下文。
- ❌ 上报前围绕同一问题尝试超过一次，形成非受控重试。
- ❌ blame 后未评估风险级别和任务上下文，直接无条件退出或无条件继续。
- ❌ blame 后必须等待用户回复才决定后续动作。
- ❌ 命令未返回 ok 就在对话中声称 blame 成功。
- ❌ 没有 blame id 就在对话中编造 `[BLAME-...]` 引用。
- ❌ 把多个独立问题塞进同一次 blame。
- ❌ 把问题**记录**到 AGENTS.md / memory / 代码注释 / TODO，当作「已处理」。记录 ≠ 解决，blame 才是解决通道。
- ❌ 反复重试同一失败动作而不 blame，违反 proactive 原则。
- ❌ 首次尝试失败后仍通过第二次尝试、猜测性结论或部分实现延迟上报。
- ❌ 处于边界场景或决策不确定状态，首次尝试后仍选择「不 blame」。
- ❌ 以上报量成本为由抑制必要 blame（试运行阶段：后台有人工跟进，鼓励 blame）。
- ❌ 引用上一轮会话或历史记忆中的旧规则覆盖本文件当前版本。
- ❌ 真实路径调用失败后**不改投**兜底 `/blame/PLAN-000`，直接沉默或继续任务。
- ❌ 兜底调用时遗漏 `### OriginalTarget` 或 `### FallbackReason`，导致后台无法定位原始上下文。
- ❌ 兜底 context 的 `### FallbackReason` 对报错文本做「总结/美化」，而非原样附上 stderr。

## 7. 与 AGENTS.md「回避状态」检测的边界

当 Agent 出现敷衍模式（频繁出现「完成了」「明白了」「对不起」等极简回复）时，往往伴随未上报的问题；应回头检查是否漏 blame。

## 8. Skill Type

**Rigid**：判定范围一旦命中，必须按四段（Try Exactly Once Before Blame → Call → Wait → Continue/Exit）执行，不得弱化、合并或跳过。模板字段不得自创或删减。试运行阶段的「首次尝试后上报 / 从严判定 / 上报前恰好一次尝试」为强约束，Agent 不得自行放宽。
