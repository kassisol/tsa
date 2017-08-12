package cert

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-cert/ca/database"
	"github.com/juliengk/go-utils"
	"github.com/kassisol/tsa/api/config"
	"github.com/spf13/cobra"
)

var certListFilter []string

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
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

	db, err := database.NewBackend("sqlite", config.CaDir)
	if err != nil {
		log.Fatal(err)
	}
	defer db.End()

	filters := utils.ConvertSliceToMap("=", certListFilter)

	certificates := db.List(filters)

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
