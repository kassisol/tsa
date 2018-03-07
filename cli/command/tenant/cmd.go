package tenant

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tenant",
		Short: "Manage tenants",
		Long:  tenantDescription,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	cmd.AddCommand(
		newListCommand(),
		newCreateCommand(),
		newDeleteCommand(),
		newMemberCommand(),
	)

	return cmd
}

var tenantDescription = `
The **tsa tenant** command has subcommands for managing tenants.

To see help for a subcommand, use:

    tsa tenant [command] --help

For full details on using tsa tenant visit Harbormaster's online documentation.

`
