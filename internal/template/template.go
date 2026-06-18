package template

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type RenderData struct {
	Stage         string
	Status        string
	From          string
	CreatedAt     string
	SealedAt      string
	Next          string
	Completed     []string
	Decisions     []string
	NextTasks     []string
	Constraints   []string
	OpenQuestions []string
}

const batonTemplate = `# Baton — {{.Stage}}

**Stage**: {{.Stage}}
**Status**: {{.Status}}
{{- if .CreatedAt}}
**Created**: {{.CreatedAt}}
{{- end}}
{{- if .From}}
**From**: {{.From}}
{{- end}}
{{- if .SealedAt}}
**Sealed**: {{.SealedAt}}
{{- end}}
{{- if .Next}}
**Next**: {{.Next}}
{{- end}}

## Completed
{{range .Completed}}- {{.}}
{{else}}(none yet)
{{end}}
## Decisions
{{range .Decisions}}- {{.}}
{{else}}(none yet)
{{end}}
## Next Tasks
{{range .NextTasks}}- {{.}}
{{else}}(none yet)
{{end}}
## Constraints — Do Not Change
{{range .Constraints}}- {{.}}
{{else}}(none)
{{end}}
## Open Questions
{{range .OpenQuestions}}- {{.}}
{{else}}(none)
{{end}}`

// RenderContext produces a compact, plain-text (no ANSI) summary of the baton
// for injection into an agent's context — via `taskbaton context` or a hook.
//
// When constraintsOnly is true only the "Constraints — Do Not Change" block is
// emitted (used to re-surface sealed constraints on every prompt). When there is
// nothing worth injecting, it returns "" — callers should emit nothing and exit 0.
func RenderContext(d RenderData, constraintsOnly bool) string {
	var b strings.Builder

	writeList := func(heading string, items []string) {
		if len(items) == 0 {
			return
		}
		fmt.Fprintf(&b, "%s\n", heading)
		for _, it := range items {
			fmt.Fprintf(&b, "- %s\n", it)
		}
		b.WriteString("\n")
	}

	if constraintsOnly {
		writeList("taskbaton — constraints that must not change:", d.Constraints)
		return strings.TrimSpace(b.String())
	}

	// Nothing useful to hand off — let the caller stay silent.
	if len(d.Constraints) == 0 && len(d.NextTasks) == 0 {
		return ""
	}

	header := fmt.Sprintf("taskbaton — current stage: %s (%s", d.Stage, d.Status)
	if d.From != "" {
		header += ", from " + d.From
	}
	if d.Next != "" {
		header += " → " + d.Next
	}
	header += ")"
	fmt.Fprintf(&b, "%s\n\n", header)

	writeList("Constraints — Do Not Change:", d.Constraints)
	writeList("Next Tasks:", d.NextTasks)

	return strings.TrimSpace(b.String())
}

func Render(d RenderData) (string, error) {
	t, err := template.New("baton").Parse(batonTemplate)
	if err != nil {
		return "", fmt.Errorf("template render: %w", err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, d); err != nil {
		return "", fmt.Errorf("template render: %w", err)
	}
	return buf.String(), nil
}
