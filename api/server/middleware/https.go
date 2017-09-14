package middleware

import (
	"errors"
	"net/http"

	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type HttpsConfig struct {
	Skipper middleware.Skipper
}

var (
	DefaultHttpsConfig = HttpsConfig{
		Skipper: middleware.DefaultSkipper,
	}

	ErrHttpsOnly = errors.New("Only https scheme is supported")
)

func Https() echo.MiddlewareFunc {
	c := DefaultHttpsConfig
	return HttpsWithConfig(c)
}

func HttpsWithConfig(config HttpsConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultHttpsConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			if c.Scheme() != "https" {
				r := jsonapi.NewErrorResponse(1000, ErrHttpsOnly.Error())

				return api.JSON(c, http.StatusInternalServerError, r)
			}

			return next(c)
		}
	}
}
