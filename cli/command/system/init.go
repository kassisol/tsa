package system

import (
	"crypto/md5"
	"fmt"
	"os"
	"path"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-cert/ca"
	"github.com/juliengk/go-cert/ca/database"
	"github.com/juliengk/go-cert/helpers"
	"github.com/juliengk/go-cert/pkix"
	"github.com/juliengk/go-utils/filedir"
	"github.com/juliengk/go-utils/random"
	"github.com/juliengk/go-utils/readinput"
	"github.com/juliengk/go-utils/validation"
	"github.com/kassisol/tsa/cli/command"
	clivalidation "github.com/kassisol/tsa/cli/validation"
	"github.com/kassisol/tsa/storage"
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
	serverAPIFQDN  string
	serverAPIBind  string
	serverAPIPort  string
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

	flags.StringVarP(&serverAPIFQDN, "api-fqdn", "", "", "API FQDN")
	flags.StringVarP(&serverAPIBind, "api-bind", "", "0.0.0.0", "API Bind Interface")
	flags.StringVarP(&serverAPIPort, "api-port", "", "443", "API Port")

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
	var apifqdn string
	var apibind string
	var apiport string

	if len(args) > 0 {
		cmd.Usage()
		os.Exit(-1)
	}

	if filedir.FileExists(command.DBFilePath) {
		log.Info("Initialization already done")
		os.Exit(0)
	}

	err := filedir.CreateDirIfNotExist(command.AppPath, 0700)
	if err != nil {
		log.Fatal(err)
	}

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

	if len(serverAPIFQDN) <= 0 {
		apifqdn = readinput.ReadInput("API FQDN")
	} else {
		apifqdn = serverAPIFQDN
	}

	if len(serverAPIBind) <= 0 {
		apibind = readinput.ReadInput("API Bind")
	} else {
		apibind = serverAPIBind
	}

	if len(serverAPIPort) <= 0 {
		apiport = readinput.ReadInput("API Port")
	} else {
		apiport = serverAPIPort
	}

	// DB
	s, err := storage.NewDriver("sqlite", command.DBFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer s.End()

	// Input validations
	// IV - CA Type
	if err = clivalidation.IsValidCAType(catype); err != nil {
		log.Fatal(err)
	}

	// IV - Organizational Unit
	if err = helpers.IsValidCAOrgUnit(ou); err != nil {
		log.Fatal(err)
	}

	// IV - E-mail
	if err = validation.IsValidEmail(email); err != nil {
		log.Fatal(err)
	}

	// IV - API FQDN
	if err = validation.IsValidFQDN(apifqdn); err != nil {
		log.Fatal(err)
	}

	// IV - API Bind
	if err = validation.IsValidIP(apibind); err != nil {
		log.Fatal(err)
	}

	// IV - API Port
	port, err := strconv.Atoi(apiport)
	if err != nil {
		log.Fatal(err)
	}
	if err = validation.IsValidPort(port); err != nil {
		log.Fatal(err)
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

	s.AddConfig("api_fqdn", apifqdn)
	s.AddConfig("api_bind", apibind)
	s.AddConfig("api_port", apiport)

	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()_-+="
	s.AddConfig("jwk", random.RandString(letterBytes, 24))

	s.AddConfig("auth_type", "none")

	// Initialize CA
	caSubject := pkix.NewSubject(country, state, locality, o, caou, cacn)

	caNDN := pkix.NewDNSNames()

	caNE := pkix.NewEmails()
	caNE.AddEmail(email)

	caIP := pkix.NewIPs()

	caAltnames := pkix.NewSubjectAltNames(caNDN, caNE, caIP)

	d, _ := strconv.Atoi(duration)
	caDate := ca.CreateDate(d)
	caSN := 1

	caTemplate, err := ca.CreateTemplate(true, caSubject, caAltnames, caDate, caSN, "")
	if err != nil {
		log.Fatal(err)
	}

	newCA, err := ca.InitCA(command.AppPath, caTemplate)
	if err != nil {
		log.Fatal(err)
	}

	ca.CreateCRLFile(command.CaCrlFile)

	// Create certificate for API
	apiDuration := 12

	err = filedir.CreateDirIfNotExist(command.ApiCertsDir, 0750)
	if err != nil {
		log.Fatal(err)
	}

	// Key pair
	apiKey, err := pkix.NewKey(4096)
	if err != nil {
		log.Fatal(err)
	}

	apiKeyBytes, err := apiKey.ToPEM()
	if err != nil {
		log.Fatal(err)
	}

	err = pkix.ToPEMFile(command.ApiKeyFile, apiKeyBytes, 0400)
	if err != nil {
		log.Fatal(err)
	}

	// CSR
	apiSubject := pkix.NewSubject(country, state, locality, o, ou, apifqdn)

	apiNDN := pkix.NewDNSNames()

	apiNE := pkix.NewEmails()
	apiNE.AddEmail(email)

	apiIP := pkix.NewIPs()

	apiAltnames := pkix.NewSubjectAltNames(apiNDN, apiNE, apiIP)

	apiCsr, err := pkix.NewCertificateRequest(apiKey, apiSubject, apiAltnames)
	if err != nil {
		log.Fatal(err)
	}

	// CA
	/*newCA, err := ca.NewCA(command.AppPath)
	if err != nil {
		log.Fatal(err)
	}*/

	apiPubKey := apiCsr.GetPublicKey()
	apiSubjectAltNames := apiCsr.GetSubjectAltNames()
	apiDate := ca.CreateDate(apiDuration)
	apiSN, err := newCA.IncrementSerialNumber()
	if err != nil {
		log.Fatal(err)
	}

	apiTemplate, err := ca.CreateTemplate(false, apiSubject, apiSubjectAltNames, apiDate, apiSN, "")
	if err != nil {
		log.Fatal(err)
	}

	apiDerBytes, err := ca.IssueCertificate(apiTemplate, newCA.Certificate.Crt, apiPubKey, newCA.Key.Private)
	if err != nil {
		log.Fatal(err)
	}

	// create certificate PEM file
	apiCertificate, err := pkix.NewCertificateFromDER(apiDerBytes)
	if err != nil {
		log.Fatal(err)
	}

	apiCrtBytes, err := apiCertificate.ToPEM()
	if err != nil {
		log.Fatal(err)
	}

	err = pkix.ToPEMFile(command.ApiCrtFile, apiCrtBytes, 0444)
	if err != nil {
		log.Fatal(err)
	}

	// create serial number file
	newCA.WriteSerialNumber(apiSN)

	// save information to CA database
	db, err := database.NewBackend("sqlite", command.CaDir)
	if err != nil {
		log.Fatal(err)
	}
	defer db.End()

	indexExpireDate := ca.DatabaseDateTimeFormat(apiDate.Expire)
	indexDN := apiCsr.SubjectToString()

	certName := fmt.Sprintf("%x", md5.Sum([]byte(indexDN)))
	certNameFile := fmt.Sprintf("%s.crt", certName)
	certNamePath := path.Join(command.CaCertsDir, certNameFile)

	db.New(apiSN, indexExpireDate, certNameFile, indexDN)

	err = pkix.ToPEMFile(certNamePath, apiCrtBytes, 0444)
	if err != nil {
		log.Fatal(err)
	}
}

var initDescription = `
Initialize config

`
