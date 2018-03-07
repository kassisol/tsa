package cert

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-utils"
	"github.com/kassisol/tsa/cli/session"
	"github.com/kassisol/tsa/client"
	"github.com/spf13/cobra"
)

var certListFilter []string

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls [tenant_id]",
		Aliases: []string{"list"},
		Short:   "List certificates issued",
		Long:    listDescription,
		Run:     runList,
	}

	flags := cmd.Flags()
	flags.StringSliceVarP(&certListFilter, "filter", "f", []string{}, "Filter output based on conditions provided")

	return cmd
}

func runList(cmd *cobra.Command, args []string) {
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

	tenantID, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(err)
	}

	clt, err := client.New(srv.Server.TSAURL)
	if err != nil {
		log.Fatal(err)
	}

	filters := utils.ConvertSliceToMap("=", certListFilter)

	certificates, err := clt.CertList(srv.Token, tenantID, filters)
	if err != nil {
		log.Fatal(err)
	}

	if len(certificates) > 0 {
		tw := tabwriter.NewWriter(os.Stdout, 20, 1, 2, ' ', 0)
		fmt.Fprintln(tw, "STATUS FLAG\tEXPIRATION DATE\tREVOCATION DATE\tREVOCATION REASON\tSERIAL NUMBER\tFILENAME\tDISTINGUISHED NAME")

		for _, certificate := range certificates {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%d\t%s\t%s\n", certificate.StatusFlag, certificate.ExpirationDate, certificate.RevocationDate, certificate.RevocationReason, certificate.SerialNumber, certificate.Filename, certificate.DistinguishedName)
		}

		tw.Flush()
	}
}

var listDescription = `
List certificates issued

`
