package commands

import (
	"github.com/kassisol/tsa/cli/command/access"
	"github.com/kassisol/tsa/cli/command/auth"
	"github.com/kassisol/tsa/cli/command/cert"
	"github.com/kassisol/tsa/cli/command/server"
	"github.com/kassisol/tsa/cli/command/system"
	"github.com/kassisol/tsa/cli/command/tenant"
	"github.com/spf13/cobra"
)

func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		access.NewCommand(),
		auth.NewCommand(),
		cert.NewCommand(),
		server.NewCommand(),
		system.NewInfoCommand(),
		system.NewPasswdCommand(),
		system.NewVersionCommand(),
		tenant.NewCommand(),
	)
}
