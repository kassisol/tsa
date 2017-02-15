package api

import (
	"net/http"

	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/labstack/echo"
)

func IndexHandle(c echo.Context) error {
	directory := Directory{
		CAInfo:     "/info",
		NewAuthz:   "/new-authz",
		NewApp:     "/acme/new-app",
		RevokeCert: "/acme/revoke-cert",
	}

	response := jsonapi.NewSuccessResponse(directory)

	return api.JSON(c, http.StatusOK, response)
}
