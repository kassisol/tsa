package adf

import (
	"path"

	"github.com/juliengk/go-utils/filedir"
)

type DaemonConfig struct {
	App App
	CA  CA
	API TLSOptions
}

type CA struct {
	Dir     CADir
	TLS     TLSOptions
	CRLFile string
}

type CADir struct {
	Root    string
	Private string
	Certs   string
}

func (c *DaemonConfig) Init() error {
	c.App.Dir.Root = "/var/lib/tsa"

	c.CA.Dir.Root = path.Join(c.App.Dir.Root, "ca")
	c.CA.Dir.Private = path.Join(c.CA.Dir.Root, "private")
	c.CA.Dir.Certs = path.Join(c.CA.Dir.Root, "certs")

	c.CA.TLS.CrtFile = path.Join(c.CA.Dir.Certs, "ca.crt")
	c.CA.CRLFile = path.Join(c.CA.Dir.Root, "CRL.crl")

	c.App.Dir.Certs = path.Join(c.App.Dir.Root, "certs")

	c.API.KeyFile = path.Join(c.App.Dir.Certs, "api.key")
	c.API.CrtFile = path.Join(c.App.Dir.Certs, "api.crt")

	if err := filedir.CreateDirIfNotExist(c.App.Dir.Root, false, 0700); err != nil {
		return err
	}

	if err := filedir.CreateDirIfNotExist(c.App.Dir.Certs, false, 0750); err != nil {
		return err
	}

	return nil
}
