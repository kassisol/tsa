package version

import (
	"os"
	"runtime"
	"strconv"
	"text/template"
	"time"
)

var (
	Version   = "unknown-version"
	GitCommit = "unknown-commit"
	GitState  = "unknown-state"
	BuildDate = "unknown-builddate"
)

var versionTemplate = `Client:
 Version:    {{.Client.Version}}
 Git commit: {{.Client.GitCommit}}{{if eq .Client.GitState "dirty"}}
 Git State:  {{.Client.GitState}}{{end}}
 Built:      {{.Client.BuildDate}}
 Go version: {{.Client.GoVersion}}
 OS/Arch:    {{.Client.OS}}/{{.Client.Arch}}{{$len := len .ServerError}}{{if eq $len 0}}

Server:
 Version:    {{.Server.Version}}
 Git commit: {{.Server.GitCommit}}{{if eq .Server.GitState "dirty"}}
 Git State:  {{.Server.GitState}}{{end}}
 Built:      {{.Server.BuildDate}}
 Go version: {{.Server.GoVersion}}
 OS/Arch:    {{.Server.OS}}/{{.Server.Arch}}{{end}}{{if gt $len 0}}

{{.ServerError}}{{end}}
`

type VersionInfo struct {
	Version   string
	GoVersion string
	GitCommit string
	GitState  string
	BuildDate string
	OS        string
	Arch      string
}

type VersionDisplay struct {
	ServerError string
	Client      *VersionInfo
	Server      *VersionInfo
}

func New() *VersionInfo {
	i, err := strconv.ParseInt(BuildDate, 10, 64)
	if err != nil {
		panic(err)
	}

	tu := time.Unix(i, 0)

	return &VersionInfo{
		Version:   Version,
		GoVersion: runtime.Version(),
		GitCommit: GitCommit,
		GitState:  GitState,
		BuildDate: tu.String(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

func NewDisplay(server *VersionInfo, err string) *VersionDisplay {
	display := new(VersionDisplay)

	// Get client version
	client := New()
	display.Client = client

	// Get server version
	display.ServerError = err
	display.Server = server

	return display
}

func (v *VersionDisplay) Show() {
	tmpl, err := template.New("version").Parse(versionTemplate)
	if err != nil {
		panic(err)
	}

	if err := tmpl.Execute(os.Stdout, v); err != nil {
		panic(err)
	}
}
