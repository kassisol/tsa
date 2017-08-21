package adf

import (
	"path"

	"github.com/juliengk/go-utils/filedir"
	"github.com/juliengk/go-utils/user"
)

type ServerConfig struct {
	AppDir string
}

func (c *ServerConfig) Init() error {
	u := user.New()

	c.AppDir = path.Join(u.HomeDir, ".tsa")

	if err := filedir.CreateDirIfNotExist(c.AppDir, false, 0750); err != nil {
		return err
	}

	return nil
}
