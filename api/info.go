package api

import (
	"io/ioutil"
	"net/http"

	"github.com/juliengk/go-cert/errors"
	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/cli/command"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/labstack/echo"
)

func InfoHandle(c echo.Context) error {
	// Read CA certificate file
	cert, err := ioutil.ReadFile(command.CaCrtFile)
	if err != nil {
		e := errors.New(errors.RootError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	// Response
	response := jsonapi.NewSuccessResponse(string(cert))

	return api.JSON(c, http.StatusOK, response)
}
