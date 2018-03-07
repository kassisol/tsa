package daemon

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/spf13/cobra"
)

var (
	serverBindAddress string
	serverBindPort    int
	serverFQDN        string
	serverTLS         bool
	serverTLSDuration int
	serverTLSGen      bool
	serverTLSCert     string
	serverTLSKey      string
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tsad",
		Short: "Starts a CA (Certification Authority) server",
		Long:  daemonDescription,
		Run:   runDaemon,
	}

	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		log.Fatal(err)
	}

	flags := cmd.Flags()
	flags.StringVarP(&serverBindAddress, "bind-address", "a", "0.0.0.0", "Bind Address")
	flags.IntVarP(&serverBindPort, "bind-port", "p", 80, "Bind Port")
	flags.StringVarP(&serverFQDN, "fqdn", "f", "", "API FQDN")
	flags.BoolVarP(&serverTLS, "tls", "t", false, "Enable TLS certificates")
	flags.IntVarP(&serverTLSDuration, "tls-duration", "", 60, "Certificate duration")
	flags.BoolVarP(&serverTLSGen, "tlsgen", "", false, "Auto generate self-signed TLS certificates")
	flags.StringVarP(&serverTLSCert, "tlscert", "", cfg.API.CrtFile, "Path to TLS certificate file")
	flags.StringVarP(&serverTLSKey, "tlskey", "", cfg.API.KeyFile, "Path to TLS key file")

	return cmd
}

var daemonDescription = `
Starts a CA (Certification Authority) server

`
