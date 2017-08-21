package auth

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/spf13/cobra"
)

func newEnableCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable authentication",
		Long:  enableDescription,
		Run:   runEnable,
	}

	return cmd
}

func runEnable(cmd *cobra.Command, args []string) {
	if len(args) < 1 || len(args) > 1 {
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

	s.RemoveConfig("auth_type", "ALL")
	s.AddConfig("auth_type", args[0])
}

var enableDescription = `
Enable authentication

`
