package api

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/tsa/auth"
	"github.com/kassisol/tsa/cli/command"
	"github.com/kassisol/tsa/storage"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func authorization(username, password string, c echo.Context) (bool, error) {
	s, err := storage.NewDriver("sqlite", command.DBFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer s.End()

	a, err := auth.NewDriver(s.GetConfig("auth_type")[0].Value)
	if err != nil {
		log.Warning(err)
	}

	loginStatus, err := a.Login(username, password)
	if err != nil {
		return false, err
	}

	if loginStatus > 0 {
		c.Set("username", username)

		admin := false
		if loginStatus == 1 {
			admin = true
		}
		c.Set("admin", admin)

		return true, nil
	}

	return false, nil
}

func API(jwk []byte, addr string) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Directory
	e.GET("/", IndexHandle)
	// CA Info
	e.GET("/info", InfoHandle)

	// Revocation file
	e.GET("/crl/CRL.crl", CRLHandle)

	// Authz
	h := middleware.BasicAuth(authorization)(AuthzHandle)
	e.GET("/new-authz", h)

	// Restricted
	r := e.Group("/acme")
	r.Use(middleware.JWT(jwk))

	// New certificate
	r.POST("/new-app", NewCertHandle)

	// Revoke
	r.POST("/revoke-cert", RevokeCertHandle)

	e.Logger.Fatal(e.StartTLS(addr, command.ApiCrtFile, command.ApiKeyFile))
}
