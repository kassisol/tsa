package middleware

import (
	"errors"
	"net/http"

	aerrors "github.com/juliengk/go-cert/errors"
	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/cli/command"
	apierr "github.com/kassisol/tsa/errors"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/kassisol/tsa/storage"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type CAInitConfig struct {
	Skipper middleware.Skipper
}

var (
	DefaultCAInitConfig = CAInitConfig{
		Skipper: middleware.DefaultSkipper,
	}

	ErrCAInit = errors.New("CA initialization should be done first")
)

func CAInit() echo.MiddlewareFunc {
	c := DefaultCAInitConfig
	return CAInitWithConfig(c)
}

func CAInitWithConfig(config CAInitConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultCAInitConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			s, err := storage.NewDriver("sqlite", command.DBFilePath)
			if err != nil {
				e := aerrors.New(apierr.DatabaseError, aerrors.ReadFailed)
				r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

				return api.JSON(c, http.StatusInternalServerError, r)
			}
			defer s.End()

			if len(s.ListConfigs("ca")) == 0 {
				r := jsonapi.NewErrorResponse(1000, ErrCAInit.Error())

				return api.JSON(c, http.StatusInternalServerError, r)
			}

			return next(c)
		}
	}
}
