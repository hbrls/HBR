# Loop (Execution Layer)

Combines the Claude Code Task system with a local tracking directory for main-window controlled loop execution.

## Role Division

| Role | Responsibility |
|------|----------------|
| **Main Window** | Planning: init → fill context → generate tasks → call task-executor |
| **task-executor** | Execute tasks: read context → load knowledge as needed → execute → capture learnings if needed |
| **Main Window (Continuation)** | Continue loop by checking pending tasks and dispatching subagents |

## Startup Flow (Main Window)

### Step 1: Initialize loop directory (prepare-only)

```bash
<SYSTEM_SKILL_ROOT>/tools/loop/scripts/init-loop.sh <slug>
# Example:
<SYSTEM_SKILL_ROOT>/tools/loop/scripts/init-loop.sh login-feature
```

> Note: init-loop.sh runs quickly. Run it in the foreground so you can read the `LOOP_DIR` output.

### Step 2: Fill context (Main Window)

In the loop directory (`.claude/skills/pensieve/loop/{date}-{slug}/`):

1. **Create and fill `_context.md`** (see format below; to avoid "Read before Write" friction, init-loop.sh no longer creates a template file)
2. **Create documents as needed**
   - `requirements.md` — only when requirements are unclear (goal/scope/constraints not confirmed)
   - `design.md` — only when implementation is not obvious
   - `plan.md` — code exploration notes (need code understanding)

### _context.md format

```markdown
# Conversation Context

## Pre-Context

### Interaction History
[Record the conversation before entering loop]

| Turn | Model Attempt | User Feedback |
|------|--------------|---------------|
| 1 | Proposed plan A | Rejected, asked simpler |
| 2 | Proposed plan B | Approved, entered loop |

### Final Consensus
[Shared understanding at loop entry]
- Goal: XXX
- Scope: YYY
- Constraints: ZZZ

### Understanding & Assumptions
[Model’s expectations]
- Expected modules involved
- Expected implementation approach
- Expected difficulties

### Document References
| Type | Path |
|------|------|
| requirements | none / path |
| design | none / path |
| plan | none / path |

### Context Links (optional)
- 基于：[[前置决策或知识]]
- 导致：[[后续决策、流程或文档]]
- 相关：[[相关主题]]

---

## Post-Context

> If execution matches the plan, leave blank or note "no deviation".

### Deviations
[Pre‑assumptions vs reality]
- Before: XXX
- Found: YYY
- Adjustment: ZZZ

### Interventions
[Manual interventions during execution]
```

### Step 3: Generate tasks (Main Window)

Before splitting tasks, read project maxims:

- All files under `<USER_DATA_ROOT>/maxims/` (`*.md`)

Then generate tasks based on context + maxims:

| Workload | Task count |
|----------|------------|
| A few lines | 1 |
| One module | 2–3 |
| Multiple modules | 4–6 |

Each task includes:
- subject (imperative, e.g., "Implement user login")
- description (source + action + completion criteria)
- activeForm (progressive, e.g., "Implementing user login")

After Phase 2 confirmation, create tasks directly in Claude Task system.
Do not use a standalone markdown/bullet list as a substitute for real task creation.
If you need to present tasks to the user, show a snapshot from created tasks (task id + subject).

Hard rule: do not generate task lists before maxims are loaded.

### Step 4: Continue from Main Window

After creating the first real task, the main window should continue by selecting the next pending task and dispatching one subagent at a time.

> Hooks are no longer required for loop continuation.

### Step 5: Execute tasks

Call an agent for each task:

```
Task agent=task-executor prompt="
task_id: 1
context: .claude/skills/pensieve/loop/{date}-{slug}/_context.md
system_skill_root: <SYSTEM_SKILL_ROOT>
user_data_root: .claude/skills/pensieve
"
```

The agent returns after one task. The main window checks pending tasks and continues with the next one.

## Two Storage Systems

| Storage | Content | Purpose |
|---------|---------|---------|
| `~/.claude/tasks/<uuid>/` | Task state (JSON) | Claude Code native |
| `.claude/skills/pensieve/loop/{date}-{slug}/` | Context + docs | Project‑level tracking and learnings (never overwritten) |

## Directory Structure

```
~/.claude/tasks/<uuid>/          # Claude Code Tasks (status)
    ├── 1.json
    ├── 2.json
    └── ...

.claude/skills/pensieve/loop/           # Project tracking (context + learnings)
    └── 2026-01-23-login/        # One directory per loop
        ├── _context.md          # Conversation context, interventions
        ├── requirements.md      # Requirements (if any)
        └── design.md            # Design notes (if any)
```

## Loop Continuation Mechanism

Main window continues until tasks complete:

```
Agent runs one task → returns
    ↓
Main window checks task state
    ↓
├── pending task exists → dispatch next subagent
└── all complete → wrap up
```

## Reinforcement Message

Injected on each continuation:

```markdown
## Loop Continue

**Pipeline**: develop
**Progress**: [2/5] completed
**Current task**: #3 Implement user login

---

## Task Content
{task description}

---

**Execution requirements**:
1. Complete the task
2. TaskUpdate → completed
3. If intervention occurs, record in _context.md
```

## Phase Selection

| Task characteristics | Phase combination |
|---------------------|-------------------|
| Clear, small scope | tasks |
| Need code understanding | plan → tasks |
| Need technical design | plan → design → tasks |
| Unclear requirements | plan → requirements → design → tasks |

## Closed-Loop Learning (Main Window)

After the agent returns, run self-improve:

```
Pre‑assumptions → execution → post‑deviations → capture learnings
```

### Flow

1. Read `tools/self-improve/_self-improve.md`
2. Compare `_context.md` pre/post sections
3. Fill Post‑Context (deviations)
4. If meaningful deviation exists, ask user to capture
5. Upon consent, write using README format
6. If captured as `decision`/`pipeline`, ensure at least one `基于/导致/相关` link is present

### Post-Context Example

| Pre‑assumption | Actual finding | Adjustment |
|---------------|----------------|------------|
| Two code paths are identical | RPWindow adds icon + styles | Add variant prop |
| 9 components independent | Helper components not reusable | Split into 3 window components |

### What counts as a deviation

**Meaningful** (worth capturing):
- Architecture assumption was wrong
- Edge cases missed
- Tool/framework limits discovered late

**Not meaningful** (do not capture):
- Typos, small adjustments
- One‑off special cases
