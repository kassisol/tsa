package tenant

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/api/errors"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/api/types"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/labstack/echo"
)

func GroupListHandle(c echo.Context) error {
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

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("%s", r)
			res := jsonapi.NewErrorResponse(1000, err)

			api.JSON(c, http.StatusUnprocessableEntity, res)
		}
	}()

	// Get POST data
	tg := new(types.TenantGroup)

	if err := c.Bind(tg); err != nil {
		r := jsonapi.NewErrorResponse(1000, "Cannot unmarshal JSON")

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	if err := s.AddGroupToTenant(tg.Tenant, tg.Group); err != nil {
		panic(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func GroupAddHandle(c echo.Context) error {
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

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("%s", r)
			res := jsonapi.NewErrorResponse(1000, err)

			api.JSON(c, http.StatusUnprocessableEntity, res)
		}
	}()

	// Get POST data
	tg := new(types.TenantGroup)

	if err := c.Bind(tg); err != nil {
		r := jsonapi.NewErrorResponse(1000, "Cannot unmarshal JSON")

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	if err := s.AddGroupToTenant(tg.Tenant, tg.Group); err != nil {
		panic(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func GroupDeleteHandle(c echo.Context) error {
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

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("%s", r)
			res := jsonapi.NewErrorResponse(1000, err)

			api.JSON(c, http.StatusUnprocessableEntity, res)
		}
	}()

	iid, _ := strconv.Atoi(c.Param("id"))
	gid, _ := strconv.Atoi(c.Param("gid"))

	if err := s.RemoveGroupFromTenant(iid, gid); err != nil {
		panic(err)
	}

	return c.NoContent(http.StatusNoContent)
}
