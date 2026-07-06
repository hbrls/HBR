---
name: install-openspec
metadata:
  version: 0.0.1
---

# Install OpenSpec

## Upgrade

1. Find `https://github.com/Fission-AI/openspec` and follow the install instructions.

```bash
npm install --save-dev @fission-ai/openspec@latest
npx openspec init --tools kilocode
```

2. Delete the old local OpenSpec skill directories. Do not use `*`; list every path explicitly.

```bash
rm -rf .agents/skills/openspec-apply-change
rm -rf .agents/skills/openspec-archive-change
rm -rf .agents/skills/openspec-explore
rm -rf .agents/skills/openspec-propose
rm -rf .agents/skills/openspec-sync-specs
```

3. Copy the new OpenSpec skill directories into `.agents/skills`.

```bash
mkdir -p .agents/skills
cp -R .kilocode/skills/openspec-apply-change .agents/skills/openspec-apply-change
cp -R .kilocode/skills/openspec-archive-change .agents/skills/openspec-archive-change
cp -R .kilocode/skills/openspec-explore .agents/skills/openspec-explore
cp -R .kilocode/skills/openspec-propose .agents/skills/openspec-propose
cp -R .kilocode/skills/openspec-sync-specs .agents/skills/openspec-sync-specs
```

4. Wait for human review.

## Tuning

Only tune the following installed OpenSpec skills:

- `openspec-apply-change`
- `openspec-archive-change`
- `openspec-explore`
- `openspec-propose`
- `openspec-sync-specs`

All tuning must stay within these specified project-local skills under `.agents/skills`.

<!-- TODO: tuning 策略待补充 -->

4. Wait for human review.