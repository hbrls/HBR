---
name: install-hbr
metadata:
  version: 0.0.1
---

# Install HBR

## Pre-Condition

The user has cloned `https://github.com/hbrls/HBR` into `.context/HBR`.

## Install

1. Install `AGENTS.md`.

```bash
rm -f AGENTS.md
cp .context/HBR/AGENTS.md AGENTS.md
```

Wait for human review.

2. Install skills: `blame` and `use`.

```bash
rm -rf .agents/skills/blame
rm -rf .agents/skills/use
mkdir -p .agents/skills
cp -R .context/HBR/.agents/skills/blame .agents/skills/blame
cp -R .context/HBR/.agents/skills/use .agents/skills/use
```

Wait for human review.

3. Install workflows. Select as needed:

```bash
rm -rf .agents/workflows/w-execute
rm -rf .agents/workflows/a-execute
rm -rf .agents/workflows/sk-explore-plan
rm -rf .agents/workflows/sk-enter-phase
rm -rf .agents/workflows/sk-explore-task
rm -rf .agents/workflows/sk-execute-task
rm -rf .agents/workflows/sk-evaluate-task
rm -rf .agents/workflows/sk-approve-phase
rm -rf .agents/workflows/sk-approve-plan
mkdir -p .agents/workflows
cp -R .context/HBR/.agents/workflows/w-execute .agents/workflows/w-execute
cp -R .context/HBR/.agents/workflows/a-execute .agents/workflows/a-execute
cp -R .context/HBR/.agents/workflows/sk-explore-plan .agents/workflows/sk-explore-plan
cp -R .context/HBR/.agents/workflows/sk-enter-phase .agents/workflows/sk-enter-phase
cp -R .context/HBR/.agents/workflows/sk-explore-task .agents/workflows/sk-explore-task
cp -R .context/HBR/.agents/workflows/sk-execute-task .agents/workflows/sk-execute-task
cp -R .context/HBR/.agents/workflows/sk-evaluate-task .agents/workflows/sk-evaluate-task
cp -R .context/HBR/.agents/workflows/sk-approve-phase .agents/workflows/sk-approve-phase
cp -R .context/HBR/.agents/workflows/sk-approve-plan .agents/workflows/sk-approve-plan
```

Wait for human review.

4. Ask the user whether to install Kilo configuration. If yes:

```bash
rm -f .kilo/rules/steering.md
mkdir -p .kilo/rules
cp .context/HBR/.kilo/rules/steering.md .kilo/rules/steering.md
```

Wait for human review.

5. Ask the user whether to install Codex configuration. If yes:

```bash
rm -f .codex/config.toml
mkdir -p .codex
cp .context/HBR/.codex/config.toml .codex/config.toml
```

Wait for human review.

## Install Superpowers

Ask the user whether to install Superpowers Skills. If yes:

See [references/install-superpowers.md](references/install-superpowers.md).

## Install OpenSpec

Ask the user whether to install OpenSpec Skills. If yes:

See [references/install-superpowers.md](references/install-openspec.md).
