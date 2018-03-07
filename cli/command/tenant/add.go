package tenant

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-utils"
	"github.com/juliengk/go-utils/readinput"
	"github.com/kassisol/tsa/cli/session"
	"github.com/kassisol/tsa/client"
	"github.com/spf13/cobra"
)

var (
	tenantAuthGroups []string
	tenantDuration   int
	tenantCountry    string
	tenantState      string
	tenantLocality   string
	tenantOrg        string
	tenantOrgUnit    string
)

func newCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [tenant name]",
		Short: "Add tenant",
		Long:  addDescription,
		Run:   runAdd,
	}

	flags := cmd.Flags()
	flags.StringSliceVarP(&tenantAuthGroups, "auth-group", "", []string{}, "Auth Groups")
	flags.IntVarP(&tenantDuration, "duration", "", 120, "Duration")
	flags.StringVarP(&tenantCountry, "country", "", "", "Country")
	flags.StringVarP(&tenantState, "state", "", "", "State")
	flags.StringVarP(&tenantLocality, "city", "", "", "Locality")
	flags.StringVarP(&tenantOrg, "org", "", "", "Organization")
	flags.StringVarP(&tenantOrgUnit, "org-unit", "", "", "Organizational Unit")

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) {
	var authGroups []string
	var caType string
	var duration int
	var country string
	var state string
	var locality string
	var o string
	var ou string

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

	caType = "root"
	duration = tenantDuration

	if len(tenantAuthGroups) <= 0 {
		ag := readinput.ReadInput("Auth Groups")

		authGroups = utils.CreateSlice(ag, ",")
	} else {
		authGroups = tenantAuthGroups
	}

	if len(tenantCountry) <= 0 {
		country = readinput.ReadInput("Country")
	} else {
		country = tenantCountry
	}

	if len(tenantState) <= 0 {
		state = readinput.ReadInput("State")
	} else {
		state = tenantState
	}

	if len(tenantLocality) <= 0 {
		locality = readinput.ReadInput("City")
	} else {
		locality = tenantLocality
	}

	if len(tenantOrg) <= 0 {
		o = readinput.ReadInput("Organization (O)")
	} else {
		o = tenantOrg
	}

	if len(tenantOrgUnit) <= 0 {
		ou = readinput.ReadInput("Organizational Unit (OU)")
	} else {
		ou = tenantOrgUnit
	}

	clt, err := client.New(srv.Server.TSAURL)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := clt.TenantCreate(srv.Token, args[0], authGroups, caType, country, state, locality, o, ou, duration); err != nil {
		log.Fatal(err)
	}
}

var addDescription = `
Add tenant

`
