package middleware

import (
	"errors"
	"net/http"

	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/api/server/httputils"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/kassisol/tsa/pkg/token"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type AdminOnlyConfig struct {
	Skipper middleware.Skipper
}

var (
	DefaultAdminOnlyConfig = AdminOnlyConfig{
		Skipper: DefaultSkipper,
	}

	ErrAdminOnly = errors.New("Section restricted to admin only")
)

func AdminOnly() echo.MiddlewareFunc {
	c := DefaultAdminOnlyConfig
	return AdminOnlyWithConfig(c)
}

func AdminOnlyWithConfig(config AdminOnlyConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultAdminOnlyConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			// Get JWT Claims
			authHeader := c.Request().Header.Get("Authorization")
			jwt, _ := token.JWTFromHeader(authHeader, "Bearer")

			jwk, err := httputils.GetTokenSigningKey()
			if err != nil {
				return api.JSON(c, http.StatusInternalServerError, err)
			}

			t := token.New(jwk, true)
			claims, _ := t.GetCustomClaims(jwt)

			if !claims.Admin {
				r := jsonapi.NewErrorResponse(1000, ErrAdminOnly.Error())

				return api.JSON(c, http.StatusInternalServerError, r)
			}

			return next(c)
		}
	}
}
