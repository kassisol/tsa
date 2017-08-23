package auth

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/cli/session"
	"github.com/kassisol/tsa/client"
	"github.com/spf13/cobra"
)

func newAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [key] [value]",
		Short: "Add auth configuration",
		Long:  addDescription,
		Run:   runAdd,
	}

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) {
	if len(args) < 2 || len(args) > 2 {
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

	if _, err := clt.AuthCreate(srv.Token, args[0], args[1]); err != nil {
		log.Fatal(err)
	}
}

var addDescription = `
Add auth configuration

`
