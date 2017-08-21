package access

import (
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/cli/session"
	"github.com/spf13/cobra"
)

var sessionRemoveAll bool

func newRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm [session_id]",
		Aliases: []string{"remove"},
		Short:   "Remove session",
		Long:    removeDescription,
		Run:     runRemove,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&sessionRemoveAll, "all", "a", false, "Remove all sessions")

	return cmd
}

func runRemove(cmd *cobra.Command, args []string) {
	if !sessionRemoveAll && (len(args) < 1 || len(args) > 1) {
		cmd.Usage()
		os.Exit(-1)
	}

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.End()

	if sessionRemoveAll {
		if err := sess.Clear(); err != nil {
			log.Fatal(err)
		}
	} else {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}

		if err := sess.Remove(uint(id)); err != nil {
			log.Fatal(err)
		}
	}
}

var removeDescription = `
Remove session

`
