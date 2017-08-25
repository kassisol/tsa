package crl

import (
	"net/http"

	"github.com/juliengk/go-utils/filedir"
	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/labstack/echo"
)

func CRLHandle(c echo.Context) error {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	if !filedir.FileExists(cfg.CA.CRLFile) {
		r := jsonapi.NewErrorResponse(1000, "CRL file does not exist")

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	return c.File(cfg.CA.CRLFile)
}
