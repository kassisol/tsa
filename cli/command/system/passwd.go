package system

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-utils/password"
	"github.com/juliengk/go-utils/readinput"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/spf13/cobra"
)

func NewPasswdCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "passwd",
		Short: "Change admin password",
		Long:  passwdDescription,
		Run:   runPasswd,
	}

	return cmd
}

func runPasswd(cmd *cobra.Command, args []string) {
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
		log.Fatal(err)
	}
	defer s.End()

	oldPassword := readinput.ReadPassword("Old Password")
	newPassword := readinput.ReadPassword("New Password")
	confirmPassword := readinput.ReadPassword("Confirm Password")

	// Input validations
	// IV - Check old password
	if !password.ComparePassword([]byte(oldPassword), []byte(s.GetConfig("admin_password")[0].Value)) {
		log.Fatal("Old password invalid")
	}

	// IV - Check that new and confirm passwords matches
	if newPassword != confirmPassword {
		log.Fatal("New passwords does not match")
	}

	s.RemoveConfig("admin_password", "ALL")
	s.AddConfig("admin_password", password.GeneratePassword(newPassword))
}

var passwdDescription = `
Change user password

`
