# Changelog

All notable changes to taskbaton are documented here.

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
