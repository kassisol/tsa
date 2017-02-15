package none

import (
	"github.com/kassisol/tsa/auth"
	"github.com/kassisol/tsa/auth/driver"
)

func init() {
	auth.RegisterDriver("none", New)
}

type Config struct{}

func New() (driver.Auther, error) {
	return &Config{}, nil
}
