package server

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-utils/validation"
	"github.com/juliengk/stack/client"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/tsa/cli/storage"
	"github.com/spf13/cobra"
)

func newAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [name] [tsa url]",
		Short: "Add TSA server",
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

	cfg := adf.NewServer()
	if err := cfg.Init(); err != nil {
		log.Fatal(err)
	}

	s, err := storage.NewDriver("sqlite", cfg.AppDir)
	if err != nil {
		log.Fatal(err)
	}
	defer s.End()

	// Input Validations
	// IV - Server name
	if err = validation.IsValidName(args[0]); err != nil {
		log.Fatal(err)
	}

	// IV - TSA URL
	if _, err := client.ParseUrl(args[1]); err != nil {
		log.Fatal(err)
	}

	s.AddServer(args[0], args[1], "")
}

var addDescription = `
Add TSA server

`
