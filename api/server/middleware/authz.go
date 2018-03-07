package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/juliengk/go-utils"
	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/api/server/httputils"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/kassisol/tsa/pkg/tenant"
	"github.com/kassisol/tsa/pkg/token"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type AuthzConfig struct {
	Skipper middleware.Skipper
	Type    string
}

var (
	DefaultAuthzConfig = AuthzConfig{
		Skipper: DefaultSkipper,
		Type:    "anonymous",
	}

	allowedTypes = []string{
		"anonymous",
		"authenticated",
		"user",
		"admin",
	}

	ErrAuthz = errors.New("Section restricted")
)

func Authz(aType string) echo.MiddlewareFunc {
	c := AuthzConfig{
		Skipper: DefaultSkipper,
		Type:    aType,
	}
	return AuthzWithConfig(c)
}

func AuthzWithConfig(config AuthzConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultAuthzConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			if !utils.StringInSlice(config.Type, allowedTypes, false) {
				r := jsonapi.NewErrorResponse(1000, fmt.Sprintf("Type '%s' is not valid", config.Type))

				return api.JSON(c, http.StatusInternalServerError, r)
			}

			if config.Type == "anonymous" {
				return next(c)
			}

			cfg, err := tenant.New()
			if err != nil {
				r := jsonapi.NewErrorResponse(1000, err.Error())

				return api.JSON(c, http.StatusInternalServerError, r)
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

			cfg.SetUserGroups(claims.Groups)

			if config.Type == "authenticated" {
				if len(claims.Groups) > 0 {
					return next(c)
				}
			}

			if config.Type == "user" {
				id, _ := strconv.Atoi(c.Param("id"))
				if err := cfg.SetTenant(id); err != nil {
					return api.JSON(c, http.StatusInternalServerError, err)
				}

				if cfg.IsAdmin() {
					return next(c)
				}

				if cfg.IsMember() {
					return next(c)
				}
			}

			if config.Type == "admin" {
				if cfg.IsAdmin() {
					return next(c)
				}
			}

			r := jsonapi.NewErrorResponse(1000, ErrAuthz.Error())
			return api.JSON(c, http.StatusInternalServerError, r)
		}
	}
}
