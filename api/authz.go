package api

import (
	"net/http"

	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/errors"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/kassisol/tsa/pkg/token"
	"github.com/labstack/echo"
)

func AuthzHandle(c echo.Context) error {
	username := c.Get("username").(string)
	admin := c.Get("admin").(bool)

	// Create the Claims
	jwk, err := token.GetSigningKey()
	if err != nil {
		e := errors.New(errors.DatabaseError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	ss, err := token.New(jwk, username, admin)
	if err != nil {
		e := errors.New(errors.DatabaseError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	// Response
	response := jsonapi.NewSuccessResponse(ss)

	return api.JSON(c, http.StatusOK, response)
}
