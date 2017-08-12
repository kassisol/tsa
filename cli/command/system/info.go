package system

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-cert/ca/database"
	"github.com/juliengk/go-cert/helpers"
	"github.com/juliengk/go-cert/pkix"
	"github.com/kassisol/tsa/api/config"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/version"
	"github.com/spf13/cobra"
)

func NewInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Display information about TSA",
		Long:  infoDescription,
		Run:   runInfo,
	}

	return cmd
}

func runInfo(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		cmd.Usage()
		os.Exit(-1)
	}

	s, err := storage.NewDriver("sqlite", config.AppPath)
	if err != nil {
		log.Fatal(err)
	}
	defer s.End()

	if len(s.ListConfigs("ca")) > 0 {
		db, err := database.NewBackend("sqlite", config.CaDir)
		if err != nil {
			log.Fatal(err)
		}
		defer db.End()

		expire := "0000-00-00"

		certificate, err := pkix.NewCertificateFromPEMFile(config.CaCrtFile)
		if err == nil {
			expire = helpers.ExpireDateString(certificate.Crt.NotAfter)
		}

		fmt.Println("Certificate Authority:")
		fmt.Println(" Type:", s.GetConfig("ca_type")[0].Value)
		fmt.Println(" Expire:", expire)
		fmt.Println(" Country:", s.GetConfig("ca_country")[0].Value)
		fmt.Println(" State:", s.GetConfig("ca_state")[0].Value)
		fmt.Println(" Locality:", s.GetConfig("ca_locality")[0].Value)
		fmt.Println(" Organization:", s.GetConfig("ca_org")[0].Value)
		fmt.Println(" Organizational Unit:", s.GetConfig("ca_ou")[0].Value)
		fmt.Println(" Common Name:", s.GetConfig("ca_cn")[0].Value)
		fmt.Println(" E-mail:", s.GetConfig("ca_email")[0].Value)
		fmt.Println("Certificates:", db.Count("A"))
		fmt.Println(" Valid:", db.Count("V"))
		fmt.Println(" Expired:", db.Count("E"))
		fmt.Println(" Revoked:", db.Count("R"))
	}

	fmt.Println("API:")
	fmt.Println(" FQDN:", s.GetConfig("api_fqdn")[0].Value)
	fmt.Println(" Bind Address:", s.GetConfig("api_bind")[0].Value)
	fmt.Println(" Bind Port:", s.GetConfig("api_port")[0].Value)

	fmt.Println("Auth Type:", s.GetConfig("auth_type")[0].Value)

	fmt.Println("Server Version:", version.Version)
	fmt.Println("Storage Driver: sqlite")
	fmt.Println("Logging Driver: standard")
	fmt.Println("TSA Root Dir:", config.AppPath)
}

var infoDescription = `
Display information about TSA

`
