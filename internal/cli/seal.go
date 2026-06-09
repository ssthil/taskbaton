package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/ssthil/taskbaton/internal/baton"
	"github.com/ssthil/taskbaton/internal/history"
)

func newSealCmd() *cobra.Command {
	var from, next string

	cmd := &cobra.Command{
		Use:   "seal",
		Short: "Seal the current baton and archive it",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			out := cmd.OutOrStdout()

			if from == "" {
				return fmt.Errorf("--from is required")
			}
			if next == "" {
				return fmt.Errorf("--next is required")
			}

			batonDir, err := batonDirFromCwd()
			if err != nil {
				return err
			}
			if _, err := os.Stat(batonDir); os.IsNotExist(err) {
				return fmt.Errorf(".baton/ not found — run: taskbaton init")
			}

			b, err := baton.Read(batonDir)
			if err != nil {
				return fmt.Errorf("no current baton found: %w", err)
			}

			if b.Status == "sealed" {
				warn(out, "baton is already sealed")
				return nil
			}

			b.Status = "sealed"
			b.From = from
			b.Next = next
			b.SealedAt = time.Now().Format(time.RFC3339)

			if err := history.Archive(batonDir, b.Stage, time.Now()); err != nil {
				return err
			}

			if err := baton.Write(batonDir, b); err != nil {
				return err
			}

			success(out, "baton sealed: %s → %s", from, next)
			note(out, "archived to .baton/history/")
			info(out, "next agent should read .baton/current.md and run: taskbaton next")
			return nil
		},
	}

	cmd.Flags().StringVar(&from, "from", "", "tool/agent handing off")
	cmd.Flags().StringVar(&next, "next", "", "tool/agent receiving the baton")
	return cmd
}
