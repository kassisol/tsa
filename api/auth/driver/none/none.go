package none

import (
	"github.com/kassisol/tsa/api/auth"
	"github.com/kassisol/tsa/api/auth/driver"
)

func init() {
	auth.RegisterDriver("none", New)
}

type Config struct{}

func New() (driver.Auther, error) {
	return &Config{}, nil
}
