package tenant

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-utils"
	"github.com/kassisol/tsa/cli/session"
	"github.com/kassisol/tsa/client"
	"github.com/spf13/cobra"
)

var tenantListFilter []string

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List tenants",
		Long:    listDescription,
		Run:     runList,
	}

	flags := cmd.Flags()
	flags.StringSliceVarP(&tenantListFilter, "filter", "f", []string{}, "Filter output based on conditions provided")

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

	clt, err := client.New(srv.Server.TSAURL)
	if err != nil {
		log.Fatal(err)
	}

	filters := utils.ConvertSliceToMap("=", tenantListFilter)

	tenants, err := clt.TenantList(srv.Token, filters)
	if err != nil {
		log.Fatal(err)
	}

	if len(tenants) > 0 {
		tw := tabwriter.NewWriter(os.Stdout, 20, 1, 2, ' ', 0)
		fmt.Fprintln(tw, "ID\tNAME\tAUTH GROUPS")

		for _, tenant := range tenants {
			groups := []string{}

			for _, g := range tenant.AuthGroups {
				groups = append(groups, g.Name)
			}

			fmt.Fprintf(tw, "%d\t%s\t%s\n", tenant.ID, tenant.Name, strings.Join(groups, ", "))
		}

		tw.Flush()
	}
}

var listDescription = `
List tenants

`
