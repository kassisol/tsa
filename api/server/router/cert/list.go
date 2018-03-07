package cert

import (
	"net/http"
	"strconv"

	"github.com/juliengk/go-cert/ca/database"
	"github.com/juliengk/go-cert/ca/database/backend"
	"github.com/juliengk/go-utils/validation"
	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/api/errors"
	"github.com/kassisol/tsa/api/server/httputils"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/kassisol/tsa/pkg/tenant"
	"github.com/kassisol/tsa/pkg/token"
	"github.com/labstack/echo"
)

func ListHandle(c echo.Context) error {
	var certs []backend.CertificateResult

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

	id, _ := strconv.Atoi(c.Param("id"))

	tnt, err := tenant.New()
	if err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	if err := tnt.SetTenant(id); err != nil {
		return api.JSON(c, http.StatusInternalServerError, err)
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

	tnt.SetUserGroups(claims.Groups)

	tenant2 := tnt.GetTenant()
	if err := cfg.Tenant(tenant2.Name); err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	db, err := database.NewBackend("sqlite", cfg.CA.Dir.Root)
	if err != nil {
		e := errors.New(errors.DatabaseError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}
	defer db.End()

	qp := c.QueryParams()
	dbFilters := httputils.QueryParams2Filters(qp)

	certificates := db.List(dbFilters)

	if tnt.IsAdmin() {
		certs = certificates
	} else if tnt.IsMember() {
		for _, certificate := range certificates {
			if err := validation.IsValidFQDN(certificate.DistinguishedName); err == nil {
				if certificate.StatusFlag == "V" {
					certs = append(certs, certificate)
				}
			}
		}
	}

	response := jsonapi.NewSuccessResponse(certificates)

	return api.JSON(c, http.StatusOK, response)
}
