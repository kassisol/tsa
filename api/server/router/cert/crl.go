package cert

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/juliengk/go-utils/filedir"
	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/api/errors"
	"github.com/kassisol/tsa/api/storage"
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

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		e := errors.New(errors.DatabaseError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}
	defer s.End()

	id, _ := strconv.Atoi(c.Param("id"))

	filters := map[string]string{
		"id": c.Param("id"),
	}
	tenants := s.ListTenants(filters)

	if len(tenants) == 0 {
		r := jsonapi.NewErrorResponse(1000, fmt.Sprintf("No Tenant with ID %d", id))

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	tenant := tenants[0]

	if err := cfg.Tenant(tenant.Name); err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	if !filedir.FileExists(cfg.CA.CRLFile) {
		r := jsonapi.NewErrorResponse(1000, "CRL file does not exist")

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	return c.File(cfg.CA.CRLFile)
}
