package api

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/juliengk/go-cert/ca"
	"github.com/juliengk/go-cert/ca/database"
	"github.com/juliengk/go-cert/errors"
	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/cli/command"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/kassisol/tsa/pkg/token"
	"github.com/labstack/echo"
	"golang.org/x/crypto/ocsp"
)

func RevokeCertHandle(c echo.Context) error {
	db, err := database.NewBackend("sqlite", command.CaDir)
	if err != nil {
		e := errors.New(errors.CertStoreError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}
	defer db.End()

	// Get POST data
	revokecert := new(RevokeCert)

	if err := c.Bind(revokecert); err != nil {
		r := jsonapi.NewErrorResponse(11000, "Internal error")

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	// Get JWT Claims
	jwt, _ := token.JWTFromHeader(c, "Authorization", "Bearer")
	claims, _ := token.GetCustomClaims(jwt)

	// Validate
	rcert := db.List(map[string]string{"serial": strconv.Itoa(revokecert.SerialNumber)})[0]

	if rcert.StatusFlag != "V" {
		e := errors.New(errors.OCSPError, errors.InvalidStatus)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	reCN := regexp.MustCompile(`CN=([a-z0-9\.\-\_]+)$`)

	cn := reCN.FindStringSubmatch(rcert.DistinguishedName)[1]

	if cn != claims.Audience && !claims.Admin {
		r := jsonapi.NewErrorResponse(11000, "Cannot revoke a certificate for which you are not the owner")

		return api.JSON(c, http.StatusBadRequest, r)
	}

	// Revoke certificate
	revocationDate := ca.DatabaseDateTimeFormat(time.Now())
	revocationReason := ocsp.CessationOfOperation

	db.Revoke(revokecert.SerialNumber, revocationDate, revocationReason)

	// Response
	response := jsonapi.NewSuccessResponseWithMessage(nil, 1001, "Certificate revoked")

	return api.JSON(c, http.StatusOK, response)
}
