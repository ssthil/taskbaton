package hooks

import (
	"encoding/json"
	"strings"
	"testing"
)

// decode round-trips through JSON so tests operate on the same shape the CLI
// reads from settings.json (all maps are map[string]any, arrays are []any).
func decode(t *testing.T, s string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestInstallPreservesForeignKeys(t *testing.T) {
	settings := decode(t, `{
		"model": "opus",
		"hooks": {
			"PreToolUse": [{"matcher": "Bash", "hooks": [{"type": "command", "command": "echo hi"}]}]
		}
	}`)

	Install(settings)

	if settings["model"] != "opus" {
		t.Errorf("foreign top-level key dropped: %v", settings["model"])
	}
	hooksVal := settings["hooks"].(map[string]any)
	if _, ok := hooksVal["PreToolUse"]; !ok {
		t.Error("foreign PreToolUse hook dropped")
	}
	for _, ev := range []string{"SessionStart", "UserPromptSubmit", "SessionEnd"} {
		if _, ok := hooksVal[ev]; !ok {
			t.Errorf("taskbaton hook %s not installed", ev)
		}
	}
}

func TestInstallIsIdempotent(t *testing.T) {
	settings := map[string]any{}
	Install(settings)
	first, _ := json.Marshal(settings)
	Install(settings)
	second, _ := json.Marshal(settings)
	if string(first) != string(second) {
		t.Errorf("install not idempotent:\nfirst:  %s\nsecond: %s", first, second)
	}

	// Exactly one SessionStart group, owned by taskbaton.
	groups := settings["hooks"].(map[string]any)["SessionStart"].([]any)
	if len(groups) != 1 {
		t.Fatalf("expected 1 SessionStart group after double install, got %d", len(groups))
	}
}

func TestUninstallIsSurgical(t *testing.T) {
	settings := decode(t, `{
		"hooks": {
			"SessionStart": [{"matcher": "startup", "hooks": [{"type": "command", "command": "my-own-thing"}]}]
		}
	}`)
	Install(settings)   // adds taskbaton's SessionStart group alongside the foreign one
	Uninstall(settings) // must remove ONLY taskbaton's

	hooksVal := settings["hooks"].(map[string]any)
	groups, ok := hooksVal["SessionStart"].([]any)
	if !ok {
		t.Fatal("foreign SessionStart group was removed entirely")
	}
	for _, g := range groups {
		hs := g.(map[string]any)["hooks"].([]any)
		for _, h := range hs {
			cmd := h.(map[string]any)["command"].(string)
			if strings.Contains(cmd, marker) {
				t.Errorf("taskbaton hook survived uninstall: %s", cmd)
			}
		}
	}
	// The user's own hook must still be there.
	blob, _ := json.Marshal(settings)
	if !strings.Contains(string(blob), "my-own-thing") {
		t.Errorf("foreign hook lost during uninstall: %s", blob)
	}
}

func TestUninstallDropsEmptyHooksKey(t *testing.T) {
	settings := map[string]any{}
	Install(settings)
	Uninstall(settings)
	if _, ok := settings["hooks"]; ok {
		t.Errorf("empty hooks key should be removed, got: %v", settings["hooks"])
	}
}

func TestEnvelopeEmptyContextIsBlank(t *testing.T) {
	if got := Envelope("SessionStart", ""); got != "" {
		t.Errorf("empty context should yield empty envelope, got %q", got)
	}
	if got := Envelope("SessionStart", "   "); got != "" {
		t.Errorf("whitespace context should yield empty envelope, got %q", got)
	}
}

func TestEnvelopeShape(t *testing.T) {
	out := Envelope("UserPromptSubmit", "do not touch auth.go")
	m := decode(t, out)
	hso, ok := m["hookSpecificOutput"].(map[string]any)
	if !ok {
		t.Fatalf("missing hookSpecificOutput: %s", out)
	}
	if hso["hookEventName"] != "UserPromptSubmit" {
		t.Errorf("wrong event name: %v", hso["hookEventName"])
	}
	if !strings.Contains(hso["additionalContext"].(string), "auth.go") {
		t.Errorf("additionalContext missing payload: %v", hso["additionalContext"])
	}
}

func TestSnippetIsValidJSON(t *testing.T) {
	snip := Snippet()
	m := decode(t, snip)
	if _, ok := m["hooks"]; !ok {
		t.Errorf("snippet missing hooks key: %s", snip)
	}
	if !strings.Contains(snip, "taskbaton hooks run session-start") {
		t.Errorf("snippet missing session-start command: %s", snip)
	}
}
