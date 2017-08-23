package auth

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/cli/session"
	"github.com/kassisol/tsa/client"
	"github.com/spf13/cobra"
)

func newDisableCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable authentication",
		Long:  disableDescription,
		Run:   runDisable,
	}

	return cmd
}

func runDisable(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		cmd.Usage()
		os.Exit(-1)
	}

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.End()

	srv, err := sess.Get()
	if err != nil {
		log.Fatal(err)
	}

	clt, err := client.New(srv.Server.TSAURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := clt.AuthDisable(srv.Token); err != nil {
		log.Fatal(err)
	}
}

var disableDescription = `
Disable authentication

`
