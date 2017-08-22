package system

import (
	"net/http"

	"github.com/juliengk/go-utils/password"
	"github.com/juliengk/stack/jsonapi"
	apierr "github.com/kassisol/tsa/api/errors"
	"github.com/kassisol/tsa/api/types"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/labstack/echo"
)

func AdminChangePasswordHandle(c echo.Context) error {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		e := apierr.New(apierr.DatabaseError, apierr.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}
	defer s.End()

	// Get POST data
	p := new(types.ChangePassword)

	if err := c.Bind(p); err != nil {
		r := jsonapi.NewErrorResponse(1000, "Cannot unmarshal JSON")

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	// Validation
	if len(p.New) <= 0 {
		r := jsonapi.NewErrorResponse(1000, "Empty password is not allowed")

		return api.JSON(c, http.StatusUnprocessableEntity, r)
	}

	if !password.ValidatePassword(p.New) {
		r := jsonapi.NewErrorResponse(1000, "Password is not valid")

		return api.JSON(c, http.StatusUnprocessableEntity, r)
	}

	if !password.ComparePassword([]byte(p.Old), []byte(s.GetConfig("admin_password")[0].Value)) {
		r := jsonapi.NewErrorResponse(1000, "Old password invalid")

		return api.JSON(c, http.StatusUnprocessableEntity, r)
	}

	if p.New != p.Confirm {
		r := jsonapi.NewErrorResponse(1000, "New passwords do not match")

		return api.JSON(c, http.StatusUnprocessableEntity, r)
	}


	s.RemoveConfig("admin_password", "ALL")
	s.AddConfig("admin_password", password.GeneratePassword(p.New))

	// Response
	response := jsonapi.NewSuccessResponse("Password changed")

	return api.JSON(c, http.StatusOK, response)
}
