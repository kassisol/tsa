package middleware

import (
	"errors"
	"net/http"

	aerrors "github.com/juliengk/go-cert/errors"
	"github.com/juliengk/go-utils/password"
	"github.com/juliengk/stack/jsonapi"
	apierr "github.com/kassisol/tsa/api/errors"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type AdminPasswordConfig struct {
	Skipper middleware.Skipper
}

var (
	DefaultAdminPasswordConfig = AdminPasswordConfig{
		Skipper: DefaultSkipper,
	}

	ErrDefaultAdminPasswordSet = errors.New("Default admin password should be changed")
)

func AdminPassword() echo.MiddlewareFunc {
	c := DefaultAdminPasswordConfig
	return AdminPasswordWithConfig(c)
}

func AdminPasswordWithConfig(config AdminPasswordConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultAdminPasswordConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			cfg := adf.NewDaemon()
			if err := cfg.Init(); err != nil {
				return err
			}

			s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
			if err != nil {
				e := aerrors.New(apierr.DatabaseError, aerrors.ReadFailed)
				r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

				return api.JSON(c, http.StatusInternalServerError, r)
			}
			defer s.End()

			if password.ComparePassword([]byte("admin"), []byte(s.GetConfig("admin_password")[0].Value)) {
				r := jsonapi.NewErrorResponse(1000, ErrDefaultAdminPasswordSet.Error())

				return api.JSON(c, http.StatusInternalServerError, r)
			}

			return next(c)
		}
	}
}
