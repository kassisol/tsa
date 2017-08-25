package system

import (
	"net/http"

	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/kassisol/tsa/version"
	"github.com/labstack/echo"
)

func ServerVersionHandle(c echo.Context) error {
	version := version.New()

	response := jsonapi.NewSuccessResponse(version)

	return api.JSON(c, http.StatusOK, response)
}
