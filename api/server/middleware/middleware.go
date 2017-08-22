package middleware

import (
	"github.com/juliengk/go-utils"
	"github.com/labstack/echo"
)

func DefaultSkipper(c echo.Context) bool {
	skipEndpoints := []string{
		"/",
		"/new-authz",
		"/system/admin/password",
	}

	if utils.StringInSlice(c.Path(), skipEndpoints, false) {
		return true
	}

	return false
}
