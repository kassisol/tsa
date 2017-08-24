package cert

import (
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/cli/session"
	"github.com/kassisol/tsa/client"
	"github.com/spf13/cobra"
)

func newRevokeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke [serial number]",
		Short: "Revoke certificate",
		Long:  revokeDescription,
		Run:   runRevoke,
	}

	return cmd
}

func runRevoke(cmd *cobra.Command, args []string) {
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

	serialNumber, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(err)
	}

	if err := clt.CertRevoke(srv.Token, serialNumber); err != nil {
		log.Fatal(err)
	}
}

var revokeDescription = `
Revoke certificate

`
