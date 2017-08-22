package server

import (
	log "github.com/Sirupsen/logrus"
	mw "github.com/kassisol/tsa/api/server/middleware"
	"github.com/kassisol/tsa/api/server/router/acme"
	"github.com/kassisol/tsa/api/server/router/ca"
	"github.com/kassisol/tsa/api/server/router/crl"
	"github.com/kassisol/tsa/api/server/router/system"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func API(addr string, tls bool) {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		log.Fatal(err)
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		log.Fatal(err)
	}
	defer s.End()

	jwk := []byte(s.GetConfig("jwk")[0].Value)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(mw.AdminPassword())
	e.Use(mw.CAInit())

	// Directory
	e.GET("/", system.IndexHandle)

	// Authz
	h := middleware.BasicAuth(mw.Authorization)(system.AuthzHandle)
	e.GET("/new-authz", h)

	// System
	jwtConfig := middleware.JWTConfig{
		Skipper:    mw.DefaultSkipper,
		SigningKey: jwk,
	}

	sys := e.Group("/system")
	sys.Use(middleware.JWTWithConfig(jwtConfig))
	sys.Use(mw.AdminOnly())

	sys.GET("/info", system.InfoHandle)

	// CA public certificate
	e.GET("/ca", ca.PubCertHandle)

	// Revocation file
	e.GET("/crl/CRL.crl", crl.CRLHandle)

	// ACME
	r := e.Group("/acme")
	r.Use(middleware.JWT(jwk))

	// New certificate
	r.POST("/new-app", acme.NewCertHandle)

	// Revoke
	r.POST("/revoke-cert", acme.RevokeCertHandle)

	if tls {
		log.Fatal(e.StartTLS(addr, cfg.API.CrtFile, cfg.API.KeyFile))
	} else {
		log.Fatal(e.Start(addr))
	}
}
