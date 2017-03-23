package auth

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/cli/command"
	"github.com/kassisol/tsa/storage"
	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List authentication configurations",
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

	s, err := storage.NewDriver("sqlite", command.DBFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer s.End()

	configs := s.ListConfigs("auth")

	if len(configs) > 0 {
		w := tabwriter.NewWriter(os.Stdout, 20, 1, 2, ' ', 0)
		fmt.Fprintln(w, "KEY\tVALUE")

		for _, config := range configs {
			fmt.Fprintf(w, "%s\t%s\n", config.Key, config.Value)
		}

		w.Flush()
	}
}

var listDescription = `
List authentication configurations

`