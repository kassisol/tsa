package system

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-utils/readinput"
	"github.com/kassisol/tsa/cli/session"
	"github.com/kassisol/tsa/client"
	"github.com/spf13/cobra"
)

var (
	serverDuration int
	serverCountry  string
	serverState    string
	serverLocality string
	serverOrg      string
	serverOrgUnit  string
)

func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize config",
		Long:  initDescription,
		Run:   runInit,
	}

	flags := cmd.Flags()
	flags.IntVarP(&serverDuration, "duration", "", 120, "Duration")
	flags.StringVarP(&serverCountry, "country", "", "", "Country")
	flags.StringVarP(&serverState, "state", "", "", "State")
	flags.StringVarP(&serverLocality, "city", "", "", "Locality")
	flags.StringVarP(&serverOrg, "org", "", "", "Organization")
	flags.StringVarP(&serverOrgUnit, "org-unit", "", "", "Organizational Unit")

	return cmd
}

func runInit(cmd *cobra.Command, args []string) {
	var catype string
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

	catype = "root"
	duration = serverDuration

	if len(serverCountry) <= 0 {
		country = readinput.ReadInput("Country")
	} else {
		country = serverCountry
	}

	if len(serverState) <= 0 {
		state = readinput.ReadInput("State")
	} else {
		state = serverState
	}

	if len(serverLocality) <= 0 {
		locality = readinput.ReadInput("City")
	} else {
		locality = serverLocality
	}

	if len(serverOrg) <= 0 {
		o = readinput.ReadInput("Organization (O)")
	} else {
		o = serverOrg
	}

	if len(serverOrgUnit) <= 0 {
		ou = readinput.ReadInput("Organizational Unit (OU)")
	} else {
		ou = serverOrgUnit
	}

	clt, err := client.New(srv.Server.TSAURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := clt.CAInit(srv.Token, catype, country, state, locality, o, ou, duration); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Initialization done successfully")
}

var initDescription = `
Initialize config

`
