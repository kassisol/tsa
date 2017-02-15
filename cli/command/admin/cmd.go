package admin

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admin",
		Short: "Manage TSA config",
		Long:  adminDescription,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	cmd.AddCommand(
		newBindCommand(),
	)

	return cmd
}

var adminDescription = `
The **tsa admin** command has subcommands for managing TSA config.

To see help for a subcommand, use:

    tsa admin [command] --help

For full details on using tsa admin visit Harbormaster's online documentation.

`
