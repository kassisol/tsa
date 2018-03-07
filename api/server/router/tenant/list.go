package tenant

import (
	"net/http"

	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/api/errors"
	"github.com/kassisol/tsa/api/server/httputils"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/api/types"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/tsa/pkg/api"
	tntpkg "github.com/kassisol/tsa/pkg/tenant"
	"github.com/kassisol/tsa/pkg/token"
	"github.com/labstack/echo"
)

func ListHandle(c echo.Context) error {
	var result []types.Tenant

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

	itnt, err := tntpkg.New()
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

	itnt.SetUserGroups(claims.Groups)

	filters := make(map[string]string)

	tenants := s.ListTenants(filters)

	if itnt.IsAdmin() {
		result = tenants
	} else {
		for _, t := range tenants {
			if err := itnt.SetTenant(int(t.ID)); err != nil {
				return api.JSON(c, http.StatusInternalServerError, err)
			}

			if itnt.IsMember() {
				result = append(result, t)
			}
		}
	}

	response := jsonapi.NewSuccessResponse(result)

	return api.JSON(c, http.StatusOK, response)
}
