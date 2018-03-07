package cert

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/juliengk/go-cert/errors"
	"github.com/juliengk/stack/jsonapi"
	apierr "github.com/kassisol/tsa/api/errors"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/labstack/echo"
)

func CAPubHandle(c echo.Context) error {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		return err
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		e := apierr.New(apierr.DatabaseError, apierr.ReadFailed)
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
