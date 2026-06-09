# Requirements

{本需求用于描述产品“要达成什么”以及“为什么要做”，为后续 Specs/Tasks 定基调；应尽量避免涉及技术实现细节，重点阐明产品目标与商业目标（例如用户价值、业务指标、增长/收入/成本、风险与合规影响）；同时说明为什么要做、本次变更解决什么问题、以及当前现状是什么。}

## Goals

{你要达成什么目标？请写清“可衡量的结果”，并尽量给出可验证的口径（例如数量、范围、完成定义、验收方式）。避免写成纯愿景或抽象口号。}

- {goal 1}
- {goal 2}

## Non-Goals

{你明确“不做什么”？请列出本阶段不包含的目标与边界，用于避免范围蔓延与误解；如未来可能做，也应写清“现在不做”的理由与触发条件（可选）。}

- {non-goal 1}
- {non-goal 2}

## Scope

{本次变更“包含什么”？请列出本阶段要交付的范围与边界（功能/场景/对象/平台/数据范围等），尽量可验证；与 Non-Goals 形成互补，避免遗漏与误解。}

- {scope item 1}
- {scope item 2}

## Non-Scope

{本次变更“不包含什么”？请列出明确不在本次交付范围内、但可能在后续迭代再做的事项；必要时说明推迟原因与进入后续的触发条件（可选）。}

- {out of scope item 1}
- {out of scope item 2}

## Functional Requirements

<!-- // EARS（Easy Approach to Requirements Syntax）是一种用于编写清晰、可验证功能需求的句式模板；参考：`what-is-ears-format.md`。 -->

{FR（Functional Requirements，功能需求）用于描述“系统在特定触发/状态下必须做什么或不得做什么”。每条 FR 应尽量可验证、可复现，避免把实现细节（类/函数/库/接口形态）写进需求本身。}

{以下示例需要保持“固定句式 + 占位符”风格，以确保不同需求文档输出口径一致；建议统一使用“系统应 / 系统不得”作为规范动词，避免混用多种表达。}

### 常规（Ubiquitous）需求

- **FR-001**: 系统应 {动作}。
- **FR-002**: 系统应 {动作}。

### 事件驱动（Event-Driven）需求

- **FR-010**: 当 {触发事件} 时，系统应 {动作}。
- **FR-011**: 当 {触发事件} 时，系统应 {动作}。

### 状态驱动（State-Driven）需求

- **FR-020**: 在 {状态条件} 期间，系统应 {动作}。
- **FR-021**: 在 {状态条件} 期间，系统应 {动作}。

### 非期望行为（Unwanted Behavior）需求

- **FR-030**: 如果 {错误条件}，则系统不得 {禁止动作}。
- **FR-031**: 如果 {错误条件}，则系统不得 {禁止动作}。

## Success Metrics

{如何衡量成功？请给出可对照的指标与验收方法（例如通过/不通过、数量、覆盖范围、证据形式）；如当前值未知可写 N/A，但必须写清“如何测量”。}

{默认不写性能指标（如延迟/吞吐/资源占用等），除非该需求与性能强相关或明确要求；避免在 Requirements 阶段引入不必要的性能口径。}

| Metric      | Current | Target | How to Measure |
|-------------|---------|--------|----------------|
| {metric 1}  | {value} | {value} | {method} |

## Dependencies

- **D-001**: {dependency 1}
- **D-002**: {dependency 2}

## Constraints

- **C-001**: {constraint description}
- **C-002**: {constraint description}

## Assumptions

- **A-001**: {assumption description}
- **A-002**: {assumption description}

## References

- **REF-001**: {reference 1}
- **REF-002**: {reference 2}