# 🪄 taskbaton

<p align="center">
  <strong>Pass the baton. Every agent picks up exactly where the last one left off.</strong><br/>
  The CLI tool that carries work state between agentic tool sessions — so nothing gets lost.
</p>

<p align="center">
  <a href="https://github.com/ssthil/taskbaton">GitHub</a> ·
  <a href="https://github.com/ssthil/taskbaton/issues">Issues</a> ·
  <a href="./CHANGELOG.md">Changelog</a>
</p>

<p align="center">
  <a href="https://github.com/ssthil/taskbaton/actions/workflows/ci.yml"><img src="https://img.shields.io/github/actions/workflow/status/ssthil/taskbaton/ci.yml?branch=main&style=flat-square&label=CI&labelColor=1a1a1a" alt="CI"/></a>
  <a href="https://github.com/ssthil/taskbaton/releases"><img src="https://img.shields.io/badge/release-v0.1.0-blue?style=flat-square&labelColor=1a1a1a" alt="release"/></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-3b82f6?style=flat-square&labelColor=1a1a1a" alt="License: MIT"/></a>
</p>

---

## Every Agent Starts Blind. Until Now.

> When you switch tools mid-project — Gemini plans, Claude Code builds, Cursor reviews, Codex refactors — every new session starts from scratch. Decisions are lost. Constraints are forgotten. Work gets redone.
>
> **taskbaton** manages a `.baton/` directory in your project root. Any agent reads `current.md` at the start of a session and immediately knows the stage, what was decided, what must not change, and what to do next. The human reviews and seals every handover — you stay the checkpoint between every stage.

- **`taskbaton init`**: Scaffold `.baton/` once. Every stage that follows lives there.
- **`taskbaton new` + `review`**: Create a stage, fill in decisions and next tasks — you or your agent drafts it, you confirm.
- **`taskbaton seal`**: Lock it. Stamp the timestamp. Archive it. The baton passes.

## How the flow works

```
taskbaton init              ← scaffold .baton/ once in the project root
         │
         ▼
taskbaton new "planning"    ← Gemini CLI: start the first stage
taskbaton review            ← fill decisions, next tasks, constraints in $EDITOR
taskbaton seal \
  --from gemini-cli \       ← human reviews and locks the handover
  --next claude-code
         │
         ▼
  Claude Code reads .baton/current.md    ← knows exactly where to start
  ... builds the feature ...
taskbaton new "feature/auth"
taskbaton review            ← agent proposes, developer confirms
taskbaton seal \
  --from claude-code \
  --next cursor
         │
         ▼
  Cursor reads .baton/current.md         ← picks up without a brief
  taskbaton seal --from cursor --next codex
         │
         ▼
taskbaton log               ← full stage history, every decision preserved
```

Human is the checkpoint at every `seal`. Agent drafts, human confirms, baton passes.

## Installation

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

## Quick start

```bash
cd your-project
taskbaton init                                        # scaffold .baton/ once
taskbaton new "planning"                              # create the first stage
taskbaton review                                      # fill in .baton/current.md in $EDITOR
taskbaton seal --from gemini-cli --next claude-code   # lock and archive
taskbaton status                                      # show stage + seal state
taskbaton next                                        # print next tasks for incoming agent
```

## How it works

1. **`taskbaton init`** scaffolds `.baton/` in your project root with `current.md`, `current.json`, `config.yaml`, and a `history/` archive directory. Run once per project.

2. **`taskbaton new <stage>`** creates a blank baton for the named stage with all section headers pre-filled. Pair it with `taskbaton review` to edit — or let your agent draft it and you confirm.

3. **`taskbaton review`** opens `.baton/current.md` in `$EDITOR`. Fill in completed work, decisions (with reasons), next tasks, constraints that must not change, and open questions.

4. **`taskbaton seal --from <tool> --next <tool>`** locks the baton. Sets status to `sealed`, stamps the timestamp, writes both `.md` and `.json`, and archives the stage to `.baton/history/`. Only a human should run this — it is the handover gate.

5. **`taskbaton next`** prints the `## Next Tasks` section to stdout — the exact lines the incoming agent starts from.

6. **`taskbaton export`** outputs `current.json` to stdout — pipe it into scripts, CI checks, or agent context loaders.

## The baton format

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
- /login and /refresh endpoints, unit tests passing (84% coverage)

## Decisions
- RS256 not HS256 — asymmetric keys required by security team
- Token expiry: 15min access / 7d refresh — do not change without PM sign-off
- Refresh tokens in Redis, not DB — performance decision

## Next Tasks
- Implement RBAC on top of existing auth
- Start at src/middleware/rbac.go
- Reference: docs/roles.md for role hierarchy

## Constraints — Do Not Change
- src/middleware/auth.go — sealed by security review
- Token storage: keep in Redis
- Auth errors: always 401, never 403 (API contract)

## Open Questions
- Rate limiting strategy — ask developer before implementing
```

A machine-readable `.baton/current.json` is always generated alongside it.

## Commands

| Command | Description |
|---|---|
| `taskbaton init` | Scaffold `.baton/` in the current project |
| `taskbaton new <stage>` | Create a new open baton stage |
| `taskbaton review` | Open `.baton/current.md` in `$EDITOR` |
| `taskbaton seal --from <tool> --next <tool>` | Seal, stamp, and archive the current baton |
| `taskbaton next` | Print next tasks for the incoming agent |
| `taskbaton status` | Show current stage and seal state |
| `taskbaton log` | List all archived baton stages |
| `taskbaton export` | Output current baton as JSON |

## AGENTS.md convention

Paste this into your project's `AGENTS.md` so every agent picks it up automatically:

```markdown
## Baton Protocol
- Session start: run `taskbaton status` and read `.baton/current.md` fully
- Session end: run `taskbaton new <stage>`, fill template, get human review,
  then `taskbaton seal --from <your-tool> --next <next-tool>`
- Never start work if a sealed baton exists without reading it first
- Never seal without developer confirmation
```

A ready-made Claude Code skill lives in [`examples/claude-skill/SKILL.md`](examples/claude-skill/SKILL.md).

## Philosophy

- **Human seals, agent drafts.** The agent fills in what it did. The developer reads, edits if needed, and runs `seal`. No baton passes without a human in the loop.
- **Decisions with reasons, not just outcomes.** A constraint without a reason gets overridden by the next agent. Record the *why* so it survives the handover.
- **Plain files, no server.** `.baton/` is just a directory of text files. Check it into git, diff it, grep it, edit it by hand. No daemon, no cloud sync, no API keys.

## Part of a suite

| Tool | Layer | Problem solved |
|---|---|---|
| [human-steerkit](https://github.com/ssthil/human-steerkit) | Within a session | Human stays in control of agent decisions |
| [llmroute](https://github.com/ssthil/llmroute) | Across providers | Routes prompts to the right LLM dynamically |
| **taskbaton** | Across sessions | Preserves work state between tool switches |

Three tools. One philosophy: **humans steer, agents execute, nothing gets lost.**

## Development

```bash
make build    # build to ./bin/taskbaton
make test     # run tests with race detector
make lint     # go vet + gofmt
make snapshot # local cross-platform build via goreleaser
```

## License

MIT — see [LICENSE](LICENSE).
