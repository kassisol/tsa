package middleware

import (
	"github.com/juliengk/go-utils"
	"github.com/labstack/echo"
)

func generateEndpointsSlice(key string, restricted bool) []string {
	var result []string

	skipKeys := []string{
		"all",
		"all_restricted",
	}
	skipEndpoints := map[string][]string{
		"all": []string{
			"/",
			"/new-authz",
			"/system/admin/password",
		},
		"all_restricted": []string{
			"/system/info",
		},
		//"admin_password": []string{
		//},
		"ca_init": []string{
			"/system/ca/init",
		},
	}

	result = append(result, skipEndpoints["all"]...)

	if restricted {
		result = append(result, skipEndpoints["all_restricted"]...)
	}

	if !utils.StringInSlice(key, skipKeys, false) {
		if v, ok := skipEndpoints[key]; ok {
			result = append(result, v...)
		}
	}

	return result
}

func DefaultSkipper(c echo.Context) bool {
	skipEndpoints := generateEndpointsSlice("all", false)

	if utils.StringInSlice(c.Path(), skipEndpoints, false) {
		return true
	}

	return false
}

func DefaultAdminPasswordSkipper(c echo.Context) bool {
	skipEndpoints := generateEndpointsSlice("admin_password", true)

	if utils.StringInSlice(c.Path(), skipEndpoints, false) {
		return true
	}

	return false
}

func DefaultCAInitSkipper(c echo.Context) bool {
	skipEndpoints := generateEndpointsSlice("ca_init", true)

	if utils.StringInSlice(c.Path(), skipEndpoints, false) {
		return true
	}

	return false
}
