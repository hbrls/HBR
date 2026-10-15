# Install OpenSpec

1. Delete the old local OpenSpec skill directories. Do not use `*`; list every path explicitly.

```bash
rm -rf .agents/skills/openspec-apply-change
rm -rf .agents/skills/openspec-archive-change
rm -rf .agents/skills/openspec-explore
rm -rf .agents/skills/openspec-propose
rm -rf .agents/skills/openspec-sync-specs
```

2. Copy the new OpenSpec skill directories into `.agents/skills`.

```bash
mkdir -p .agents/skills
cp -R .context/HBR/.agents/skills/openspec-apply-change .agents/skills/openspec-apply-change
cp -R .context/HBR/.agents/skills/openspec-archive-change .agents/skills/openspec-archive-change
cp -R .context/HBR/.agents/skills/openspec-explore .agents/skills/openspec-explore
cp -R .context/HBR/.agents/skills/openspec-propose .agents/skills/openspec-propose
cp -R .context/HBR/.agents/skills/openspec-sync-specs .agents/skills/openspec-sync-specs
```

3. Wait for human review.
