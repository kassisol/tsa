// HBM TSA is an application acting as a CA (Certification Authority) server
// to issue certificates for Docker Engine with TLS enabled.
// Copyright (C) 2017 Kassisol inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-utils"
	"github.com/juliengk/go-utils/user"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/spf13/cobra"
)

var (
	serverBindAddress string
	serverBindPort    int
	serverTLS         bool
	serverTLSCN       string
	serverTLSDuration int
	serverTLSGen      bool
	serverTLSCert     string
	serverTLSKey      string
)

func newCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tsad",
		Short: "Starts a CA (Certification Authority) server",
		Long:  daemonDescription,
		Run:   runDaemon,
	}

	cmd.SetHelpTemplate(helpTemplate)
	cmd.SetUsageTemplate(usageTemplate)

	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		log.Fatal(err)
	}

	flags := cmd.Flags()
	flags.StringVarP(&serverBindAddress, "bind-address", "a", "0.0.0.0", "Bind Address")
	flags.IntVarP(&serverBindPort, "bind-port", "p", 80, "Bind Port")
	flags.BoolVarP(&serverTLS, "tls", "t", false, "Enable TLS certificates")
	flags.StringVarP(&serverTLSCN, "tlscn", "", "", "Certificate Common Name")
	flags.IntVarP(&serverTLSDuration, "tls-duration", "", 60, "Certificate duration")
	flags.BoolVarP(&serverTLSGen, "tlsgen", "", false, "Auto generate self-signed TLS certificates")
	flags.StringVarP(&serverTLSCert, "tlscert", "", cfg.API.CrtFile, "Path to TLS certificate file")
	flags.StringVarP(&serverTLSKey, "tlskey", "", cfg.API.KeyFile, "Path to TLS key file")

	return cmd
}

func main() {
	u := user.New()

	if !u.IsRoot() {
		log.Fatal("You must be root to run that command")
	}

	cmd := newCommand()
	if err := cmd.Execute(); err != nil {
		utils.Exit(err)
	}
}

var daemonDescription = `
Starts a CA (Certificate Authority) server

`

var usageTemplate = `{{ .Short | trim }}

Usage:{{ if .Runnable }}
  {{ if .HasAvailableFlags }}{{ appendIfNotPresent .UseLine "[flags]" }}{{ else }}{{ .UseLine }}{{ end }}{{ end }}{{ if .HasAvailableSubCommands }}
  {{ .CommandPath }} [command]{{ end }}{{ if gt .Aliases 0 }}

Aliases:
  {{ .NameAndAliases }}{{ end }}{{ if .HasExample }}

Examples:
  {{ .Example }}{{ end }}{{ if .HasAvailableSubCommands }}

Available Commands:{{ range .Commands }}{{ if .IsAvailableCommand }}
  {{ rpad .Name .NamePadding }} {{ .Short }}{{ end }}{{ end }}{{ end }}{{ if .HasAvailableLocalFlags }}

Flags:
{{ .LocalFlags.FlagUsages | trimRightSpace }}{{ end }}{{ if .HasAvailableInheritedFlags }}

Global Flags:
  {{ .InheritedFlags.FlagUsages | trimRightSpace }}{{ end }}{{ if .HasHelpSubCommands }}

Additional help topics:{{ range .Commands }}{{ if .IsHelpCommand }}
  {{ rpad .CommandPath .CommandPathPadding }} {{ .Short }}{{ end }}{{ end }}{{ end }}{{ if .HasAvailableSubCommands }}

Use "{{ .CommandPath }} [command] --help" for more information about a command.{{ end }}
`

var helpTemplate = `
{{ if or .Runnable .HasSubCommands }}{{ .UsageString }}{{ end }}`
