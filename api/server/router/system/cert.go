package system

import (
	"net/http"
	"strconv"
	"time"

	"github.com/juliengk/go-cert/ca"
	"github.com/juliengk/go-cert/ca/database"
	"golang.org/x/crypto/ocsp"

	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/api/errors"
	"github.com/kassisol/tsa/api/server/httputils"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/labstack/echo"
)

func CertListHandle(c echo.Context) error {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	db, err := database.NewBackend("sqlite", cfg.CA.Dir.Root)
	if err != nil {
		e := errors.New(errors.DatabaseError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}
	defer db.End()

	qp := c.QueryParams()
	filters := httputils.QueryParams2Filters(qp)

	certificates := db.List(filters)

	response := jsonapi.NewSuccessResponse(certificates)

	return api.JSON(c, http.StatusOK, response)
}

func CertRevokeHandle(c echo.Context) error {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	db, err := database.NewBackend("sqlite", cfg.CA.Dir.Root)
	if err != nil {
		e := errors.New(errors.DatabaseError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}
	defer db.End()

	serialNumber, err := strconv.Atoi(c.Param("serialnumber"))
	if err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	revocationDate := ca.DatabaseDateTimeFormat(time.Now())
	revocationReason := ocsp.CessationOfOperation

	db.Revoke(serialNumber, revocationDate, revocationReason)

	return c.NoContent(http.StatusNoContent)
}
