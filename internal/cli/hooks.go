package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/ssthil/taskbaton/internal/hooks"
)

func newHooksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hooks",
		Short: "Wire taskbaton into Claude Code lifecycle hooks",
		Long: `Install lifecycle hooks so agents pick up the baton automatically:

  • session start  — inject constraints + next tasks into the agent
  • every prompt   — re-surface "Do Not Change" constraints
  • session end    — flush the current draft (like 'taskbaton checkpoint')

'hooks install' wires Claude Code by merging into .claude/settings.json.
For Cursor or Codex, run 'hooks install --print' and paste the snippet.`,
	}
	cmd.AddCommand(newHooksInstallCmd())
	cmd.AddCommand(newHooksUninstallCmd())
	cmd.AddCommand(newHooksRunCmd())
	return cmd
}

// claudeSettingsPath returns <cwd>/.claude/settings.json.
func claudeSettingsPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, ".claude", "settings.json"), nil
}

func readSettings(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]any{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return map[string]any{}, nil
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if m == nil {
		m = map[string]any{}
	}
	return m, nil
}

func writeSettings(path string, settings map[string]any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0644)
}

func newHooksInstallCmd() *cobra.Command {
	var printOnly bool

	cmd := &cobra.Command{
		Use:          "install",
		Short:        "Merge taskbaton hooks into .claude/settings.json",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			out := cmd.OutOrStdout()

			if printOnly {
				fmt.Fprintln(out, hooks.Snippet())
				note(out, "paste the above into Cursor/Codex hook config, or run without --print to wire Claude Code")
				return nil
			}

			path, err := claudeSettingsPath()
			if err != nil {
				return err
			}
			settings, err := readSettings(path)
			if err != nil {
				return err
			}
			hooks.Install(settings)
			if err := writeSettings(path, settings); err != nil {
				return err
			}

			success(out, "taskbaton hooks installed → %s", bold(".claude/settings.json"))
			note(out, "session start injects context · prompts re-surface constraints · session end checkpoints")
			return nil
		},
	}
	cmd.Flags().BoolVar(&printOnly, "print", false, "print the hooks JSON snippet instead of writing settings.json")
	return cmd
}

func newHooksUninstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "uninstall",
		Short:        "Remove taskbaton hooks from .claude/settings.json",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			out := cmd.OutOrStdout()

			path, err := claudeSettingsPath()
			if err != nil {
				return err
			}
			settings, err := readSettings(path)
			if err != nil {
				return err
			}
			hooks.Uninstall(settings)
			if err := writeSettings(path, settings); err != nil {
				return err
			}

			success(out, "taskbaton hooks removed from %s", bold(".claude/settings.json"))
			return nil
		},
	}
}

func newHooksRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "run <session-start|user-prompt|session-end>",
		Short:  "Hook handler invoked by the harness (reads payload on stdin)",
		Hidden: true,
		Args:   cobra.ExactArgs(1),
		// A hook must never break the session: always exit 0, emit nothing on error.
		RunE: func(cmd *cobra.Command, args []string) error {
			drainStdin()
			out := cmd.OutOrStdout()

			switch args[0] {
			case "session-start":
				if env := hooks.Envelope("SessionStart", currentContext(false)); env != "" {
					fmt.Fprintln(out, env)
				}
			case "user-prompt":
				if env := hooks.Envelope("UserPromptSubmit", currentContext(true)); env != "" {
					fmt.Fprintln(out, env)
				}
			case "session-end":
				if batonDir, err := batonDirFromCwd(); err == nil {
					_, _, _ = saveCheckpoint(batonDir) // best-effort flush; ignore errors
				}
			}
			return nil
		},
	}
}

// drainStdin consumes the hook JSON payload the harness pipes in so the writer
// never sees a broken pipe. It skips reading when stdin is a terminal to avoid
// blocking if the command is run by hand.
func drainStdin() {
	fi, err := os.Stdin.Stat()
	if err != nil || (fi.Mode()&os.ModeCharDevice) != 0 {
		return
	}
	_, _ = io.Copy(io.Discard, os.Stdin)
}
