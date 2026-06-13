package template

import (
	"bytes"
	"fmt"
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
