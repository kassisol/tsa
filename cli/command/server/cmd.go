package server

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Manage TSA servers",
		Long:  serverDescription,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	cmd.AddCommand(
		newAddCommand(),
		newListCommand(),
		newRemoveCommand(),
		//newStatusCommand(),
		//newUseCommand(),
	)

	return cmd
}

var serverDescription = `
The **tsa server** command has subcommands for managing TSA servers.

To see help for a subcommand, use:

    tsa server [command] --help

For full details on using tsa server visit Harbormaster's online documentation.

`
