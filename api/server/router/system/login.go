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

func LoginHandle(c echo.Context) error {
	var ttlTemp string

	username := c.Get("username").(string)
	groups := c.Get("groups").([]string)

	qTtl := c.QueryParam("ttl")

	if len(qTtl) == 0 {
		ttlTemp = "1440"
	} else {
		ttlTemp = qTtl
	}

	ttl, err := strconv.Atoi(ttlTemp)
	if err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())

		return api.JSON(c, http.StatusUnprocessableEntity, r)
	}

	// Create the Claims
	jwk, err := httputils.GetTokenSigningKey()
	if err != nil {
		e := errors.New(errors.DatabaseError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	t := token.New(jwk, true)

	ss, err := t.Create(username, "harbormaster", groups, ttl)
	if err != nil {
		e := errors.New(errors.DatabaseError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	// Response
	response := jsonapi.NewSuccessResponse(ss)

	return api.JSON(c, http.StatusOK, response)
}
