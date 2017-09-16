package system

import (
	"net/http"

	"github.com/juliengk/go-cert/ca/database"
	"github.com/juliengk/go-cert/errors"
	"github.com/juliengk/go-cert/helpers"
	"github.com/juliengk/go-cert/pkix"
	"github.com/juliengk/stack/jsonapi"
	apierr "github.com/kassisol/tsa/api/errors"
	"github.com/kassisol/tsa/api/types"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/kassisol/tsa/pkg/host"
	"github.com/kassisol/tsa/version"
	"github.com/labstack/echo"
)

func InfoHandle(c echo.Context) error {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		e := errors.New(apierr.DatabaseError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}
	defer s.End()

	info := types.SystemInfo{
		ServerVersion:   version.Version,
		StorageDriver:   "sqlite",
		LoggingDriver:   "standard",
		TSARootDir:      cfg.App.Dir.Root,
	}

	if len(s.ListConfigs("ca")) > 0 {
		db, err := database.NewBackend("sqlite", cfg.CA.Dir.Root)
		if err != nil {
			e := errors.New(errors.CertStoreError, errors.ReadFailed)
			r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

			return api.JSON(c, http.StatusInternalServerError, r)
		}
		defer db.End()

		expire := "0000-00-00"

		certificate, err := pkix.NewCertificateFromPEMFile(cfg.CA.TLS.CrtFile)
		if err == nil {
			expire = helpers.ExpireDateString(certificate.Crt.NotAfter)
		}

		ca := types.CertificationAuthority{
			Type:               s.GetConfig("ca_type")[0].Value,
			Expire:             expire,
			Country:            s.GetConfig("ca_country")[0].Value,
			State:              s.GetConfig("ca_state")[0].Value,
			Locality:           s.GetConfig("ca_locality")[0].Value,
			Organization:       s.GetConfig("ca_org")[0].Value,
			OrganizationalUnit: s.GetConfig("ca_ou")[0].Value,
			CommonName:         s.GetConfig("ca_cn")[0].Value,
		}
		info.CA = ca

		stats := types.CertificateStats{
			Certificate: db.Count("A"),
			Valid:       db.Count("V"),
			Expired:     db.Count("E"),
			Revoked:     db.Count("R"),
		}
		info.CertificateStats = stats
	}

	apiHost := host.New(c.Request().URL, c.Request().Host)
	if s.CountConfigKey("api_fqdn") == 1 {
		apiHost = s.GetConfig("api_fqdn")[0].Value
	}

	api2 := types.API{
		FQDN:        apiHost,
		BindAddress: s.GetConfig("api_bind")[0].Value,
		BindPort:    s.GetConfig("api_port")[0].Value,
	}
	info.API = api2

	auth := types.Auth{
		Type: s.GetConfig("auth_type")[0].Value,
	}
	info.Auth = auth

	response := jsonapi.NewSuccessResponse(info)

	return api.JSON(c, http.StatusOK, response)
}
