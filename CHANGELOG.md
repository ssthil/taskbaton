# Changelog

All notable changes to taskbaton are documented here.

## [Unreleased]

### Added
- `taskbaton hooks install` — wire taskbaton into Claude Code lifecycle hooks (merges into `.claude/settings.json`, idempotent, preserves existing settings); `--print` emits the snippet for Cursor/Codex
- `taskbaton hooks uninstall` — surgically remove only taskbaton's hook entries
- Hook behaviours: inject constraints + next tasks at session start, re-surface *Constraints — Do Not Change* on every prompt, flush the draft on session end. Sealing stays human-gated
- `taskbaton context` — plain-text baton summary for agent pickup; prints nothing and exits 0 when no baton exists, so it is safe to call unconditionally
- `examples/hooks/` — Cursor/Codex/manual snippets

### Fixed
- `taskbaton init` post-init hint pointed at a nonexistent `taskbaton push` command; now points at `taskbaton new <stage>`

## [0.3.0] - 2026-06-13

### Added
- `taskbaton checkpoint` — mid-session draft save without sealing or archiving; persists current state so a usage-wall doesn't leave the human with a blank template
- Session duration nudge in `taskbaton status` — warns after 45 min if a baton is still open, with a reminder to checkpoint
- `created_at` field on every new baton (RFC3339); used for accurate session duration tracking with fallback to file mtime for older batons

## [0.2.0] - 2026-06-13

### Added
- `taskbaton mcp` — MCP server over stdio; exposes baton state as native context to Claude Code and other MCP hosts
- Resources: `baton://current` (Markdown), `baton://current.json`, `baton://history`
- Tools: `get_status`, `get_next_tasks`, `get_constraints`
- Homebrew tap via `brew tap ssthil/senthil-tools`

## [0.1.0] - 2026-06-09

### Added
- `taskbaton init` — scaffold `.baton/` in project root
- `taskbaton new <stage>` — create a new open baton stage
- `taskbaton review` — open `.baton/current.md` in `$EDITOR`
- `taskbaton seal --from <tool> --next <tool>` — lock, stamp, and archive
- `taskbaton next` — print next tasks for incoming agent
- `taskbaton status` — show current stage and seal state
- `taskbaton log` — list full stage history
- `taskbaton export` — pipe-friendly JSON output
- Dual-format output: `.md` + `.json` on every write
- History archive in `.baton/history/` with date-prefixed filenames
- GoReleaser cross-platform builds (linux/darwin/windows × amd64/arm64)
- GitHub Actions CI (multi-OS) and Release workflows
