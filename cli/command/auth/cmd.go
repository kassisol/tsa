package auth

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication informations",
		Long:  authDescription,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	cmd.AddCommand(
		newEnableCommand(),
		newDisableCommand(),
		newListCommand(),
		newAddCommand(),
		newRemoveCommand(),
	)

	return cmd
}

var authDescription = `
The **tsa auth** command has subcommands for managing authentication configurations.

To see help for a subcommand, use:

    tsa auth [command] --help

For full details on using tsa auth visit Harbormaster's online documentation.

`
