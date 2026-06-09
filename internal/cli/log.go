package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ssthil/taskbaton/internal/history"
)

func newLogCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "log",
		Short: "List all archived baton stages",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runLog(cmd.OutOrStdout())
		},
	}
}

func runLog(out io.Writer) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	batonDir := filepath.Join(cwd, ".baton")

	files, err := history.List(batonDir)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		fmt.Fprintln(out, "no history yet")
		return nil
	}

	for _, f := range files {
		fmt.Fprintf(out, "  %s\n", strings.TrimSuffix(f, ".md"))
	}
	return nil
}
