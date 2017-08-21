package crl

import (
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/labstack/echo"
)

func CRLHandle(c echo.Context) error {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		return err
	}

	return c.File(cfg.CA.CRLFile)
}
