package api

import (
	"github.com/kassisol/tsa/api/config"
	"github.com/labstack/echo"
)

func CRLHandle(c echo.Context) error {
	return c.File(config.CaCrlFile)
}
