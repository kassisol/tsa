package system

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
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

	fmt.Println("API:")
	fmt.Println(" FQDN:", info.API.FQDN)
	fmt.Println(" Bind Address:", info.API.BindAddress)
	fmt.Println(" Bind Port:", info.API.BindPort)

	fmt.Println("Auth Type:", info.Auth.Type)

	fmt.Println("Server Version:", info.ServerVersion)
	fmt.Println("ID:", info.ID)
	fmt.Println("Storage Driver:", info.StorageDriver)
	fmt.Println("Logging Driver:", info.LoggingDriver)
	fmt.Println("TSA Root Dir:", info.TSARootDir)
}

var infoDescription = `
Display information about TSA

`
