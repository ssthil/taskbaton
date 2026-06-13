// Package cli wires the Cobra command tree for taskbaton.
package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

// SetVersion lets the build inject the release version.
func SetVersion(v string) {
	if v != "" {
		version = v
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "taskbaton",
		Short: "Pass work state between agentic tool sessions",
		Long: `taskbaton manages a .baton/ directory in your project root.
Any agent reads current.md at the start of a session and immediately
knows where it stands. Human reviews and seals every handover before
it passes — the human stays the checkpoint between every stage.`,
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runStatus(cmd.OutOrStdout())
		},
	}

	root.AddCommand(newInitCmd())
	root.AddCommand(newNewCmd())
	root.AddCommand(newReviewCmd())
	root.AddCommand(newSealCmd())
	root.AddCommand(newNextCmd())
	root.AddCommand(newStatusCmd())
	root.AddCommand(newLogCmd())
	root.AddCommand(newExportCmd())
	root.AddCommand(newMCPCmd())
	return root
}

func printBanner(out io.Writer) {
	const w = 52
	fmt.Fprintln(out)
	boxTop(out, w)
	boxRow(out, w, "  taskbaton "+version, func(s string) string { return bold(cyan(s)) })
	boxRow(out, w, "  the baton your AI tools actually pass", dim)
	boxBottom(out, w)
}

func runStatus(out io.Writer) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	batonDir := cwd + "/.baton"
	if _, err := os.Stat(batonDir); os.IsNotExist(err) {
		printBanner(out)
		fmt.Fprintln(out)
		fmt.Fprintln(out, bold("Get started"))
		step(out, 1, "taskbaton init", "scaffold .baton/ in your project")
		step(out, 2, "taskbaton new <stage>", "create the first baton stage")
		step(out, 3, "taskbaton review", "fill in decisions and next tasks")
		step(out, 4, "taskbaton seal --from <tool> --next <tool>", "lock and archive")
		fmt.Fprintln(out)
		fmt.Fprintln(out, bold("Tips"))
		tip(out, "taskbaton status", "show current stage and seal state")
		tip(out, "taskbaton log", "show full stage history")
		tip(out, "taskbaton next", "print next tasks for incoming agent")
		tip(out, "taskbaton export --json", "pipe-friendly JSON output")
		fmt.Fprintln(out)
		return nil
	}
	return runStatus2(out, batonDir)
}

func runStatus2(out io.Writer, batonDir string) error {
	printBanner(out)
	fmt.Fprintln(out)
	note(out, "run %s for stage details · %s to seal",
		bold("taskbaton status"), bold("taskbaton seal --from <tool> --next <tool>"))
	return nil
}

// Execute builds the command tree and runs it.
func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
