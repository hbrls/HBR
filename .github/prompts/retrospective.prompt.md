---
description: Execute Command: .agents/commands/retrospective/COMMAND.md
agent: agent
metadata:
   version: 0.0.1
---

# Retrospective 复盘

1. Read the command definition:
   `view_file .agents/commands/retrospective/COMMAND.md`
2. Execute the command strictly following the instructions in the file.
3. **CRITICAL**: When the command requires user interaction, elicitation, or waiting for feedback (e.g. "WAIT FOR USER RESPONSE"), you MUST use the `notify_user` tool. 
   - Present the options (e.g. 1-9) in the `Message` argument.
   - Set `BlockedOnUser` to `true`.
   - Do NOT proceed until you receive the user's response via the tool output in the next turn.
