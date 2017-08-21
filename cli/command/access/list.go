package access

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/cli/session"
	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List sessions",
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

	sessions := sess.List()

	if len(sessions) > 0 {
		w := tabwriter.NewWriter(os.Stdout, 20, 1, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tACTIVE\tSERVER\tEXPIRE")

		for _, s := range sessions {
			expire := sess.GetExpire(s.Token)

			fmt.Fprintf(w, "%d\t%t\t%s\t%s\n", s.ID, s.Active, s.Server, expire.String())
		}

		w.Flush()
	}
}

var listDescription = `
List sessions

`
