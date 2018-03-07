package system

import (
	"net/http"

	"github.com/juliengk/go-cert/errors"
	"github.com/juliengk/stack/jsonapi"
	apierr "github.com/kassisol/tsa/api/errors"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/api/types"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/kassisol/tsa/pkg/host"
	"github.com/kassisol/tsa/version"
	"github.com/labstack/echo"
)

func InfoHandle(c echo.Context) error {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		e := errors.New(apierr.DatabaseError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}
	defer s.End()

	info := types.SystemInfo{
		ServerVersion: version.Version,
		ID:            s.GetConfig("server_id")[0].Value,
		StorageDriver: "sqlite",
		LoggingDriver: "standard",
		TSARootDir:    cfg.App.Dir.Root,
	}

	apiHost := host.New(c.Request().URL, c.Request().Host)
	if s.CountConfigKey("api_fqdn") == 1 {
		apiHost = s.GetConfig("api_fqdn")[0].Value
	}

	api2 := types.API{
		FQDN:        apiHost,
		BindAddress: s.GetConfig("api_bind")[0].Value,
		BindPort:    s.GetConfig("api_port")[0].Value,
	}
	info.API = api2

	auth := types.Auth{
		Type: s.GetConfig("auth_type")[0].Value,
	}
	info.Auth = auth

	response := jsonapi.NewSuccessResponse(info)

	return api.JSON(c, http.StatusOK, response)
}
