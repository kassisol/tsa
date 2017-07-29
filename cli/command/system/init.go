package system

import (
	"crypto/md5"
	"fmt"
	"net"
	"os"
	"path"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-cert/ca"
	"github.com/juliengk/go-cert/ca/database"
	"github.com/juliengk/go-cert/helpers"
	"github.com/juliengk/go-cert/pkix"
	"github.com/juliengk/go-utils/filedir"
	"github.com/juliengk/go-utils/ip"
	"github.com/juliengk/go-utils/readinput"
	"github.com/juliengk/go-utils/validation"
	"github.com/kassisol/tsa/cli/command"
	clivalidation "github.com/kassisol/tsa/cli/validation"
	"github.com/kassisol/tsa/pkg/token"
	"github.com/kassisol/tsa/storage"
	"github.com/spf13/cobra"
)

var (
	serverDuration         string
	serverCountry          string
	serverState            string
	serverLocality         string
	serverOrg              string
	serverOrgUnit          string
	serverEmail            string
	serverAPIFQDN          string
	serverAPIBind          string
	serverAPIPort          string
)

type configCommon struct {
	Country  string
	State    string
	Locality string
	O        string
	OU       string
	Email    string
}

type configCustom struct {
	KeyFile  string
	CrtFile  string
	FQDN     string
	Port     string
	IP       []string
	Duration int
}

type configCert struct {
	configCommon
	configCustom
}

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

	if err := filedir.CreateDirIfNotExist(command.AppPath, false, 0700); err != nil {
		log.Fatal(err)
	}

	defer func() {
		if r := recover(); r != nil {
			os.RemoveAll(command.AppPath)

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

	if len(serverAPIFQDN) <= 0 {
		apifqdn = readinput.ReadInput("API FQDN")
		if len(apifqdn) <= 0 {
			af, err := os.Hostname()
			if err != nil {
				panic(err)
			}
			apifqdn = af
		}
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
		panic(err)
	}
	defer s.End()

	// Input validations
	// IV - CA Type
	if err = clivalidation.IsValidCAType(catype); err != nil {
		panic(err)
	}

	// IV - Organizational Unit
	if err = helpers.IsValidCAOrgUnit(ou); err != nil {
		panic(err)
	}

	// IV - E-mail
	if err = validation.IsValidEmail(email); err != nil {
		panic(err)
	}

	// IV - API FQDN
	if err = validation.IsValidFQDN(apifqdn); err != nil {
		panic(err)
	}

	// IV - API Bind
	if err = validation.IsValidIP(apibind); err != nil {
		panic(err)
	}

	// IV - API Port
	port, err := strconv.Atoi(apiport)
	if err != nil {
		panic(err)
	}
	if err = validation.IsValidPort(port); err != nil {
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

	s.AddConfig("api_fqdn", apifqdn)
	s.AddConfig("api_bind", apibind)
	s.AddConfig("api_port", apiport)

	s.AddConfig("jwk", token.GenerateJWK("", 24))

	s.AddConfig("auth_type", "none")

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

	newCA, err := ca.InitCA(command.AppPath, caTemplate)
	if err != nil {
		panic(err)
	}

	ca.CreateCRLFile(command.CaCrlFile)

	ccommon := configCommon{
		Country:  country,
		State:    state,
		Locality: locality,
		O:        o,
		OU:       ou,
		Email:    email,
	}

	// Create certificate for API
	if err = filedir.CreateDirIfNotExist(command.ApiCertsDir, false, 0750); err != nil {
		log.Fatal(err)
	}

	ips := []string{}

	interfaces := ip.New()
	interfaces.Get()
	interfaces.IgnoreIntf([]string{"lo"})

	for _, intf := range interfaces {
		if len(intf.V4) > 0 {
			ips = append(ips, intf.V4[0])
		}
	}

	apiccustom := configCustom{
		KeyFile:  command.ApiKeyFile,
		CrtFile:  command.ApiCrtFile,
		FQDN:     apifqdn,
		Port:     apiport,
		IP:       ips,
		Duration: 60,
	}

	apiconfig := configCert{
		ccommon,
		apiccustom,
	}

	if err = createCert(newCA, apiconfig); err != nil {
		panic(err)
	}
}

func createCert(newCA *ca.CA, config configCert) error {
	// Key pair
	key, err := pkix.NewKey(4096)
	if err != nil {
		return err
	}

	keyBytes, err := key.ToPEM()
	if err != nil {
		return err
	}

	if err = pkix.ToPEMFile(config.KeyFile, keyBytes, 0400); err != nil {
		return err
	}

	// CSR
	subject := pkix.NewSubject(config.Country, config.State, config.Locality, config.O, config.OU, config.FQDN)

	ndn := pkix.NewDNSNames()

	ne := pkix.NewEmails()
	ne.AddEmail(config.Email)

	ips := pkix.NewIPs()

	if len(config.IP) > 0 {
		for _, ip := range config.IP {
			ips.AddIP(net.ParseIP(ip))
		}
	}

	altnames := pkix.NewSubjectAltNames(*ndn, *ne, *ips)

	csr, err := pkix.NewCertificateRequest(key, subject, altnames)
	if err != nil {
		return err
	}

	// CA
	pubKey := csr.GetPublicKey()
	subjectAltNames := csr.GetSubjectAltNames()
	date := ca.CreateDate(config.Duration)
	sn, err := newCA.IncrementSerialNumber()
	if err != nil {
		return err
	}

	crlUrl := fmt.Sprintf("https://%s:%s/crl/CRL.crl", config.FQDN, config.Port)

	template, err := ca.CreateTemplate(false, subject, subjectAltNames, date, sn, crlUrl)
	if err != nil {
		return err
	}

	derBytes, err := ca.IssueCertificate(template, newCA.Certificate.Crt, pubKey, newCA.Key.Private)
	if err != nil {
		return err
	}

	// create certificate PEM file
	certificate, err := pkix.NewCertificateFromDER(derBytes)
	if err != nil {
		return err
	}

	crtBytes, err := certificate.ToPEM()
	if err != nil {
		return err
	}

	if err = pkix.ToPEMFile(config.CrtFile, crtBytes, 0444); err != nil {
		return err
	}

	// create serial number file
	newCA.WriteSerialNumber(sn)

	// save information to CA database
	db, err := database.NewBackend("sqlite", command.CaDir)
	if err != nil {
		return err
	}
	defer db.End()

	indexExpireDate := ca.DatabaseDateTimeFormat(date.Expire)
	indexDN := csr.SubjectToString()

	certName := fmt.Sprintf("%x", md5.Sum([]byte(indexDN)))
	certNameFile := fmt.Sprintf("%s.crt", certName)
	certNamePath := path.Join(command.CaCertsDir, certNameFile)

	db.New(sn, indexExpireDate, certNameFile, indexDN)

	if err = pkix.ToPEMFile(certNamePath, crtBytes, 0444); err != nil {
		return err
	}

	return nil
}

var initDescription = `
Initialize config

`
