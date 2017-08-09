package server

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-utils/filedir"
	"github.com/juliengk/go-utils/validation"
	"github.com/kassisol/tsa/api"
	"github.com/kassisol/tsa/cli/command"
	"github.com/kassisol/tsa/pkg/tls"
	"github.com/kassisol/tsa/pkg/token"
	"github.com/kassisol/tsa/storage"
	"github.com/spf13/cobra"
)

var (
	serverBindAddress string
	serverBindPort    int
	serverTLS         bool
	serverTLSCN       string
	serverTLSDuration int
	serverTLSGen      bool
	serverTLSCert     string
	serverTLSKey      string
)

func newStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts a CA (Certificate Authority) server",
		Long:  startDescription,
		Run:   runStart,
	}

	flags := cmd.Flags()
	flags.StringVarP(&serverBindAddress, "bind-address", "a", "0.0.0.0", "Bind Address")
	flags.IntVarP(&serverBindPort, "bind-port", "p", 80, "Bind Port")
	flags.BoolVarP(&serverTLS, "tls", "t", false, "Enable TLS certificates")
	flags.StringVarP(&serverTLSCN, "tlscn", "", "", "Certificate Common Name")
	flags.IntVarP(&serverTLSDuration, "tls-duration", "", 60, "Certificate duration")
	flags.BoolVarP(&serverTLSGen, "tlsgen", "", false, "Auto generate self-signed TLS certificates")
	flags.StringVarP(&serverTLSCert, "tlscert", "", command.ApiCrtFile, "Path to TLS certificate file")
	flags.StringVarP(&serverTLSKey, "tlskey", "", command.ApiKeyFile, "Path to TLS key file")

	return cmd
}

func serverInitConfig() error {
	if filedir.FileExists(command.DBFilePath) {
		log.Info("Server initialization already done")

		return nil
	}

	if err := filedir.CreateDirIfNotExist(command.AppPath, false, 0700); err != nil {
		return err
	}

	if err := filedir.CreateDirIfNotExist(command.ApiCertsDir, false, 0750); err != nil {
		return err
	}

	// DB
	s, err := storage.NewDriver("sqlite", command.DBFilePath)
	if err != nil {
		return err
	}
	defer s.End()

	s.AddConfig("jwk", token.GenerateJWK("", 24))
	s.AddConfig("auth_type", "none")

	return nil
}

func runStart(cmd *cobra.Command, args []string) {
	var bindPort int
	var tlsCN string

	if len(args) > 0 {
		cmd.Usage()
		os.Exit(-1)
	}

	if err := serverInitConfig(); err != nil {
		log.Fatal(err)
	}

	bindPort = serverBindPort
	if serverTLS && bindPort == 80 {
		bindPort = 443
	}

	if len(serverTLSCN) <= 0 {
		af, err := os.Hostname()
		if err != nil {
			log.Fatal(err)
		}
		tlsCN = af
	} else {
		tlsCN = serverTLSCN
	}

	// Input validations
	// IV - API Bind address
	if err := validation.IsValidIP(serverBindAddress); err != nil {
		log.Fatal(err)
	}

	// IV - API Port
	if err := validation.IsValidPort(bindPort); err != nil {
		log.Fatal(err)
	}

	// IV - API FQDN
	if err := validation.IsValidFQDN(tlsCN); err != nil {
		log.Fatal(err)
	}

	// Create API certificates
	config, err := tls.New(tlsCN, serverTLSDuration, serverTLSKey, serverTLSCert)
	if err != nil {
		log.Fatal(err)
	}

	if serverTLSGen {
		if !config.CertificateExist() || (config.CertificateExist() && config.IsCertificateExpire()) {
			if err := config.CreateSelfSignedCertificate(); err != nil {
				log.Fatal(err)
			}
		}
	}

	if serverTLS {
		if !config.CertificateExist() {
			log.Fatal("No certificate found")
		}
	}

	addr := fmt.Sprintf("%s:%d", serverBindAddress, bindPort)

	api.API(addr, serverTLS)
}

var startDescription = `
Starts a CA (Certificate Authority) server

`
