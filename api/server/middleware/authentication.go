package middleware

import (
	log "github.com/Sirupsen/logrus"
	pass "github.com/juliengk/go-utils/password"
	"github.com/kassisol/tsa/api/auth"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/labstack/echo"
)

func Authentication(username, password string, c echo.Context) (bool, error) {
	var groups []string

	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		return false, err
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		log.Fatal(err)
	}
	defer s.End()

	authType := s.GetConfig("auth_type")[0].Value
	a, err := auth.NewDriver(authType)
	if err != nil {
		log.Warning(err)

		if username != "admin" {
			return false, err
		}
	}

	if username == "admin" {
		if pass.ComparePassword([]byte(password), []byte(s.GetConfig("admin_password")[0].Value)) {
			groups = []string{"admin"}
		}
	} else {
		groups, err = a.Login(username, password)
		if err != nil {
			log.Warning(err)

			return false, err
		}
	}

	c.Set("username", username)
	c.Set("groups", groups)

	return true, nil
}
