package server

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/cli/storage"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/spf13/cobra"
)

func newRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm [name]",
		Aliases: []string{"remove"},
		Short:   "Remove TSA server",
		Long:    removeDescription,
		Run:     runRemove,
	}

	return cmd
}

func runRemove(cmd *cobra.Command, args []string) {
	if len(args) < 1 || len(args) > 1 {
		cmd.Usage()
		os.Exit(-1)
	}

	cfg := adf.NewServer()
	if err := cfg.Init(); err != nil {
		log.Fatal(err)
	}

	s, err := storage.NewDriver("sqlite", cfg.AppDir)
	if err != nil {
		log.Fatal(err)
	}
	defer s.End()

	s.RemoveServer(args[0])
}

var removeDescription = `
Remove TSA server

`
