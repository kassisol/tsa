package access

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/cli/session"
	"github.com/spf13/cobra"
)

func newStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Session status",
		Long:  statusDescription,
		Run:   runStatus,
	}

	return cmd
}

func runStatus(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		cmd.Usage()
		os.Exit(-1)
	}

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.End()

	fmt.Println(sess.Status())
}

var statusDescription = `
Session status

`
