# Amuro `clarify` ViewTab：内存事件投影实施方案

> updated_by: HBR - DeepSeek-4 (DeepSeek Harness)
> updated_at: 2026-10-30 16:09:45

<!-- // 读者提示：本文件交实施 Agent 执行。目标仓库是 Amuro，不是 HBR。
     // 只实现 Client 侧 conversation.view 投影；Host 侧 clarify 工具已经实现并验证可调用，不得重写。
     // 简化前提沿用需求讨论：单活动 Session、用户在场、不做性能折叠、展示当前已加载窗口中的全部数据。 -->

## 1. 目标

在 `/Users/shuaizhexu/github/Amuro/packages/dsh-plugin-amuro` 中新增一个 Conversation ViewTab：

1. 监听共享 Session 事件流中的全部 `tool/call`；
2. 在事件投影入口过滤 `event.data.name === 'clarify'`；
3. 从 `arguments` 解析 `type` 与 `payload`，按调用顺序显示；
4. 派生列表只保存在内存 target 快照中，不写文件、不追加 Session 事件；
5. Tab 打开期间实时追加新记录；首次打开时重建当前已加载事件窗口里的历史记录。

## 2. 已确认的机制边界

### 2.1 没有 by-tool-name 监听接口

DSH 的表式 `match` 只按 `event.type` 预路由。所有工具调用共享 `tool/call`，因此必须保留：

```ts
event.data.name === 'clarify'
```

这层过滤。`toolCallFocus` 只负责按 `callId` 定位已有记录，不是订阅过滤器，本方案不使用它。

### 2.2 过滤必须放在 match

不要在 React 渲染时过滤全部工具调用。非 clarify 调用应在 `match` 返回 `null`，从而：

- 不创建 Context；
- 不进入 clarify target builder；
- 不触发无关快照更新或 React 重渲染。

### 2.3 只消费 tool/call

clarify 是 LLM → 用户的单向输出，展示内容已经完整存在于 `tool/call.data.arguments`。

本方案**不监听 `tool/result`**：其结果只是「已投递」确认，不是用户回答，也不是新的展示内容。

### 2.4 内存快照不等于原始事件不持久

原始 `tool/call` 仍由 DSH 正常写入 Session 事件日志；本任务不修改该行为。

ViewTab 自己生成的 `{ records: [...] }` 快照只在内存中存在。页面重载、Session 切换或 target 激活时，从共享事件窗口重新计算，不建立第二份持久数据。

## 3. 代码落点

现有 Client 入口为 `src/client/index.tsx`，目前只有 Sidebar/Main 注册。本任务新增 Conversation View 链路，建议拆分为：

| 文件 | 改动 |
|---|---|
| `src/client/index.tsx` | 扩展 Client inject；注册事件 Definition、target view、Session Hook 与 `conversation.view` Slot |
| `src/client/clarifyProjection.ts` | 新增：clarify 事件 Definition、记录类型、target builder 与空快照 |
| `src/client/ClarifyView.tsx` | 新增：纯展示组件，只消费快照，不监听 Session |
| `src/client/Amuro.module.css` 或独立样式文件 | 按需增加列表、空状态、类型与 payload 样式 |
| `src/index.js` | **不改**：Host clarify 工具保持现状 |
| `cordis.patch.yml` | **不改**：插件已经挂载 |

不要求新增 Cordis 配置，也不要把功能做进现有 `AmuroPage` Sidebar 页面；本需求的目标是 Conversation ViewTab。

## 4. 事件 Definition 契约

新增一个稳定 Definition，建议：

```text
kind:   amuro-clarify-call
target: amuro-clarify
```

### 4.1 match

使用表式 match，仅声明 `tool/call`：

- `name !== 'clarify'`：返回 `null`；
- `name === 'clarify'`：返回 `{ id: String(callId), role: 'start' }`。

`callId` 是记录主键；同一调用即使后续产生 result，也不会在本 Definition 中生成第二条记录。

### 4.2 start

从匹配事件中生成不可变记录状态：

```ts
interface ClarifyRecord {
  callId: string
  seq: number
  time: number
  type?: string
  payload?: string
  argumentsText: string
  parseError?: string
}
```

规则：

- `seq` 用于排序，`time` 用于显示，`callId` 用于 React key；
- 正常路径：`JSON.parse(arguments)` 后读取 `type` 与 `payload`；
- 解析失败或字段不是字符串时，不得抛错中断整个 Conversation assembler；保留原始 `argumentsText` 并标记 `parseError`，由 UI 显示为无效记录；
- 不读取 execute 返回值，不推断投递状态。

### 4.3 update

本 Definition 当前只认领 start 事件。`update` 按框架契约保留并返回现有 state，不合成额外业务状态。

### 4.4 buildViewNode

必须提供 `buildViewNode`；只声明 `target` 而不提供它，会在注册时被框架拒绝。

节点要求：

- `key` 使用 Context 的稳定 key；
- `target === 'amuro-clarify'`；
- `anchorSeq` 使用调用的 `seq`；
- `data` 保存完整 `ClarifyRecord`；
- state 为空时返回 `null`。

## 5. target builder 与快照

注册 View Definition：

```text
target: amuro-clarify
create: () => ClarifySnapshotBuilder
```

Builder 只需要支持框架要求的两个入口：

### 5.1 replace(input)

- 清空旧节点；
- 按 `node.key` 重建 Map；
- 返回完整新快照。

用于 target 首次激活、重新绑定或全量重建。

### 5.2 apply(input)

- 按 `node.key` upsert 本轮变化节点；
- 返回新快照。

每次工具调用只产生一个不可变节点，因此无需复杂的 splice、删除或分页算法。

### 5.3 snapshot 形状

```ts
interface ClarifySnapshot {
  records: ClarifyRecord[]
}
```

输出前按 `seq` 升序排列，必要时以 `callId` 作为稳定的次级排序键。快照与记录数组应保持不可变，避免 React 使用引用比较时漏更新。

空快照必须是稳定常量：

```ts
const EMPTY_CLARIFY_SNAPSHOT = { records: [] }
```

避免无 Session 或空列表时每次 `getSnapshot()` 都创建新对象。

## 6. Session Hook 与 ViewTab 注册

### 6.1 Client 服务依赖

在现有 Client `inject` 中增加：

```text
uiConversation
uiSession
```

保留现有 `slots`。事件与 target 分别通过：

```text
ctx.uiConversation.events.register(...)
ctx.uiConversation.views.register(...)
```

注册。

### 6.2 Hook

通过 `ctx.uiSession.provide('amuro-clarify', ...)` 暴露：

```ts
interface ClarifyViewProps {
  useClarify: () => ClarifySnapshot
}
```

Hook 内部：

1. 从 `useConversationContext()` 取得当前 Conversation；
2. 用 `useSyncExternalStore` 订阅 `conversation.views.target('amuro-clarify')`；
3. 无 Conversation 时返回稳定空快照；
4. 不在组件中自行创建事件订阅或保存记录副本。

`hasValue` 只判断当前 Conversation 是否存在；真实数据由 target source 提供。

### 6.3 conversation.view Slot

在现有 Sidebar/Main 注册之外，新增：

```text
name:  conversation.view
id:    plugin-amuro-clarify
label: Clarify
order: 置于现有 Conversation Tabs 之后即可
```

Slot 组件直接接收 Slot 传入的 `useClarify` 并渲染 `ClarifyView`。图标复用 Amuro 现有官方 primitives 图标方案即可，不新增图标依赖。

## 7. 展示要求

最小展示单元包括：

- `type`：作为记录类型标签；
- `payload`：作为主体文本，保留换行；
- `time`：可选的辅助时间信息；
- 空状态：当前窗口没有 clarify 调用时显示明确空态；
- 无效记录：存在 `parseError` 时显示原始 arguments 与错误标记，不让整个 Tab 崩溃。

首版不需要编辑、删除、已读状态、分页、搜索或持久化，也不需要显示「已投递」结果行。

## 8. 构建与验证

### 8.1 构建

在 Amuro 插件包执行现有构建流程：

```bash
pnpm bundle
```

确认 `lib/client.js` 已更新。Host 源码未改，不应产生 Host 行为变化。

Client 插件更新后，通过刷新当前 DSH Web GUI 验证；不要假设修改 `src/` 会自动生效，也不要启动替代 Web 服务。

### 8.2 功能验收序列

1. **历史重建**：在已有 clarify 调用的 Session 中首次打开 Clarify Tab，当前已加载窗口内的记录全部出现；
2. **名称过滤**：产生一次 clarify 调用与至少一次其他工具调用，Tab 只增加 clarify 记录；
3. **字段解析**：正常调用的 `type`、`payload` 与 ToolCall 参数一致；
4. **实时追加**：Tab 保持打开时再次调用 clarify，新记录按 `seq` 追加且旧记录不重复；
5. **无结果噪音**：不出现单独的「已投递」记录；
6. **内存重建**：刷新页面后记录能从 Session 事件窗口重建，无需插件自己的存储文件；
7. **Session 隔离**：切换 Session 后只显示该 Session 的 clarify 记录；
8. **异常容错**：构造无法解析的 clarify arguments 时，仅该记录显示异常，Conversation 与其他 Tab 不受影响。

### 8.3 建议的低成本自测

对纯函数/Builder 做最小自测：

- match 接受 clarify、拒绝其他工具；
- start 正确解析双字符串参数；
- replace 后 apply 不重复、不丢失；
- 多次 upsert 同一 callId 仍只有一条；
- 空快照引用稳定；
- 非法 JSON 不抛出到 assembler。

## 9. 明确不做

- 不实现 by-tool-name 框架扩展；保留 match 中的一行 name filter；
- 不监听 `tool/result`，不维护 `callId → name` 配对表；
- 不调用 `session.append`，不写 localStorage/IndexedDB/文件；
- 不把过滤逻辑放进 React；
- 不修改 clarify 工具的 description、schema、execute 或 output；
- 不实现 MCP、Codex、历史分页、性能折叠或跨 Session 聚合；
- 不改动 HBR 仓库中的任何插件代码。

## 10. 完成判定

当 Clarify ViewTab 能仅依据 `tool/call` 重建并实时展示当前 Session 的 clarify 参数列表，其他工具调用不会产生 Context 或快照更新，且插件自身没有任何派生数据持久化时，本任务完成。
