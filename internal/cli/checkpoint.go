package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/ssthil/taskbaton/internal/baton"
)

func newCheckpointCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "checkpoint",
		Short: "Save current baton state without sealing",
		Long: `Persist the current open baton to disk right now.

Use this mid-session before a tool's usage limit is reached — gives you
a saved draft to review rather than a blank template if the agent cuts out.

The baton stays open. Edit it with 'taskbaton review', seal it with
'taskbaton seal --from <tool> --next <tool>' when ready to hand off.`,
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			out := cmd.OutOrStdout()

			batonDir, err := batonDirFromCwd()
			if err != nil {
				return err
			}
			if _, err := os.Stat(batonDir); os.IsNotExist(err) {
				return fmt.Errorf(".baton/ not found — run: taskbaton init")
			}

			b, err := baton.Read(batonDir)
			if err != nil {
				return fmt.Errorf("no current baton — run: taskbaton new <stage> first")
			}
			if b.Status == "sealed" {
				warn(out, "baton is already sealed — run: taskbaton new <stage> to start the next stage")
				return nil
			}

			if err := baton.Write(batonDir, b); err != nil {
				return err
			}

			elapsed := sessionAge(batonDir, b)
			success(out, "checkpoint saved at %s  %s  stage: %s  %s  open %s",
				time.Now().Format("15:04:05"),
				gray(glyphDot), bold(b.Stage),
				gray(glyphDot), gray(elapsed))
			note(out, "run %s to add notes · %s to finalize",
				bold("taskbaton review"),
				bold("taskbaton seal --from <tool> --next <tool>"))
			return nil
		},
	}
}

// sessionAge returns a human-readable elapsed time since the baton was created.
// Falls back to file mtime if CreatedAt is not set (older batons).
func sessionAge(batonDir string, b baton.Baton) string {
	var since time.Duration
	if b.CreatedAt != "" {
		if t, err := time.Parse(time.RFC3339, b.CreatedAt); err == nil {
			since = time.Since(t)
		}
	}
	if since == 0 {
		if fi, err := os.Stat(batonDir + "/current.json"); err == nil {
			since = time.Since(fi.ModTime())
		}
	}
	if since == 0 {
		return "unknown duration"
	}
	mins := int(since.Minutes())
	if mins < 1 {
		return "< 1 min"
	}
	if mins < 60 {
		return fmt.Sprintf("%d min", mins)
	}
	return fmt.Sprintf("%dh %dm", mins/60, mins%60)
}
