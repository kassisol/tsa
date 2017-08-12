package none

import (
	"github.com/kassisol/tsa/api/auth/driver"
)

func (c *Config) Login(username, password string) (driver.LoginStatus, error) {
	return driver.None, nil
}
