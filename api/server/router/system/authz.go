package system

import (
	"net/http"
	"strconv"

	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/api/errors"
	"github.com/kassisol/tsa/api/server/httputils"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/kassisol/tsa/pkg/token"
	"github.com/labstack/echo"
)

func AuthzHandle(c echo.Context) error {
	username := c.Get("username").(string)
	admin := c.Get("admin").(bool)

	qTtl := c.QueryParam("ttl")
	ttl, _ := strconv.Atoi(qTtl)

	// Create the Claims
	jwk, err := httputils.GetTokenSigningKey()
	if err != nil {
		e := errors.New(errors.DatabaseError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	t := token.New(jwk, true)

	ss, err := t.Create(username, "harbormaster", admin, ttl)
	if err != nil {
		e := errors.New(errors.DatabaseError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	// Response
	response := jsonapi.NewSuccessResponse(ss)

	return api.JSON(c, http.StatusOK, response)
}
