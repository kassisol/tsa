package access

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "access",
		Short: "Manage session to TSA servers",
		Long:  accessDescription,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	cmd.AddCommand(
		newLoginCommand(),
		newListCommand(),
		newRemoveCommand(),
		newStatusCommand(),
		newUseCommand(),
		newUnuseCommand(),
	)

	return cmd
}

var accessDescription = `
The **tsa access** command has subcommands for managing session to TSA servers.

To see help for a subcommand, use:

    tsa access [command] --help

For full details on using tsa access visit Harbormaster's online documentation.

`
