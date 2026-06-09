package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func batonDirFromCwd() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, ".baton"), nil
}

func newReviewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "review",
		Short: "Open .baton/current.md in $EDITOR",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runReview(cmd.OutOrStdout())
		},
	}
}

func runReview(out io.Writer) error {
	batonDir, err := batonDirFromCwd()
	if err != nil {
		return err
	}
	mdPath := filepath.Join(batonDir, "current.md")
	if _, err := os.Stat(mdPath); os.IsNotExist(err) {
		return fmt.Errorf("no baton found — run: taskbaton init")
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = "vi"
	}

	c := exec.Command(editor, mdPath)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}

	info(out, "after editing, run: taskbaton seal --from <tool> --next <tool>")
	return nil
}
