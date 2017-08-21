package auth

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/spf13/cobra"
)

func newRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm [key] [value]",
		Aliases: []string{"remove"},
		Short:   "Remove authentication configuration",
		Long:    removeDescription,
		Run:     runRemove,
	}

	return cmd
}

func runRemove(cmd *cobra.Command, args []string) {
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

	key := args[0]
	value := args[1]

	if len(args[1]) == 0 {
		value = "ALL"
	}

	s.RemoveConfig(key, value)
}

var removeDescription = `
Remove authentication configuration

`
