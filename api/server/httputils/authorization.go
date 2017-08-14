package httputils

import (
	log "github.com/Sirupsen/logrus"
	pass "github.com/juliengk/go-utils/password"
	"github.com/kassisol/tsa/api/auth"
	"github.com/kassisol/tsa/api/auth/driver"
	"github.com/kassisol/tsa/api/config"
	"github.com/kassisol/tsa/api/storage"
	"github.com/labstack/echo"
)

func Authorization(username, password string, c echo.Context) (bool, error) {
	var loginStatus driver.LoginStatus

	s, err := storage.NewDriver("sqlite", config.AppPath)
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
