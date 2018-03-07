package tenant

import (
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/cli/session"
	"github.com/kassisol/tsa/client"
	"github.com/spf13/cobra"
)

func newDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm [tenant_id]",
		Aliases: []string{"remove"},
		Short:   "Remove tenant",
		Long:    removeDescription,
		Run:     runRemove,
	}

	return cmd
}

func runRemove(cmd *cobra.Command, args []string) {
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

	tid, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(err)
	}

	if err := clt.TenantDelete(srv.Token, tid); err != nil {
		log.Fatal(err)
	}
}

var removeDescription = `
Remove tenant

`
