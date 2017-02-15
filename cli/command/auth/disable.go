package auth

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/cli/command"
	"github.com/kassisol/tsa/storage"
	"github.com/spf13/cobra"
)

func newDisableCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable authentication",
		Long:  disableDescription,
		Run:   runDisable,
	}

	return cmd
}

func runDisable(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		cmd.Usage()
		os.Exit(-1)
	}

	s, err := storage.NewDriver("sqlite", command.DBFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer s.End()

	s.RemoveConfig("auth_type", "ALL")
	s.AddConfig("auth_type", "none")
}

var disableDescription = `
Disable authentication

`
