package system

import (
	"os"

	"github.com/kassisol/tsa/cli/session"
	"github.com/kassisol/tsa/client"
	"github.com/kassisol/tsa/version"
	"github.com/spf13/cobra"
)

func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show the TSA version information",
		Long:  versionDescription,
		Run:   runVersion,
	}

	return cmd
}

func runVersion(cmd *cobra.Command, args []string) {
	sess, err := session.New()
	if err != nil {
		display(&version.VersionInfo{}, err.Error())
	}
	defer sess.End()

	srv, err := sess.Get()
	if err != nil {
		display(&version.VersionInfo{}, err.Error())
	}

	clt, err := client.New(srv.Server.TSAURL)
	if err != nil {
		display(&version.VersionInfo{}, err.Error())
	}

	server, err := clt.GetServerVersion()
	if err != nil {
		display(&version.VersionInfo{}, err.Error())
	}

	display(server, "")
}

func display(server *version.VersionInfo, srvErr string) {
	code := 0
	if len(srvErr) > 0 {
		code = 1
	}

	info := version.NewDisplay(server, srvErr)
	info.Show()

	os.Exit(code)
}

var versionDescription = `
All software has versions. This is TSA's

`
