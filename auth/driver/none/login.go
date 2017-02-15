package none

import (
	"github.com/kassisol/tsa/auth/driver"
)

func (c *Config) Login(username, password string) driver.LoginStatus {
	return driver.None
}
