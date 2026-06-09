---
name: sk-enter-phase
description: Execute Workflow: .agents/workflows/sk-enter-phase/WORKFLOW.md
metadata:
    version: 0.0.3
---

# Enter Phase

**FAIL-FAST REQUIREMENT** 检查 `# LLM Coding Level Board`，当前使用基模的评分必须优于或等于 `L1`。如果满足，回复 **{实际使用的基模名称} + 所以我来砍，阿剑也不会生气吧**。如果不满足，必须立刻停止并退出。

Read and Execute the Workflow, strictly following the instructions in the file: `.agents/workflows/sk-enter-phase/WORKFLOW.md`.
