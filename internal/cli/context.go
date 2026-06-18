package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ssthil/taskbaton/internal/baton"
	"github.com/ssthil/taskbaton/internal/template"
)

func newContextCmd() *cobra.Command {
	var constraintsOnly bool

	cmd := &cobra.Command{
		Use:   "context",
		Short: "Print baton context for an agent to pick up",
		Long: `Emit a compact, plain-text summary of the current baton — stage,
constraints, and next tasks — ready to inject into an agent's context.

This is the portable core used by 'taskbaton hooks'. Pipe it into any tool
that accepts extra context. When no baton exists it prints nothing and exits 0,
so it is safe to call unconditionally at session start.

  taskbaton context                     full pickup context
  taskbaton context --constraints-only  just the "Do Not Change" block`,
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			text := currentContext(constraintsOnly)
			if text != "" {
				fmt.Fprintln(cmd.OutOrStdout(), text)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&constraintsOnly, "constraints-only", false,
		"emit only the Constraints — Do Not Change block")
	return cmd
}

// currentContext loads the current baton and renders it for injection. It never
// errors: a missing .baton/, missing baton, or empty baton all yield "" so the
// caller (a CLI command or a hook) can stay silent and exit 0.
func currentContext(constraintsOnly bool) string {
	batonDir, err := batonDirFromCwd()
	if err != nil {
		return ""
	}
	if _, err := os.Stat(batonDir); err != nil {
		return ""
	}
	b, err := baton.Read(batonDir)
	if err != nil {
		return ""
	}
	return template.RenderContext(template.RenderData{
		Stage:         b.Stage,
		Status:        b.Status,
		From:          b.From,
		SealedAt:      b.SealedAt,
		Next:          b.Next,
		Completed:     b.Completed,
		Decisions:     b.Decisions,
		NextTasks:     b.NextTasks,
		Constraints:   b.Constraints,
		OpenQuestions: b.OpenQuestions,
	}, constraintsOnly)
}
