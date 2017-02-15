package api

import (
	"encoding/json"

	"github.com/juliengk/stack/jsonapi"
	"github.com/labstack/echo"
)

func JSON(c echo.Context, code int, i interface{}) error {
	mime := jsonapi.BuildVendorMIME("harbormaster")

	c.Response().Header().Set(echo.HeaderContentType, mime)
	c.Response().WriteHeader(code)

	return json.NewEncoder(c.Response()).Encode(i)
}
