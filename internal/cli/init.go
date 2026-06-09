package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ssthil/taskbaton/internal/config"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Scaffold .baton/ in the current project",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()

			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("cannot determine working directory: %w", err)
			}

			batonDir := filepath.Join(cwd, ".baton")

			if _, err := os.Stat(batonDir); err == nil {
				warn(out, "already initialized — %s exists", batonDir)
				return nil
			}

			if err := os.Mkdir(batonDir, 0700); err != nil {
				return fmt.Errorf("create .baton/: %w", err)
			}

			historyDir := filepath.Join(batonDir, "history")
			if err := os.Mkdir(historyDir, 0700); err != nil {
				return fmt.Errorf("create .baton/history/: %w", err)
			}

			reader := bufio.NewReader(os.Stdin)

			defaultName := filepath.Base(cwd)
			projectName := prompt(reader, "Project name", defaultName)
			author := prompt(reader, "Author", "")

			cfg := config.Config{
				ProjectName: projectName,
				Author:      author,
			}

			if err := config.Save(batonDir, cfg); err != nil {
				return err
			}

			success(out, "Initialized %s", bold(".baton/"))
			info(out, "Project: %s", bold(projectName))
			if author != "" {
				info(out, "Author:  %s", author)
			}
			note(out, "Run %s to record your first task handoff.", bold("taskbaton push"))

			return nil
		},
	}
}

func prompt(reader *bufio.Reader, label, defaultVal string) string {
	if defaultVal != "" {
		fmt.Fprintf(os.Stderr, "%s [%s]: ", label, defaultVal)
	} else {
		fmt.Fprintf(os.Stderr, "%s: ", label)
	}
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return defaultVal
	}
	return line
}
