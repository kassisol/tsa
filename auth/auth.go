package auth

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kassisol/tsa/auth/driver"
)

type Initialize func() (driver.Auther, error)

var initializers = make(map[string]Initialize)

var supportedDrivers = func() string {
	drivers := make([]string, 0, len(initializers))

	for v := range initializers {
		drivers = append(drivers, string(v))
	}

	sort.Strings(drivers)

	return strings.Join(drivers, ",")
}()

func NewDriver(driver string) (driver.Auther, error) {
	if init, exists := initializers[driver]; exists {
		return init()
	}

	return nil, fmt.Errorf("The Auth Driver: %s is not supported. Supported drivers are %s", driver, supportedDrivers)
}

func RegisterDriver(driver string, init Initialize) {
	initializers[driver] = init
}
