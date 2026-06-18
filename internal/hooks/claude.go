// Package hooks holds the Claude Code hook adapter for taskbaton: the settings.json
// merge/uninstall logic and the hook stdout envelope. All transformations are pure
// (operate on a decoded settings map) so they can be unit-tested without disk IO.
package hooks

import (
	"encoding/json"
	"strings"
)

// marker identifies a hook command owned by taskbaton. install/uninstall only ever
// touch entries whose command contains it — foreign hooks are left untouched.
const marker = "taskbaton hooks run"

// spec describes one Claude Code hook taskbaton registers.
type spec struct {
	event   string // Claude Code event name, e.g. "SessionStart"
	matcher string // matcher string ("" → omit, fires on every occurrence)
	sub     string // taskbaton subcommand, e.g. "session-start"
	timeout int    // seconds; 0 → omit
}

// specs is the full set of hooks taskbaton installs. Mechanism choices:
//   - SessionStart on startup|resume: inject pickup context once per session.
//   - UserPromptSubmit (no matcher): re-surface constraints on every prompt.
//   - SessionEnd on clear|logout|other: flush the draft once on the way out
//     (SessionEnd fires once at exit, unlike Stop which fires every turn).
var specs = []spec{
	{event: "SessionStart", matcher: "startup|resume", sub: "session-start"},
	{event: "UserPromptSubmit", matcher: "", sub: "user-prompt", timeout: 30},
	{event: "SessionEnd", matcher: "clear|logout|other", sub: "session-end"},
}

// Envelope returns the JSON a context-injecting hook (SessionStart, UserPromptSubmit)
// writes to stdout. When ctx is empty it returns "" so the caller emits nothing —
// an empty stdout is a valid "no context, allow" response on exit 0.
func Envelope(event, ctx string) string {
	if strings.TrimSpace(ctx) == "" {
		return ""
	}
	out := map[string]any{
		"hookSpecificOutput": map[string]any{
			"hookEventName":     event,
			"additionalContext": ctx,
		},
	}
	data, err := json.Marshal(out)
	if err != nil {
		return ""
	}
	return string(data)
}

// EventName maps a taskbaton hook subcommand to its Claude Code event name.
func EventName(sub string) string {
	for _, s := range specs {
		if s.sub == sub {
			return s.event
		}
	}
	return ""
}

// hookEntry builds the {type, command, timeout?} object for a spec.
func (s spec) hookEntry() map[string]any {
	e := map[string]any{
		"type":    "command",
		"command": "taskbaton hooks run " + s.sub,
	}
	if s.timeout > 0 {
		e["timeout"] = s.timeout
	}
	return e
}

// group builds the {matcher?, hooks:[...]} matcher-group for a spec.
func (s spec) group() map[string]any {
	g := map[string]any{
		"hooks": []any{s.hookEntry()},
	}
	if s.matcher != "" {
		g["matcher"] = s.matcher
	}
	return g
}

// Install merges taskbaton's hooks into a decoded settings map, preserving every
// foreign key and hook. It first strips any existing taskbaton entries, so it is
// idempotent — running it twice yields the same result.
func Install(settings map[string]any) map[string]any {
	if settings == nil {
		settings = map[string]any{}
	}
	Uninstall(settings) // clean slate for our own entries

	hooksVal, _ := settings["hooks"].(map[string]any)
	if hooksVal == nil {
		hooksVal = map[string]any{}
	}

	for _, s := range specs {
		arr, _ := hooksVal[s.event].([]any)
		arr = append(arr, s.group())
		hooksVal[s.event] = arr
	}

	settings["hooks"] = hooksVal
	return settings
}

// Uninstall removes only taskbaton-owned hook entries from a decoded settings map.
// Matcher-groups left with no hooks are dropped; events left with no groups are
// dropped; an emptied "hooks" key is removed entirely.
func Uninstall(settings map[string]any) map[string]any {
	if settings == nil {
		return settings
	}
	hooksVal, ok := settings["hooks"].(map[string]any)
	if !ok {
		return settings
	}

	for event, v := range hooksVal {
		groups, ok := v.([]any)
		if !ok {
			continue
		}
		kept := make([]any, 0, len(groups))
		for _, g := range groups {
			gm, ok := g.(map[string]any)
			if !ok {
				kept = append(kept, g)
				continue
			}
			hs, ok := gm["hooks"].([]any)
			if !ok {
				kept = append(kept, g)
				continue
			}
			keptHooks := make([]any, 0, len(hs))
			for _, h := range hs {
				if hm, ok := h.(map[string]any); ok {
					if cmd, _ := hm["command"].(string); strings.Contains(cmd, marker) {
						continue // drop taskbaton's hook
					}
				}
				keptHooks = append(keptHooks, h)
			}
			if len(keptHooks) == 0 {
				continue // group held only taskbaton hooks → drop it
			}
			gm["hooks"] = keptHooks
			kept = append(kept, gm)
		}
		if len(kept) == 0 {
			delete(hooksVal, event)
		} else {
			hooksVal[event] = kept
		}
	}

	if len(hooksVal) == 0 {
		delete(settings, "hooks")
	} else {
		settings["hooks"] = hooksVal
	}
	return settings
}

// Snippet returns the pretty-printed JSON of just taskbaton's hooks block, for
// `hooks install --print` and for pasting into Cursor/Codex config by hand.
func Snippet() string {
	hooksVal := map[string]any{}
	for _, s := range specs {
		arr, _ := hooksVal[s.event].([]any)
		hooksVal[s.event] = append(arr, s.group())
	}
	data, err := json.MarshalIndent(map[string]any{"hooks": hooksVal}, "", "  ")
	if err != nil {
		return ""
	}
	return string(data)
}
