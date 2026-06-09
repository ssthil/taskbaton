package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export current baton as JSON",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runExport(cmd)
		},
	}
	cmd.Flags().Bool("json", true, "output as JSON (reserved for future formats)")
	return cmd
}

func runExport(cmd *cobra.Command) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	jsonPath := filepath.Join(cwd, ".baton", "current.json")

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no baton found — run: taskbaton init")
		}
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return nil
}
