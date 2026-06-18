# taskbaton hooks — harness snippets

`taskbaton hooks install` wires Claude Code automatically. For other harnesses,
paste the relevant snippet below. Every harness ends up calling the same three
`taskbaton hooks run …` commands, so behaviour is identical:

- **session start** → injects constraints + next tasks
- **every prompt** → re-surfaces *Constraints — Do Not Change*
- **session end** → flushes the current draft (like `taskbaton checkpoint`)

`taskbaton hooks run` reads the harness payload on stdin, emits the right output
format, and always exits 0 — a hook never breaks your session. `seal` stays
human-only; hooks never seal.

Generate the canonical Claude Code JSON any time with:

```bash
taskbaton hooks install --print
```

## Claude Code

Run `taskbaton hooks install` (writes `.claude/settings.json`), or paste the
`--print` output into `.claude/settings.json` by hand.

## Cursor

Cursor follows the same hook model. Add to your Cursor hooks config:

```json
{
  "hooks": {
    "beforeSubmitPrompt": [
      { "command": "taskbaton hooks run user-prompt" }
    ],
    "sessionStart": [
      { "command": "taskbaton hooks run session-start" }
    ],
    "sessionEnd": [
      { "command": "taskbaton hooks run session-end" }
    ]
  }
}
```

> Event names vary by Cursor version — map them to your version's lifecycle
> events. The `taskbaton hooks run …` commands stay the same.

## Codex

Map Codex's lifecycle hooks to the same three commands:

| Codex lifecycle | Command |
|---|---|
| session/turn start | `taskbaton hooks run session-start` |
| before user prompt | `taskbaton hooks run user-prompt` |
| session end        | `taskbaton hooks run session-end` |

## Manual / any tool

No hook system? Add to your `AGENTS.md` instead:

```markdown
- Session start: run `taskbaton context`
- Session end:   run `taskbaton checkpoint`, then ask the human to `seal`
```
