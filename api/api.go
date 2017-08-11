package api

import (
	log "github.com/Sirupsen/logrus"
	pass "github.com/juliengk/go-utils/password"
	mw "github.com/kassisol/tsa/api/middleware"
	"github.com/kassisol/tsa/auth"
	"github.com/kassisol/tsa/auth/driver"
	"github.com/kassisol/tsa/cli/command"
	"github.com/kassisol/tsa/storage"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func authorization(username, password string, c echo.Context) (bool, error) {
	var loginStatus driver.LoginStatus

	s, err := storage.NewDriver("sqlite", command.DBFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer s.End()

	authType := s.GetConfig("auth_type")[0].Value
	if authType == "none" {
		log.Warning("No authentication configured")
	}

	if username == "admin" {
		if pass.ComparePassword([]byte(password), []byte(s.GetConfig("admin_password")[0].Value)) {
			loginStatus = 1
		}
	} else {
		a, err := auth.NewDriver(authType)
		if err != nil {
			log.Warning(err)
		}

		loginStatus, err = a.Login(username, password)
		if err != nil {
			log.Warning(err)

			return false, err
		}
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

func API(addr string, tls bool) {
	s, err := storage.NewDriver("sqlite", command.DBFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer s.End()

	jwk := []byte(s.GetConfig("jwk")[0].Value)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(mw.AdminPassword())

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

	if tls {
		log.Fatal(e.StartTLS(addr, command.ApiCrtFile, command.ApiKeyFile))
	} else {
		log.Fatal(e.Start(addr))
	}
}
