# Code Structure

> updated_by: Windsurf - Cascade
> updated_at: 2026-02-21 14:36:00

## Project Layout

The RE10 project is structured around defining workflows for an AI agent. The layout reflects this focus on configuration and instruction over traditional source code.

```
RE10/
├── .agents/commands/              # AI command definitions
│   ├── [command-name]/     # Directory for a specific command
│   │   └── COMMAND.md      # The main definition file for the command
│   └── ...
├── .skills/                # Reusable AI skill definitions
│   ├── [skill-name]/
│   │   └── SKILL.md
│   └── ...
├── .steering/              # Project governance and high-level direction
│   ├── constitution.md     # Core principles and decision-making
│   ├── product.md          # Product vision and goals
│   ├── design.md           # Design principles and architecture
│   └── structure.md        # This file: code structure and conventions
├── .windsurf/              # Supporting configurations for the agent
│   └── workflows/
├── AGENTS.md               # Core instructions for the AI agent
└── README.md               # Project overview and setup (to be created)
```

## Naming Conventions

Consistency in naming is crucial for the AI agent to correctly parse the project.

| Type | Convention | Example |
|---|---|---|
| Directories | `kebab-case` | `.agents/commands/task-dispatcher` |
| Core Files | `UPPERCASE.md` | `COMMAND.md`, `SKILL.md`, `AGENTS.md` |
| Other Files | `kebab-case.md` | `references/constitution.md` |

## Code Patterns

As this project primarily consists of Markdown-based definitions, traditional code patterns do not apply. Instead, the structure within the `.md` files themselves follows a pattern:

- **Frontmatter**: YAML frontmatter is used for metadata (`name`, `description`, `version`).
- **Headings**: Markdown headings (`#`, `##`) are used to structure the content logically.
- **Clarity**: Content is written to be unambiguous for both human and AI interpretation.

## Changelog

<!-- // 这是一个 Living Document，如无必要，无需维护变更历史。 -->
