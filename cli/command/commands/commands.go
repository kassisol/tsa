package commands

import (
	"github.com/kassisol/tsa/cli/command/auth"
	"github.com/kassisol/tsa/cli/command/cert"
	"github.com/kassisol/tsa/cli/command/system"
	"github.com/spf13/cobra"
)

func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		auth.NewCommand(),
		cert.NewCommand(),
		system.NewInfoCommand(),
		system.NewInitCommand(),
		system.NewPasswdCommand(),
		system.NewVersionCommand(),
	)
}
