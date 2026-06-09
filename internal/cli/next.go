package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/ssthil/taskbaton/internal/baton"
)

func newNextCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "next",
		Short: "Print next tasks for the incoming agent",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runNext(cmd.OutOrStdout())
		},
	}
}

func runNext(out io.Writer) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	batonDir := filepath.Join(cwd, ".baton")

	b, err := baton.Read(batonDir)
	if err != nil {
		return err
	}

	if len(b.NextTasks) == 0 {
		fmt.Fprintln(out, "(no next tasks recorded)")
		return nil
	}

	fmt.Fprintln(out, bold("Next Tasks"))
	for _, t := range b.NextTasks {
		fmt.Fprintf(out, "  • %s\n", t)
	}
	return nil
}
