# Design Template

**Owner**: Main Window

**Purpose**: Record design choice and rationale when implementation is not obvious.

**When to write**:
- Implementation is not obvious

---

## Writing Rules

- Create at `.claude/skills/pensieve/loop/{date}-{slug}/design.md`
- Do necessary code exploration during design
- After user confirmation, record the path in `_context.md`

## Must Include

- **Context references**: sources that justify the design
- **Approach overview**: what we will do
- **Decision rationale**: why this approach
- **Alternatives**: options considered but rejected
- **Risks**: what could go wrong

---

## Document Format

```markdown
# Design

## Context References

| Type | Path | Key Excerpt |
|------|------|----------|
| requirement | `.claude/skills/pensieve/loop/{date}-{slug}/requirements.md` | Core problem: ... |
| decision | `decisions/2026-01-20-xxx.md` | "Quoted text" |
| maxim | `maxims/linus.md` | "Good taste" |
| conversation | _context.md#requirements-clarification | User said: "..." |
| skill | `skills/xxx/SKILL.md` | Relevant capability |

## Approach Overview

[High‑level description; avoid too much detail]

## 上下文链接（recommended）
- 基于：[[前置决策或知识]]
- 导致：[[后续决策或流程]]
- 相关：[[相关主题]]

## Decision Rationale

| Decision Point | Choice | Rationale | Evidence |
|--------|------|------|------|
| [Issue 1] | [Choice] | [Why] | Context reference |
| [Issue 2] | [Choice] | [Why] | Context reference |

## Alternatives

### Option A: [Name]
- Description: [How]
- Rejection reason: [Why not]

### Option B: [Name]
- Description: [How]
- Rejection reason: [Why not]

## Risks

| Risk | Impact | Mitigation |
|------|------|------|
| [Risk 1] | [Impact] | [Mitigation] |
| [Risk 2] | [Impact] | [Mitigation] |
```

---

## Context Logging

**Must be recorded in `_context.md`**:
- Key design decisions and discussion
- User feedback on options
- Design changes and reasons

**Example**:
```markdown
## Key Conversation

### Design Decision (2026-01-23)

Discussion: state synchronization

Options:
1. Main-window continuation — read pending tasks and dispatch subagents one by one
2. .active file — simple JSON file for local state only

User preference: Option 2, because "simpler is better"

Final choice: .active file
```

---

## Common Pitfalls

| Pitfall | Response |
|---------|----------|
| Only describe the approach, not the rationale | "Why" matters more than "how" |
| No alternatives | Consider at least one alternative |
| Ignoring risks | Be honest about what can go wrong |
| Too detailed | This is design, not implementation — stay high‑level |
| Detached from requirements | Every decision should trace back to requirements |