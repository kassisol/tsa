package access

import (
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/cli/session"
	"github.com/spf13/cobra"
)

func newUnuseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unuse [session_id]",
		Short: "Unuse session",
		Long:  unuseDescription,
		Run:   runUnuse,
	}

	return cmd
}

func runUnuse(cmd *cobra.Command, args []string) {
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

	if err := sess.Unuse(uint(id)); err != nil {
		log.Fatal(err)
	}
}

var unuseDescription = `
Unuse session

`
