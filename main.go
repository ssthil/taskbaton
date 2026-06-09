// Command taskbaton manages work-state handover between agentic tool sessions.
package main

import "github.com/ssthil/taskbaton/internal/cli"

// version is injected at build time via -ldflags "-X main.version=$(VERSION)".
var version = "dev"

func main() {
	cli.SetVersion(version)
	cli.Execute()
}
