package commands

import (
	"github.com/kassisol/tsa/cli/command/admin"
	"github.com/kassisol/tsa/cli/command/auth"
	"github.com/kassisol/tsa/cli/command/cert"
	"github.com/kassisol/tsa/cli/command/server"
	"github.com/kassisol/tsa/cli/command/system"
	"github.com/spf13/cobra"
)

func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		admin.NewCommand(),
		auth.NewCommand(),
		cert.NewCommand(),
		server.NewCommand(),
		system.NewInfoCommand(),
		system.NewInitCommand(),
		system.NewVersionCommand(),
	)
}
