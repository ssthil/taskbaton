# Changelog

All notable changes to taskbaton are documented here.

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
