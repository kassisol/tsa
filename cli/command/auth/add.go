package auth

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/auth"
	"github.com/kassisol/tsa/cli/command"
	"github.com/kassisol/tsa/storage"
	"github.com/spf13/cobra"
)

func newAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [key] [value]",
		Short: "Add auth configuration",
		Long:  addDescription,
		Run:   runAdd,
	}

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) {
	if len(args) < 2 || len(args) > 2 {
		cmd.Usage()
		os.Exit(-1)
	}

	s, err := storage.NewDriver("sqlite", command.DBFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer s.End()

	authType := s.GetConfig("auth_type")[0].Value

	if authType == "none" {
		log.Fatal("Make sure to enable a backend auth before adding configuration")
	}

	a, err := auth.NewDriver(authType)
	if err != nil {
		log.Fatal(err)
	}

	if err = a.AddConfig(args[0], args[1]); err != nil {
		log.Fatal(err)
	}
}

var addDescription = `
Add auth configuration

`
