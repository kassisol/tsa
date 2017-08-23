package auth

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/cli/session"
	"github.com/kassisol/tsa/client"
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

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.End()

	srv, err := sess.Get()
	if err != nil {
		log.Fatal(err)
	}

	clt, err := client.New(srv.Server.TSAURL)
	if err != nil {
		log.Fatal(err)
	}

	configs, err := clt.AuthList(srv.Token)
	if err != nil {
		log.Fatal(err)
	}

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
