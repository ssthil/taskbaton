# taskbaton

[![CI](https://github.com/ssthil/taskbaton/actions/workflows/ci.yml/badge.svg)](https://github.com/ssthil/taskbaton/actions/workflows/ci.yml)
[![release](https://img.shields.io/badge/release-v0.1.0-blue)](https://github.com/ssthil/taskbaton/releases)
[![MIT License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

> Pass the baton. Every agentic tool picks up exactly where the last one left off.
> **The baton your AI tools actually pass.**

When you switch agentic tools mid-project — Gemini plans → Claude Code builds → Cursor reviews → Codex refactors — every tool starts blind. Decisions are lost, constraints are forgotten, work gets undone.

`taskbaton` manages a `.baton/` directory in your project root. Any agent reads `current.md` at the start of a session and immediately knows where it stands. Human reviews and seals every handover before it passes — the human stays the checkpoint between every stage.

---

## Install

**Homebrew (macOS / Linux)**
```bash
brew tap ssthil/taskbaton
brew install taskbaton
```

**Go**
```bash
go install github.com/ssthil/taskbaton@latest
```

**Binary releases** — download from [Releases](https://github.com/ssthil/taskbaton/releases), verify with `checksums.txt`.

**Upgrade**
```bash
brew upgrade taskbaton
# or
go install github.com/ssthil/taskbaton@latest
```

---

## Quick start

```bash
cd your-project
taskbaton init              # scaffold .baton/
taskbaton new "planning"    # create the first stage
taskbaton review            # fill in decisions and next tasks in $EDITOR
taskbaton seal --from gemini-cli --next claude-code
```

The next agent reads `.baton/current.md` and runs `taskbaton next`.

---

## The workflow

```
── Gemini plans ──────────────────────────────────────────────────────────
  taskbaton init
  taskbaton new "planning"
  taskbaton review                  ← fill in decisions + next tasks
  taskbaton seal --from gemini-cli --next claude-code

── Claude Code builds ────────────────────────────────────────────────────
  (reads .baton/current.md at session start)
  ... builds the feature ...
  taskbaton new "feature/auth"
  taskbaton review                  ← agent proposes, dev reviews
  taskbaton seal --from claude-code --next cursor

── Cursor reviews ────────────────────────────────────────────────────────
  (reads .baton/current.md — knows exactly where to start)
  taskbaton seal --from cursor --next codex
```

Human is the checkpoint at every `seal`. Agent drafts, human confirms, baton passes.

---

## Commands

| Command | Description |
|---|---|
| `taskbaton init` | Scaffold `.baton/` in the current project |
| `taskbaton new <stage>` | Create a new open baton stage |
| `taskbaton review` | Open `.baton/current.md` in `$EDITOR` |
| `taskbaton seal --from <tool> --next <tool>` | Seal the current baton and archive it |
| `taskbaton next` | Print next tasks for the incoming agent |
| `taskbaton status` | Show current stage and seal state |
| `taskbaton log` | List all archived baton stages |
| `taskbaton export` | Export current baton as JSON |

---

## Baton format

Every `.baton/current.md` follows this structure:

```markdown
# Baton — feature/auth

**Stage**: feature/auth
**Status**: sealed
**From**: claude-code
**Sealed**: 2026-06-08T10:42:00+08:00
**Next**: cursor

## Completed
- JWT middleware in src/middleware/auth.go

## Decisions
- RS256 not HS256 — asymmetric keys required by security team

## Next Tasks
- Implement RBAC on top of existing auth

## Constraints — Do Not Change
- src/middleware/auth.go — sealed by security review

## Open Questions
- Rate limiting strategy — ask developer before implementing
```

A machine-readable `.baton/current.json` is always generated alongside it.

---

## AGENTS.md convention

Paste this into your project's `AGENTS.md`:

```markdown
## Baton Protocol
- Session start: run `taskbaton status` and read `.baton/current.md`
- Session end: run `taskbaton new <stage>`, fill template, get human review,
  then `taskbaton seal --from <your-tool> --next <next-tool>`
- Never start work if a sealed baton exists without reading it first
```

---

## Project layout

```
.baton/
  current.md          ← active baton (agent reads this on start)
  current.json        ← machine-readable twin
  config.yaml         ← project config (author, etc.)
  history/
    2026-06-08_feature-auth.md
    2026-06-07_planning.md
```

---

## Development

```bash
make build    # build to ./bin/taskbaton
make test     # run tests with race detector
make lint     # go vet + gofmt
make snapshot # cross-platform snapshot via goreleaser
```

---

## Part of a suite

| Tool | Layer | Problem solved |
|---|---|---|
| human-steerkit | Within a session | Human stays in control of agent decisions |
| llmroute | Across providers | Routes prompts to the right LLM dynamically |
| taskbaton | Across sessions | Preserves work state between tool switches |

Three tools. One philosophy: **humans steer, agents execute, nothing gets lost.**

---

## License

MIT — see [LICENSE](LICENSE).
