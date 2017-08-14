package system

import (
	"net/http"

	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/api/types"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/labstack/echo"
)

func IndexHandle(c echo.Context) error {
	directory := types.Directory{
		CAInfo:     "/ca",
		NewAuthz:   "/new-authz",
		NewApp:     "/acme/new-app",
		RevokeCert: "/acme/revoke-cert",
	}

	response := jsonapi.NewSuccessResponse(directory)

	return api.JSON(c, http.StatusOK, response)
}
