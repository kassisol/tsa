package cert

import (
	"os"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-cert/ca"
	"github.com/juliengk/go-cert/ca/database"
	"github.com/kassisol/tsa/cli/command"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ocsp"
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

	serialNumber, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewBackend("sqlite", command.CaDir)
	if err != nil {
		log.Fatal(err)
	}
	defer db.End()

	// Revoke certificate
	revocationDate := ca.DatabaseDateTimeFormat(time.Now())
	revocationReason := ocsp.CessationOfOperation

	db.Revoke(serialNumber, revocationDate, revocationReason)
}

var revokeDescription = `
Revoke certificate

`
