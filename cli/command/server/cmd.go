package server

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Starts a CA (Certificate Authority) server",
		Long:  serverDescription,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	cmd.AddCommand(
		newStartCommand(),
	)

	return cmd
}

var serverDescription = `
The **tsa server** command has subcommands for starting a CA (Certificate Authority) server.

To see help for a subcommand, use:

    tsa server [command] --help

For full details on using tsa server visit Harbormaster's online documentation.

`

var serverConfig string
