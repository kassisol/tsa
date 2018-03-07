package server

import (
	log "github.com/Sirupsen/logrus"
	mw "github.com/kassisol/tsa/api/server/middleware"
	"github.com/kassisol/tsa/api/server/router/cert"
	"github.com/kassisol/tsa/api/server/router/system"
	"github.com/kassisol/tsa/api/server/router/tenant"
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
	e.Use(mw.Https())
	e.Use(mw.AdminPassword())

	// Index
	e.GET("/", system.IndexHandle)

	// Version
	e.GET("/version", system.ServerVersionHandle)

	// Login
	h := middleware.BasicAuth(mw.Authentication)(system.LoginHandle)
	e.GET("/login", h)

	// System
	jwtConfig := middleware.JWTConfig{
		Skipper:    mw.DefaultSkipper,
		SigningKey: jwk,
	}

	sys := e.Group("/system")
	sys.Use(middleware.JWTWithConfig(jwtConfig))
	sys.Use(mw.Authz("admin"))

	sys.GET("/info", system.InfoHandle)
	sys.PUT("/admin/password", system.AdminChangePasswordHandle)

	sys.GET("/auth", system.AuthListHandle)
	sys.POST("/auth", system.AuthCreateHandle)
	sys.DELETE("/auth/:key", system.AuthDeleteHandle)
	sys.PUT("/auth/enable/:type", system.AuthEnableHandle)
	sys.PUT("/auth/disable", system.AuthDisableHandle)

	// Tenant
	tnt := e.Group("/tenant")

	tnt.GET("/", tenant.ListHandle, middleware.JWTWithConfig(jwtConfig), mw.Authz("authenticated"))
	tnt.POST("/", tenant.CreateHandle, middleware.JWTWithConfig(jwtConfig), mw.Authz("admin"))
	tnt.DELETE("/:id", tenant.DeleteHandle, middleware.JWTWithConfig(jwtConfig), mw.Authz("admin"))

	tnt.GET("/:id/group", tenant.GroupListHandle, middleware.JWTWithConfig(jwtConfig), mw.Authz("admin"))
	tnt.POST("/:id/group", tenant.GroupAddHandle, middleware.JWTWithConfig(jwtConfig), mw.Authz("admin"))
	tnt.DELETE("/:id/group/:gid", tenant.GroupDeleteHandle, middleware.JWTWithConfig(jwtConfig), mw.Authz("admin"))

	// CA public certificate
	tnt.GET("/:id/ca", cert.CAPubHandle, mw.Authz("anonymous"))

	// Revocation file
	tnt.GET("/:id/crl/CRL.crl", cert.CRLHandle, mw.Authz("anonymous"))

	tnt.GET("/:id/cert", cert.ListHandle, middleware.JWTWithConfig(jwtConfig), mw.Authz("user"))
	tnt.POST("/:id/cert", cert.SignHandle, middleware.JWTWithConfig(jwtConfig), mw.Authz("user"))
	tnt.POST("/:id/cert/revoke", cert.RevokeHandle, middleware.JWTWithConfig(jwtConfig), mw.Authz("user"))

	if tls {
		log.Fatal(e.StartTLS(addr, cfg.API.CrtFile, cfg.API.KeyFile))
	} else {
		log.Fatal(e.Start(addr))
	}
}
