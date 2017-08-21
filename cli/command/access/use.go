package access

import (
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/cli/session"
	"github.com/spf13/cobra"
)

func newUseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use [session_id]",
		Short: "Use session",
		Long:  useDescription,
		Run:   runUse,
	}

	return cmd
}

func runUse(cmd *cobra.Command, args []string) {
	if len(args) < 1 || len(args) > 1 {
		cmd.Usage()
		os.Exit(-1)
	}

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.End()

	id, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(err)
	}

	if err := sess.Use(uint(id)); err != nil {
		log.Fatal(err)
	}
}

var useDescription = `
Use session

`
