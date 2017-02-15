package admin

import (
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-utils/filedir"
	"github.com/juliengk/go-utils/validation"
	"github.com/kassisol/tsa/cli/command"
	"github.com/kassisol/tsa/storage"
	"github.com/spf13/cobra"
)

var adminBindType string

func newBindCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bind [value]",
		Short: "Update bind address and port",
		Long:  bindDescription,
		Run:   runBind,
	}

	flags := cmd.Flags()
	flags.StringVarP(&adminBindType, "type", "t", "address", "Bind type (address or port)")

	return cmd
}

func runBind(cmd *cobra.Command, args []string) {
	if len(args) < 1 && len(args) > 1 {
		cmd.Usage()
		os.Exit(-1)
	}

	if !filedir.FileExists(command.DBFilePath) {
		log.Fatal("Initialization needs to be done first")
	}

	// DB
	s, err := storage.NewDriver("sqlite", command.DBFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer s.End()

	// Input validations
	if adminBindType != "address" && adminBindType != "port" {
		log.Fatal("Bind type should be either \"address\" or \"port\"")
	}

	if adminBindType == "address" {
		// IV - API Bind
		if err = validation.IsValidIP(args[0]); err != nil {
			log.Fatal(err)
		}
	}

	if adminBindType == "port" {
		// IV - API Port
		port, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}
		if err = validation.IsValidPort(port); err != nil {
			log.Fatal(err)
		}
	}

	key := "api_bind"
	if adminBindType == "port" {
		key = "api_port"
	}

	// Remove old password
	s.RemoveConfig(key, "ALL")

	// Save new password
	s.AddConfig(key, args[0])
}

var bindDescription = `
Update bind address and port

`
