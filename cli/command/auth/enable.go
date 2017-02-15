package auth

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/cli/command"
	"github.com/kassisol/tsa/storage"
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

	s, err := storage.NewDriver("sqlite", command.DBFilePath)
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
