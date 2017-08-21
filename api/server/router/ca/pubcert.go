package ca

import (
	"io/ioutil"
	"net/http"

	"github.com/juliengk/go-cert/errors"
	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/labstack/echo"
)

func PubCertHandle(c echo.Context) error {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		return err
	}

	// Read CA certificate file
	cert, err := ioutil.ReadFile(cfg.CA.TLS.CrtFile)
	if err != nil {
		e := errors.New(errors.RootError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	// Response
	response := jsonapi.NewSuccessResponse(string(cert))

	return api.JSON(c, http.StatusOK, response)
}
