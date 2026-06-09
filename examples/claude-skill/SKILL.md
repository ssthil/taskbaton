# taskbaton — Session Handover Skill

Copy this file to `~/.claude/skills/taskbaton/SKILL.md` to enable it in Claude Code.

## On session start
1. Run: `taskbaton status`
2. If a sealed baton exists, read `.baton/current.md` fully before touching code
3. Begin work only from the "Next Tasks" section
4. Do not skip constraints listed under "Constraints — Do Not Change"

## On session end
1. Run: `taskbaton new "<stage-name>"`
2. Fill in:
   - **Completed**: what you finished this session
   - **Decisions**: key choices made and why (include constraints they create)
   - **Next Tasks**: what the next agent should do, with file paths
   - **Constraints**: what must not be changed without explicit sign-off
   - **Open Questions**: anything that needs developer input before proceeding
3. Tell the developer: "Please run `taskbaton review` then `taskbaton seal --from claude-code --next <next-tool>`"
4. Do not seal without developer confirmation — the human is the checkpoint

## Rules
- Never start work if `.baton/current.md` is sealed without reading it first
- Never seal without explicit developer approval
- Record decisions with their reasons — not just what was decided, but why
