---
description: 提交或复盘时自动将可复用结论沉淀到 knowledge/decision/maxim/pipeline，直接写入用户数据。
---

# 自改进（Auto Self-Improve）

> 工具边界见 `<SYSTEM_SKILL_ROOT>/references/tool-boundaries.md` | 共享规则见 `<SYSTEM_SKILL_ROOT>/references/shared-rules.md`

将经验与模式写入 Pensieve 四类用户数据：`maxim / decision / pipeline / knowledge`。

## Tool Contract

### Use when
- 提交 pipeline（`run-when-committing.md`）调用（自动触发）
- loop 完成后的收尾沉淀
- 用户明确要求"沉淀、记录、复盘、规范化"

### Failure fallback
- 发现结构性问题（旧路径并行/目录缺失/格式大面积不符）：跳过写入，建议运行 `doctor`
- 无法一次判断分类：按三层拆分（IS → `knowledge`，WANT → `decision`，MUST → `maxim`）

---

## 语义分层分配

先判定语义层，再决定写入类型：
1. **knowledge（IS）**：系统事实、现状、边界、机制
2. **decision（WANT）**：项目偏好、策略取舍、目标方向
3. **maxim（MUST）**：跨项目的硬约束与底线

同一洞察可拆分写入多层。`pipeline` 仅表达 HOW（执行顺序与验证闭环），不替代 IS/WANT/MUST。

### Pipeline 新建条件
- 同类任务在多个会话/loop 中重复出现
- 执行顺序会显著影响结果（步骤不可随意交换）
- 每步都有可验证完成标准

---

## 核心原则
- **自动沉淀**：由 pipeline 触发时直接执行
- **先读后写**：创建任何文件前先读对应 README
- **分类稳定**：只使用 `maxim / decision / pipeline / knowledge`
- **结论优先**：标题与第一句话能独立表达结论
- **准则单文件**：每条 `maxim` 独立文件
- **流程只做编排**：`pipeline` 只保留 task 编排与验证；理论背景外链
- **目标是减探索**：沉淀应让"下一次从症状到定位"的路径更短

---

## 定位加速知识模型

当问题属于"需要探索代码库才能回答"且内容属于 IS（事实层）时，`knowledge` 应覆盖：
1. **状态转换**：某动作触发后，数据和视图如何变化
2. **症状 → 根因 → 定位**：看到什么现象，去哪里查，为什么
3. **边界与所有权**：谁有写权限、谁只能调用、跨模块如何流转
4. **不存在/已移除**：哪些能力不在当前系统，避免重复误判
5. **反模式与禁区**：看起来可行但会失败的路径及原因

### 定位加速检查项（探索型问题适用）
- 至少 1 条"症状 → 根因 → 定位"映射
- 至少 1 条"边界与所有权"约束
- 至少 1 条"反模式/不要做什么"
- 给出可验证信号（日志、测试、运行结果、可观测行为之一）

---

## Phase 1: 提取与分类

**Goal**: 从会话上下文 + diff 中提取洞察并确定分类。

**Actions**:
1. 读取会话上下文与 `git diff --cached`
2. 提取核心洞察（可以是多条），判定语义层（IS/WANT/MUST），多层时拆分写入
3. 判断是否需要 `pipeline`（HOW）：仅当满足条件时新建
4. 路径规则：
   - `maxim`：`.claude/skills/pensieve/maxims/{one-sentence-conclusion}.md`
   - `decision`：`.claude/skills/pensieve/decisions/{date}-{conclusion}.md`
   - `pipeline`：`.claude/skills/pensieve/pipelines/run-when-*.md`
   - `knowledge`：`.claude/skills/pensieve/knowledge/{name}/content.md`

---

## Phase 2: 读取规范 + 写入

**Goal**: 按规范写入并保持链接连通。

**Actions**:
1. 读取目标 README（`<SYSTEM_SKILL_ROOT>/{type}/README.md`），按其格式生成内容：结论式标题 + 一句话结论 + 核心内容 + 语义链接（`<SYSTEM_SKILL_ROOT>/references/shared-rules.md` § 语义链接规则）
2. 类型特定要求：
   - `decision`：包含定位加速三项（下次少问什么 / 下次少查什么 / 失效条件）
   - 探索型 `knowledge`：包含定位加速检查项
   - `pipeline`：自检"不影响 task 编排的段落是否已外链"
3. 写入目标路径，在关联文档补反向链接（若需要）
4. 维护项目级 SKILL：
   ```bash
   bash <SYSTEM_SKILL_ROOT>/tools/project-skill/scripts/maintain-project-skill.sh --event self-improve --note "auto-improve: {file1,file2,...}"
   ```
5. 输出简短摘要（写入路径 + 沉淀类型）