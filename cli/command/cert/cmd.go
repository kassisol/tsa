package cert

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cert",
		Short: "Manage certificates",
		Long:  certDescription,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	cmd.AddCommand(
		newListCommand(),
		newRevokeCommand(),
	)

	return cmd
}

var certDescription = `
The **tsa cert** command has subcommands for managing certificates.

To see help for a subcommand, use:

    tsa cert [command] --help

For full details on using tsa cert visit Harbormaster's online documentation.

`
