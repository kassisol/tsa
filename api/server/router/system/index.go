package system

import (
	"net/http"

	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/labstack/echo"
)

func IndexHandle(c echo.Context) error {
	response := jsonapi.NewSuccessResponse("TSA")

	return api.JSON(c, http.StatusOK, response)
}
