package auth

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/cli/session"
	"github.com/kassisol/tsa/client"
	"github.com/spf13/cobra"
)

func newRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm [key] [value]",
		Aliases: []string{"remove"},
		Short:   "Remove authentication configuration",
		Long:    removeDescription,
		Run:     runRemove,
	}

	return cmd
}

func runRemove(cmd *cobra.Command, args []string) {
	if len(args) < 1 || len(args) > 2 {
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

	if len(args) == 1 {
		args = append(args, "")
	}

	if err := clt.AuthDelete(srv.Token, args[0], args[1]); err != nil {
		log.Fatal(err)
	}
}

var removeDescription = `
Remove authentication configuration

`
