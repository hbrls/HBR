# Install Superpowers

1. Delete the old local Superpowers skill directories. Do not use `*`; list every path explicitly.

```bash
rm -rf .agents/skills/brainstorming
rm -rf .agents/skills/dispatching-parallel-agents
rm -rf .agents/skills/executing-plans
rm -rf .agents/skills/finishing-a-development-branch
rm -rf .agents/skills/receiving-code-review
rm -rf .agents/skills/requesting-code-review
rm -rf .agents/skills/subagent-driven-development
rm -rf .agents/skills/systematic-debugging
rm -rf .agents/skills/test-driven-development
rm -rf .agents/skills/using-superpowers
rm -rf .agents/skills/verification-before-completion
rm -rf .agents/skills/writing-plans
```

2. Copy the new Superpowers skill directories into `.agents/skills`.

```bash
mkdir -p .agents/skills
cp -R .context/HBR/.agents/skills/brainstorming .agents/skills/brainstorming
cp -R .context/HBR/.agents/skills/dispatching-parallel-agents .agents/skills/dispatching-parallel-agents
cp -R .context/HBR/.agents/skills/executing-plans .agents/skills/executing-plans
cp -R .context/HBR/.agents/skills/finishing-a-development-branch .agents/skills/finishing-a-development-branch
cp -R .context/HBR/.agents/skills/receiving-code-review .agents/skills/receiving-code-review
cp -R .context/HBR/.agents/skills/requesting-code-review .agents/skills/requesting-code-review
cp -R .context/HBR/.agents/skills/subagent-driven-development .agents/skills/subagent-driven-development
cp -R .context/HBR/.agents/skills/systematic-debugging .agents/skills/systematic-debugging
cp -R .context/HBR/.agents/skills/test-driven-development .agents/skills/test-driven-development
cp -R .context/HBR/.agents/skills/using-superpowers .agents/skills/using-superpowers
cp -R .context/HBR/.agents/skills/verification-before-completion .agents/skills/verification-before-completion
cp -R .context/HBR/.agents/skills/writing-plans .agents/skills/writing-plans
```

3. Wait for human review.
