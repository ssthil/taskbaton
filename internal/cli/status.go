package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/ssthil/taskbaton/internal/baton"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current stage and seal state",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runStatusCmd(cmd.OutOrStdout())
		},
	}
}

func runStatusCmd(out io.Writer) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	batonDir := filepath.Join(cwd, ".baton")

	b, err := baton.Read(batonDir)
	if err != nil {
		fmt.Fprintln(out, "no baton found — run: taskbaton init")
		return nil
	}

	fmt.Fprintf(out, "  %-8s %s\n", "Stage:", b.Stage)
	fmt.Fprintf(out, "  %-8s %s\n", "Status:", b.Status)
	if b.From != "" {
		fmt.Fprintf(out, "  %-8s %s\n", "From:", b.From)
	}
	if b.Next != "" {
		fmt.Fprintf(out, "  %-8s %s\n", "Next:", b.Next)
	}
	if b.SealedAt != "" {
		fmt.Fprintf(out, "  %-8s %s\n", "Sealed:", b.SealedAt)
	}
	return nil
}
