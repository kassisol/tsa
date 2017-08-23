package system

import (
	"net/http"

	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/api/auth"
	"github.com/kassisol/tsa/api/errors"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/api/types"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/labstack/echo"
)

func AuthListHandle(c echo.Context) error {
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

	configs := s.ListConfigs("auth")

	response := jsonapi.NewSuccessResponse(configs)

	return api.JSON(c, http.StatusOK, response)
}

func AuthCreateHandle(c echo.Context) error {
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

	// Get POST data
	config := new(types.ServerConfig)

	if err := c.Bind(config); err != nil {
		r := jsonapi.NewErrorResponse(1000, "Cannot unmarshal JSON")

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	authType := s.GetConfig("auth_type")[0].Value

	if authType == "none" {
		r := jsonapi.NewErrorResponse(1000, "Make sure to enable a backend auth before adding configuration")

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	a, err := auth.NewDriver(authType)
	if err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	if err = a.AddConfig(config.Key, config.Value); err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	response := jsonapi.NewSuccessResponse(config)

	return api.JSON(c, http.StatusCreated, response)
}

func AuthDeleteHandle(c echo.Context) error {
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

	key := c.Param("key")
	value := c.QueryParam("value")

	s.RemoveConfig(key, value)

	return c.NoContent(http.StatusNoContent)
}

func AuthEnableHandle(c echo.Context) error {
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

	atype := c.Param("type")

	s.RemoveConfig("auth_type", "ALL")
	s.AddConfig("auth_type", atype)

	return c.NoContent(http.StatusNoContent)
}

func AuthDisableHandle(c echo.Context) error {
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

	s.RemoveConfig("auth_type", "ALL")
	s.AddConfig("auth_type", "none")

	return c.NoContent(http.StatusNoContent)
}
