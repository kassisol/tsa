package system

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/api/types"
	"github.com/kassisol/tsa/cli/session"
	"github.com/kassisol/tsa/client"
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

	info, err := clt.GetInfo(srv.Token)
	if err != nil {
		log.Fatal(err)
	}

	if info.CA != (types.CertificationAuthority{}) {
		fmt.Println("Certificate Authority:")
		fmt.Println(" Type:", info.CA.Type)
		fmt.Println(" Expire:", info.CA.Expire)
		fmt.Println(" Country:", info.CA.Country)
		fmt.Println(" State:", info.CA.State)
		fmt.Println(" Locality:", info.CA.Locality)
		fmt.Println(" Organization:", info.CA.Organization)
		fmt.Println(" Organizational Unit:", info.CA.OrganizationalUnit)
		fmt.Println(" Common Name:", info.CA.CommonName)
		fmt.Println("Certificates:", info.CertificateStats.Certificate)
		fmt.Println(" Valid:", info.CertificateStats.Valid)
		fmt.Println(" Expired:", info.CertificateStats.Expired)
		fmt.Println(" Revoked:", info.CertificateStats.Revoked)
	}

	fmt.Println("API:")
	fmt.Println(" FQDN:", info.API.FQDN)
	fmt.Println(" Bind Address:", info.API.BindAddress)
	fmt.Println(" Bind Port:", info.API.BindPort)

	fmt.Println("Auth Type:", info.Auth.Type)

	fmt.Println("Server Version:", info.ServerVersion)
	fmt.Println("Storage Driver:", info.StorageDriver)
	fmt.Println("Logging Driver:", info.LoggingDriver)
	fmt.Println("TSA Root Dir:", info.TSARootDir)
}

var infoDescription = `
Display information about TSA

`
