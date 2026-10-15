---
name: sk-approve-phase
metadata:
  version: 0.0.3
---

## Approve Phase

**FAIL-FAST REQUIREMENT** Check if the tools exist, stop immediately and exit on any error.

```bash
$ backstage-gitea plan --help
```

## Initialization

**MUST** 读取相关的 `Plan` 信息在 `.context/current-plan-metadata.yaml` 和 `.context/current-plan.md`。读取失败时必须立刻停止并退出。

**MUST** 读取相关的 `Phase` 信息在 `.context/current-phase-metadata.yaml` 和 `.context/current-phase.md`。读取失败时必须立刻停止并退出。

## Workflow

- Step 1: 如果 `.context/current-phase.md` 有改动，使用 `backstage-gitea` 工具将 `.context/current-phase.md` 上传更新到 Gitea 项目中。

```bash
$ backstage-gitea agent PUT /:appId/:planId/:phaseId --status PASS --context .context/current-phase.md
```

- Step 2: 如果 `.context/current-plan.md` 有改动，使用 `backstage-gitea` 工具将 `.context/current-plan.md` 上传更新到 Gitea 项目中。

```bash
$ backstage-gitea agent PUT /:appId/:planId --context .context/current-plan.md
```
