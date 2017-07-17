package none

import (
	"github.com/kassisol/tsa/auth/driver"
)

func (c *Config) Login(username, password string) (driver.LoginStatus, error) {
	return driver.None, nil
}
