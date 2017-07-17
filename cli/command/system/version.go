package system

import (
	"github.com/kassisol/tsa/version"
	"github.com/spf13/cobra"
)

func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show the TSA version information",
		Long:  versionDescription,
		Run: func(cmd *cobra.Command, args []string) {
			info := version.New()
			info.ShowVersion()
		},
	}

	return cmd
}

var versionDescription = `
All software has versions. This is TSA's

`
