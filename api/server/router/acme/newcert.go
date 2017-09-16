package acme

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-cert/ca"
	"github.com/juliengk/go-cert/ca/database"
	"github.com/juliengk/go-cert/errors"
	"github.com/juliengk/go-cert/pkix"
	"github.com/juliengk/go-utils/validation"
	"github.com/juliengk/stack/jsonapi"
	apierr "github.com/kassisol/tsa/api/errors"
	"github.com/kassisol/tsa/api/server/httputils"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/api/types"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/kassisol/tsa/pkg/host"
	"github.com/kassisol/tsa/pkg/token"
	"github.com/labstack/echo"
)

func NewCertHandle(c echo.Context) error {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		return err
	}

	db, err := database.NewBackend("sqlite", cfg.CA.Dir.Root)
	if err != nil {
		e := errors.New(errors.CertStoreError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}
	defer db.End()

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		e := errors.New(apierr.DatabaseError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}
	defer s.End()

	// Get JWT Claims
	authHeader := c.Request().Header.Get("Authorization")
	jwt, _ := token.JWTFromHeader(authHeader, "Bearer")

	jwk, err := httputils.GetTokenSigningKey()
	if err != nil {
		return api.JSON(c, http.StatusInternalServerError, err)
	}

	t := token.New(jwk, true)
	claims, _ := t.GetCustomClaims(jwt)

	// Get POST data
	newcert := new(types.NewCert)

	if err := c.Bind(newcert); err != nil {
		log.Info(err)
		e := errors.New(errors.CSRError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusUnprocessableEntity, r)
	}

	// Read CSR
	csr, err := pkix.NewCertificateRequestFromDER(newcert.CSR)
	if err != nil {
		e := errors.New(errors.CSRError, errors.ParseFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusUnprocessableEntity, r)
	}

	// Validate Common Name for client
	if newcert.Type == "client" && claims.Audience != csr.CR.Subject.CommonName {
		e := errors.New(errors.CertStoreError, errors.RecordFound)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusBadRequest, r)
	}

	// Validate Common Name for engine
	if newcert.Type == "engine" {
		if !claims.Admin {
			r := jsonapi.NewErrorResponse(12000, "Only members of Admin group can request certificate of type engine")

			return api.JSON(c, http.StatusBadRequest, r)
		}

		if claims.Audience == csr.CR.Subject.CommonName {
			r := jsonapi.NewErrorResponse(12000, "Cannot set CN to username for certificate of type engine")

			return api.JSON(c, http.StatusBadRequest, r)
		}

		// if type is engine, CN should be a FQDN
		if err = validation.IsValidFQDN(csr.CR.Subject.CommonName); err != nil {
			r := jsonapi.NewErrorResponse(12000, "FQDN is not valid")

			return api.JSON(c, http.StatusBadRequest, r)
		}
	}

	// Make sure that certificate not already issued and is valid
	indexDN := csr.SubjectToString()

	if db.Exists(indexDN) {
		e := errors.New(errors.CertStoreError, errors.RecordFound)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusBadRequest, r)
	}

	// Sign CSR
	newCA, err := ca.NewCA(cfg.App.Dir.Root)
	if err != nil {
		e := errors.New(errors.RootError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	caPubKey := csr.GetPublicKey()
	caSubject := csr.GetSubject()
	caSubjectAltNames := csr.GetSubjectAltNames()
	caDate := ca.CreateDate(newcert.Duration)
	caSN, err := newCA.IncrementSerialNumber()
	if err != nil {
		e := errors.New(errors.RootError, errors.IncrementFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	apiHost := host.New(c.Request().URL, c.Request().Host)
	if s.CountConfigKey("api_fqdn") == 1 {
		apiHost = s.GetConfig("api_fqdn")[0].Value
	}

	apiPort := s.GetConfig("api_port")[0].Value

	CrlUrl := fmt.Sprintf("https://%s:%s/crl/CRL.crl", apiHost, apiPort)

	template, err := ca.CreateTemplate(false, caSubject, caSubjectAltNames, caDate, caSN, CrlUrl)
	if err != nil {
		e := errors.New(errors.RootError, errors.Unknown)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	crtDerBytes, err := ca.IssueCertificate(template, newCA.Certificate.Crt, caPubKey, newCA.Key.Private)
	if err != nil {
		e := errors.New(errors.RootError, errors.Unknown)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	certificate, err := pkix.NewCertificateFromDER(crtDerBytes)
	if err != nil {
		e := errors.New(errors.RootError, errors.Unknown)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	crtBytes, err := certificate.ToPEM()
	if err != nil {
		return err
	}

	err = newCA.WriteSerialNumber(caSN)
	if err != nil {
		e := errors.New(errors.RootError, errors.WriteFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	// Save information to CA database
	indexExpireDate := ca.DatabaseDateTimeFormat(caDate.Expire)

	certName := fmt.Sprintf("%x", md5.Sum([]byte(indexDN)))
	certNameFile := fmt.Sprintf("%s.crt", certName)
	certNamePath := path.Join(cfg.CA.Dir.Certs, certNameFile)

	db.New(caSN, indexExpireDate, certNameFile, indexDN)

	err = pkix.ToPEMFile(certNamePath, crtDerBytes, 0444)
	if err != nil {
		e := errors.New(errors.CertificateError, errors.Unknown)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	// Response
	response := jsonapi.NewSuccessResponse(string(crtBytes))

	return api.JSON(c, http.StatusOK, response)
}
