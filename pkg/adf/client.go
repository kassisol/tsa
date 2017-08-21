package adf

import (
	"path"

	"github.com/juliengk/go-utils/filedir"
	"github.com/juliengk/go-utils/user"
)

type ClientConfig struct {
	App     App
	Profile Profile
	TLS     TLSOptions
}

type Profile struct {
	Name    string
	CertDir string
}

func (c *ClientConfig) Init() error {
	u := user.New()

	appDir := path.Join(u.HomeDir, ".twic")
	certsDir := path.Join(appDir, "certs")

	if err := filedir.CreateDirIfNotExist(appDir, false, 0750); err != nil {
		return err
	}

	if err := filedir.CreateDirIfNotExist(certsDir, false, 0750); err != nil {
		return err
	}

	c.App.Dir.Root = appDir
	c.App.Dir.Certs = certsDir

	return nil
}

func (c *ClientConfig) SetName(name string) error {
	c.Profile.Name = name

	certNameDir := path.Join(c.App.Dir.Certs, c.Profile.Name)

	if err := filedir.CreateDirIfNotExist(certNameDir, false, 0750); err != nil {
		return err
	}

	c.Profile.CertDir = certNameDir

	c.TLS.CaFile = path.Join(certNameDir, "ca.pem")
	c.TLS.KeyFile = path.Join(certNameDir, "key.pem")
	c.TLS.CrtFile = path.Join(certNameDir, "cert.pem")

	return nil
}
