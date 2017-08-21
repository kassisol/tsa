package system

import (
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-cert/ca"
	"github.com/juliengk/go-cert/helpers"
	"github.com/juliengk/go-cert/pkix"
	"github.com/juliengk/go-utils/readinput"
	"github.com/juliengk/go-utils/validation"
	"github.com/kassisol/tsa/api/storage"
	clivalidation "github.com/kassisol/tsa/cli/validation"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/spf13/cobra"
)

var (
	serverDuration string
	serverCountry  string
	serverState    string
	serverLocality string
	serverOrg      string
	serverOrgUnit  string
	serverEmail    string
)

func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize config",
		Long:  initDescription,
		Run:   runInit,
	}

	flags := cmd.Flags()
	flags.StringVarP(&serverDuration, "duration", "", "120", "Duration")
	flags.StringVarP(&serverCountry, "country", "", "", "Country")
	flags.StringVarP(&serverState, "state", "", "", "State")
	flags.StringVarP(&serverLocality, "city", "", "", "Locality")
	flags.StringVarP(&serverOrg, "org", "", "", "Organization")
	flags.StringVarP(&serverOrgUnit, "org-unit", "", "", "Organizational Unit")
	flags.StringVarP(&serverEmail, "email", "", "", "E-mail")

	return cmd
}

func runInit(cmd *cobra.Command, args []string) {
	var catype string
	var duration string
	var country string
	var state string
	var locality string
	var o string
	var ou string
	var email string

	if len(args) > 0 {
		cmd.Usage()
		os.Exit(-1)
	}

	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		log.Fatal(err)
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		panic(err)
	}
	defer s.End()

	if len(s.ListConfigs("ca")) > 0 {
		log.Fatal("Initialization already done")
	}

	defer func() {
		if r := recover(); r != nil {
			os.RemoveAll(cfg.CA.Dir.Root)

			s.RemoveConfig("ca_type", "ALL")
			s.RemoveConfig("ca_duration", "ALL")
			s.RemoveConfig("ca_country", "ALL")
			s.RemoveConfig("ca_state", "ALL")
			s.RemoveConfig("ca_locality", "ALL")
			s.RemoveConfig("ca_org", "ALL")
			s.RemoveConfig("ca_ou", "ALL")
			s.RemoveConfig("ca_cn", "ALL")
			s.RemoveConfig("ca_email", "ALL")

			log.Fatal(r)
		}
	}()

	catype = "root"

	if len(serverDuration) <= 0 {
		duration = readinput.ReadInput("Duration")
	} else {
		duration = serverDuration
	}

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

	if len(serverEmail) <= 0 {
		email = readinput.ReadInput("Email")
	} else {
		email = serverEmail
	}

	// Input validations
	// IV - CA Type
	if err := clivalidation.IsValidCAType(catype); err != nil {
		panic(err)
	}

	// IV - Organizational Unit
	if err := helpers.IsValidCAOrgUnit(ou); err != nil {
		panic(err)
	}

	// IV - E-mail
	if err := validation.IsValidEmail(email); err != nil {
		panic(err)
	}

	// Save inputs to DB
	s.AddConfig("ca_type", catype)
	s.AddConfig("ca_duration", duration)

	caou := helpers.UpdateOrgUnitLabel(ou)
	cacn := helpers.UpdateCommonNameLabel("root", ou)

	s.AddConfig("ca_country", country)
	s.AddConfig("ca_state", state)
	s.AddConfig("ca_locality", locality)
	s.AddConfig("ca_org", o)
	s.AddConfig("ca_ou", caou)
	s.AddConfig("ca_cn", cacn)
	s.AddConfig("ca_email", email)

	// Initialize CA
	caSubject := pkix.NewSubject(country, state, locality, o, caou, cacn)

	caNDN := pkix.NewDNSNames()

	caNE := pkix.NewEmails()
	caNE.AddEmail(email)

	caIP := pkix.NewIPs()

	caAltnames := pkix.NewSubjectAltNames(*caNDN, *caNE, *caIP)

	d, _ := strconv.Atoi(duration)
	caDate := ca.CreateDate(d)
	caSN := 1

	caTemplate, err := ca.CreateTemplate(true, caSubject, caAltnames, caDate, caSN, "")
	if err != nil {
		panic(err)
	}

	if _, err = ca.InitCA(cfg.App.Dir.Root, caTemplate); err != nil {
		panic(err)
	}

	ca.CreateCRLFile(cfg.CA.CRLFile)
}

var initDescription = `
Initialize config

`
