package main

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-utils/filedir"
	"github.com/juliengk/go-utils/password"
	"github.com/juliengk/go-utils/validation"
	"github.com/kassisol/tsa/api"
	"github.com/kassisol/tsa/api/config"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/pkg/tls"
	"github.com/kassisol/tsa/pkg/token"
	"github.com/spf13/cobra"
)

func serverInitConfig() error {
	if filedir.FileExists(config.AppPath) {
		log.Info("Server initialization already done")

		return nil
	}

	if err := filedir.CreateDirIfNotExist(config.AppPath, false, 0700); err != nil {
		return err
	}

	if err := filedir.CreateDirIfNotExist(config.ApiCertsDir, false, 0750); err != nil {
		return err
	}

	s, err := storage.NewDriver("sqlite", config.AppPath)
	if err != nil {
		return err
	}
	defer s.End()

	s.AddConfig("jwk", token.GenerateJWK("", 24))
	s.AddConfig("auth_type", "none")
	s.AddConfig("admin_password", password.GeneratePassword("admin"))

	return nil
}

func runDaemon(cmd *cobra.Command, args []string) {
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

	// Create API certificates
	conf := tls.New(serverTLSKey, serverTLSCert)

	if serverTLSGen {
		if !conf.CertificateExist() || (conf.CertificateExist() && conf.IsCertificateExpire()) {
			if err := conf.CreateSelfSignedCertificate(tlsCN, serverTLSDuration); err != nil {
				log.Fatal(err)
			}
		}
	}

	if serverTLS {
		if !conf.CertificateExist() {
			log.Fatal("No certificate found")
		}
	}

	s, err := storage.NewDriver("sqlite", config.AppPath)
	if err != nil {
		log.Fatal(err)
	}
	defer s.End()

	s.RemoveConfig("api_bind", "ALL")
	s.RemoveConfig("api_port", "ALL")
	s.RemoveConfig("api_fqdn", "ALL")
	s.AddConfig("api_bind", serverBindAddress)
	s.AddConfig("api_port", strconv.Itoa(bindPort))
	s.AddConfig("api_fqdn", tlsCN)

	addr := fmt.Sprintf("%s:%d", serverBindAddress, bindPort)

	api.API(addr, serverTLS)
}
