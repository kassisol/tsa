package server

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/cli/storage"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List TSA servers",
		Long:    listDescription,
		Run:     runList,
	}

	return cmd
}

func runList(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
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

	servers := s.ListServers(map[string]string{})

	if len(servers) > 0 {
		w := tabwriter.NewWriter(os.Stdout, 20, 1, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tTSA URL")

		for _, s := range servers {
			fmt.Fprintf(w, "%s\t%s\n", s.Name, s.TSAURL)
		}

		w.Flush()
	}
}

var listDescription = `
List TSA servers

`
