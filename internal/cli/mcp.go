package cli

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ssthil/taskbaton/internal/mcp"
)

func newMCPCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mcp",
		Short: "Start the MCP server (stdio transport)",
		Long: `Start a Model Context Protocol server over stdio.

Register this command in your MCP host config so agents can read baton
state as native context — no copy-paste required.

Claude Code (~/.claude/claude_desktop_config.json):
  {
    "mcpServers": {
      "taskbaton": {
        "command": "taskbaton",
        "args": ["mcp"]
      }
    }
  }`,
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			batonDir := filepath.Join(cwd, ".baton")
			return mcp.New(batonDir).Serve(os.Stdin, os.Stdout)
		},
	}
}
