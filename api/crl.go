package api

import (
	"github.com/kassisol/tsa/cli/command"
	"github.com/labstack/echo"
)

func CRLHandle(c echo.Context) error {
	return c.File(command.CaCrlFile)
}
