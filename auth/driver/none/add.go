package none

import (
	"fmt"
)

func (c *Config) AddConfig(key, value string) error {
	return fmt.Errorf("Method not valid for \"none\" auth driver")
}
