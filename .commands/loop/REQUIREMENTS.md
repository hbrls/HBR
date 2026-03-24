# Requirements Template

**Owner**: Main Window

**Purpose**: Clarify ambiguous requirements into a stable goal/scope/constraints baseline.

**When to write**:
- Requirements are unclear (goal/scope/constraints not confirmed)

---

## Writing Rules

- Create at `.claude/skills/pensieve/loop/{date}-{slug}/requirements.md`
- Write a first draft before asking questions
- After user confirmation, record the path in `_context.md`

---

## Document Format

```markdown
# Requirements Anchor

## Context References

| Type | Path | Key Excerpt |
|------|------|----------|
| decision | `decisions/xxx.md` | "Quoted text" |
| maxim | `maxims/linus.md` | "Good taste" |
| conversation | _context.md#requirements-clarification | User said: "..." |
| skill | `skills/xxx/SKILL.md` | Relevant capability |

## Core Problem

[One sentence: what problem are we solving?]

## 上下文链接（recommended）
- 基于：[[前置决策或知识]]
- 导致：[[后续决策或流程]]
- 相关：[[相关主题]]

## Success Criteria

- [ ] [Verifiable completion criterion 1]
- [ ] [Verifiable completion criterion 2]
- [ ] ...

## Boundaries

### In Scope
- [Clearly in scope]

### Out of Scope
- [Explicit exclusions to prevent scope creep]

## Constraints

- [Technical constraints]
- [Time constraints]
- [Resource constraints]
```

---

## Context Logging

**Must be recorded in `_context.md`**:
- Key conversations during requirement clarification
- User's original wording vs final understanding
- Any requirement changes and their reasons

**Example**:
```markdown
## Key Conversation

### Requirements Clarification (2026-01-23)

User: "I want multiple terminals to share state."

Follow‑ups:
Q: What scenario do you expect?
A: I started a task in terminal A and want to continue in terminal B.

Q: Do you need two terminals running different tasks at the same time?
A: No, only one at a time.

Q: What is the desired outcome?
A: Open a new terminal and continue where I left off.

Clarified experience:
- Start a loop in terminal A
- Close terminal A, continue in terminal B
- Progress is preserved from last stop
```

---

## Common Pitfalls

| Pitfall | Response |
|---------|----------|
| Requirements too vague | Success criteria must be verifiable |
| Unclear boundaries | "Out of scope" must be explicit |
| Forgot to update | Update the anchor when requirements change |
| Writing a PRD | This is an anchor, not a full doc — keep it short |