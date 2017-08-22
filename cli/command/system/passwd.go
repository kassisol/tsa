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

func NewPasswdCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "passwd [server name]",
		Short: "Change admin password",
		Long:  passwdDescription,
		Run:   runPasswd,
	}

	return cmd
}

func runPasswd(cmd *cobra.Command, args []string) {
	if len(args) < 1 || len(args) > 1 {
		cmd.Usage()
		os.Exit(-1)
	}

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.End()

	srv, err := sess.GetServer(args[0])
	if err != nil {
		log.Fatal(err)
	}

	oldPassword := readinput.ReadPassword("Old Password")
	newPassword := readinput.ReadPassword("New Password")
	confirmPassword := readinput.ReadPassword("Confirm Password")

	clt, err := client.New(srv.TSAURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := clt.AdminChangePassword(oldPassword, newPassword, confirmPassword); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Password changed successfully")
}

var passwdDescription = `
Change user password

`
