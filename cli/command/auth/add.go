package auth

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/api/auth"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/pkg/adf"
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

	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		log.Fatal(err)
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
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
