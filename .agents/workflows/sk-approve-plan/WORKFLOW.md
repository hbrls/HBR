---
name: sk-approve-plan
metadata:
  version: 0.0.3
---

## Approve Plan

**FAIL-FAST REQUIREMENT** Check if the tools exist, stop immediately and exit on any error.

```bash
$ backstage-gitea plan --help
```

## Initialization

**MUST** 读取相关的 `Plan` 信息在 `.context/current-plan-metadata.yaml` 和 `.context/current-plan.md`。读取失败时必须立刻停止并退出。

## Workflow

*目前是 Workflow 调试阶段，Phase 已经拆分好了。所以需要做的只是将 Phase 使用 backstage-gitea 工具上传。*

- Step 3: 使用 `backstage-gitea` 工具将 `.context/current-plan.md` 的 `Phases` 章节创建到 Gitea 项目中。每个 Phase 调用一次创建工具。

```bash
$ backstage-gitea plan POST /:appId/:planId/phases --name "Phase Name"
```

## Example

**输入数据：**

从 `.context/current-plan-metadata.yaml` 中读取 `appId` 和 `planId`：

```yaml
appId: cms-mgr
planId: PLAN-102
```

从 `.context/current-plan.md` 中读取 `Phases` 章节内容：

```markdown
## Phase

### Phase 1: 修复前导斜杠问题

### Phase 2: 修复尾部空格问题
```

**执行操作：**

提取每个 Phase 的标题，去除序号前缀（`Phase N:`），然后逐个调用 API 创建：

```bash
$ backstage-gitea plan POST /cms-mgr/PLAN-102/phases --name "修复前导斜杠问题"
# 工具返回 phaseId: PHASE-100

$ backstage-gitea plan POST /cms-mgr/PLAN-102/phases --name "修复尾部空格问题"
# 工具返回 phaseId: PHASE-200
```

**后处理：**
- 将工具自动生成的 `phaseId` 回填至 `.context/current-plan.md` 对应的 Phase 章节
- 注意：API 不支持批量创建功能，必须逐个 Phase 串行调用
