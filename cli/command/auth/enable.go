package auth

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/cli/session"
	"github.com/kassisol/tsa/client"
	"github.com/spf13/cobra"
)

func newEnableCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable [type]",
		Short: "Enable authentication",
		Long:  enableDescription,
		Run:   runEnable,
	}

	return cmd
}

func runEnable(cmd *cobra.Command, args []string) {
	if len(args) < 1 || len(args) > 1 {
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

	if err := clt.AuthEnable(srv.Token, args[0]); err != nil {
		log.Fatal(err)
	}
}

var enableDescription = `
Enable authentication

`
