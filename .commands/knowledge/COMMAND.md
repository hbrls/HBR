---
id: knowledge-readme
type: knowledge
title: Knowledge 规范
status: active
created: 2026-02-28
updated: 2026-02-28
tags: [pensieve, knowledge, spec]
---

# Knowledge（知识）

外部参考材料：技术文档、API 说明、最佳实践等。

## 目的

Knowledge 的核心价值是降低执行摩擦。

如果没有这条知识，模型会卡在哪里，代价有多大？

### 三类摩擦

| 摩擦类型 | 特征 | 例子 |
|---|---|---|
| 时间差 | 知识晚于模型训练截止 | Next.js 15 API、Claude Code 新功能 |
| 隐性知识 | 需要从结构中推断 | 命名约定、架构取舍 |
| 分散知识 | 信息存在但检索成本高 | issues、邮件、源码注释 |

## 捕获标准

核心问题：**不写下来会带来什么执行摩擦？**

| 摩擦等级 | 动作 |
|---|---|
| 高频阻塞且恢复成本高 | 必须沉淀 |
| 偶发阻塞且可快速搜索 | 不沉淀，只保留链接 |
| 一次性问题 | 不沉淀 |

### 适合沉淀的信号

| 信号 | 说明 |
|---|---|
| 模型反复问同一问题 | 时间差或隐性知识缺失 |
| 搜索结果不准/过时 | 时间差问题 |
| 每次都要从代码猜约定 | 隐性知识未显式化 |
| pipeline 依赖外部标准 | 应沉淀为可复用参考 |

## 语义定位（IS 层）

`knowledge` 只承载 IS（事实层）：系统现状、机制边界、可验证行为，即“是这样”。

当一次任务主要成本在探索与定位时，若内容属于 IS，建议沉淀：

1. 状态转换：动作触发后，数据/行为如何变化
2. 症状 -> 根因 -> 定位：看到什么现象，去哪里查，为什么
3. 边界与所有权：谁可以改、谁只能调用、跨模块如何流转
4. 不存在/已移除：哪些能力不在当前系统，避免重复探索
5. 反模式：看起来可行但会失败的路径

> 原则：先判定语义层，再落地类型。`knowledge`=IS，`decision`=WANT，`maxim`=MUST，`pipeline`=HOW。

## 关系与演化

| 方向 | 说明 |
|---|---|
| Knowledge -> Decision | 外部知识 + 项目实践形成决策 |
| Knowledge -> Pipeline | 外部标准约束执行流程 |
| Knowledge -> Maxim | 最佳实践被内化为准则 |

### Knowledge 与 Decision 的区别

| 类型 | 本质 | 判断句 |
|---|---|---|
| Knowledge | 外部输入 | “世界是这样运作的” |
| Decision | 内部输出 | “我们决定这样做” |

边界例子：
- 框架限制“有文档但难找” -> Knowledge（分散）
- 没文档、踩坑得出 -> Decision（内部经验）

## 编写规范

### 目录结构

```
.claude/skills/pensieve/knowledge/{name}/
├── content.md      # 知识正文
└── source/         # 来源文件（可选）
```

### 文件格式

```markdown
# {知识标题}

## Source
[原始链接或参考来源]

## Summary
[一句话摘要]

## Content
[正文：摘录或综合]

## When to Use
[什么场景下应查阅这条知识]

## 上下文链接（推荐）
- 基于：[[前置知识或决策]]
- 导致：[[会影响的决策或流程]]
- 相关：[[相关主题]]
```

### 探索型知识模板（推荐）

```markdown
# {主题}

## Source
[来源：代码库/文档/会话]

## Summary
[一句话说明这条知识能节省什么探索成本]

## 状态转换
[动作 -> 状态变化 -> 可观察结果]

## 症状 -> 根因 -> 定位
- **[症状]**: [根因] -> [文件/模块/入口]

## 边界与所有权
- [模块 A 负责什么]
- [模块 B 只读/只调，不直接写]

## 反模式（Do Not）
- [不要这样做 + 原因]

## 验证信号
- [日志/测试/运行行为]

## When to Use
[下一次遇到什么情况先读这条]
```

### 示例

```markdown
# Agent Design Best Practices

## Source
https://www.anthropic.com/engineering/advanced-tool-use

## Summary
Anthropic 官方的 agent 工具设计指南。

## Content
- Tool Search Tool：动态发现工具，减少 token 消耗
- Programmatic Tool Calling：用代码编排工具调用，减少上下文噪声
- Tool Use Examples：通过例子提高参数调用准确率

## When to Use
设计 agent、优化工具调用、降低 token 消耗时。
```

## 备注

- Knowledge 是输入，不是输出
- 优先链接原始来源，避免无必要复制
- 需要本地副本时，用复制/移动命令，不要手抄重写
- 定期清理过期知识
- 链接是推荐项，仅在能提升检索价值时保留

## 系统知识 vs 项目知识

- 系统知识：`skills/pensieve/knowledge/`（随插件更新）
- 项目知识：`.claude/skills/pensieve/knowledge/`（永不覆盖）
