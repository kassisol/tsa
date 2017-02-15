package system

import (
	"fmt"

	"github.com/kassisol/tsa/version"
	"github.com/spf13/cobra"
)

func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show the TSA version information",
		Long:  versionDescription,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("TSA", version.Version, "-- HEAD")
		},
	}

	return cmd
}

var versionDescription = `
All software has versions. This is TSA's

`
