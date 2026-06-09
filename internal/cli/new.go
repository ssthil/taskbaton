package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ssthil/taskbaton/internal/baton"
)

func newNewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "new <stage>",
		Short: "Create a new baton stage",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stage := args[0]
			out := cmd.OutOrStdout()

			batonDir, err := batonDirFromCwd()
			if err != nil {
				return err
			}
			if _, err := os.Stat(batonDir); os.IsNotExist(err) {
				return fmt.Errorf(".baton/ not found — run: taskbaton init")
			}

			existing, readErr := baton.Read(batonDir)
			if readErr == nil && existing.Status == "open" {
				warn(out, "an open baton already exists for stage %q — seal it first with: taskbaton seal --from <tool> --next <tool>", existing.Stage)
				return nil
			}

			b := baton.New(stage)
			if err := baton.Write(batonDir, b); err != nil {
				return err
			}

			success(out, "new baton stage %q created", stage)
			info(out, "fill in .baton/current.md then run: taskbaton review")
			return nil
		},
	}
}
