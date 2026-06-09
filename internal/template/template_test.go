package template

import (
	"strings"
	"testing"
)

func TestRenderContainsStage(t *testing.T) {
	out, err := Render(RenderData{Stage: "my-stage", Status: "active"})
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	if !strings.Contains(out, "my-stage") {
		t.Errorf("output does not contain stage name %q:\n%s", "my-stage", out)
	}
}

func TestRenderContainsAllSections(t *testing.T) {
	out, err := Render(RenderData{Stage: "s", Status: "active"})
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	headers := []string{
		"## Completed",
		"## Decisions",
		"## Next Tasks",
		"## Constraints",
		"## Open Questions",
	}
	for _, h := range headers {
		if !strings.Contains(out, h) {
			t.Errorf("output missing section header %q", h)
		}
	}
}

func TestRenderEmptySlicesFallback(t *testing.T) {
	out, err := Render(RenderData{Stage: "s", Status: "active"})
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	if strings.Contains(out, "<nil>") {
		t.Error("output contains literal \"<nil>\"")
	}
	if strings.Contains(out, "[]") {
		t.Error("output contains literal \"[]\"")
	}
	// Completed, Decisions, NextTasks use "(none yet)"; Constraints and OpenQuestions use "(none)"
	if !strings.Contains(out, "(none yet)") {
		t.Error("output missing fallback \"(none yet)\" for empty list sections")
	}
	if !strings.Contains(out, "(none)") {
		t.Error("output missing fallback \"(none)\" for empty list sections")
	}
}

func TestRenderPopulatedSlices(t *testing.T) {
	out, err := Render(RenderData{
		Stage:     "s",
		Status:    "active",
		Completed: []string{"did A"},
	})
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	if !strings.Contains(out, "- did A") {
		t.Errorf("output does not contain list item %q:\n%s", "- did A", out)
	}
}

func TestRenderSealedFields(t *testing.T) {
	d := RenderData{
		Stage:    "s",
		Status:   "sealed",
		From:     "claude-code",
		SealedAt: "2026-06-09T10:00:00+08:00",
		Next:     "cursor",
	}
	out, err := Render(d)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	for _, want := range []string{"claude-code", "2026-06-09T10:00:00+08:00", "cursor"} {
		if !strings.Contains(out, want) {
			t.Errorf("output does not contain %q:\n%s", want, out)
		}
	}
}

func TestRenderOptionalFieldsOmitted(t *testing.T) {
	out, err := Render(RenderData{Stage: "s", Status: "active", From: "", SealedAt: "", Next: ""})
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	for _, label := range []string{"**From**", "**Sealed**", "**Next**"} {
		if strings.Contains(out, label) {
			t.Errorf("output should not contain %q when field is empty:\n%s", label, out)
		}
	}
}
