package adf

import (
	"path"

	"github.com/juliengk/go-utils/filedir"
)

type EngineConfig struct {
	CertsDir string
	TLS      TLSOptions
}

func (c *EngineConfig) Init() error {
	c.CertsDir = "/etc/docker/tls"

	if err := filedir.CreateDirIfNotExist(c.CertsDir, false, 0750); err != nil {
		return err
	}

	c.TLS.CaFile = path.Join(c.CertsDir, "ca.pem")
	c.TLS.KeyFile = path.Join(c.CertsDir, "server-key.pem")
	c.TLS.CrtFile = path.Join(c.CertsDir, "server-cert.pem")

	return nil
}
