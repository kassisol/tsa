package ldap

import (
	"github.com/kassisol/tsa/auth"
	"github.com/kassisol/tsa/auth/driver"
)

func init() {
	auth.RegisterDriver("ldap", New)
}

type Count struct {
	Min int
	Max int
}

type ConfigKeys map[string]Count

type Config struct {
	Keys ConfigKeys
}

func NewConfigKeys() ConfigKeys {
	return ConfigKeys{
		"auth_type":             Count{Min: 1, Max: 1},
		"auth_host":             Count{Min: 1, Max: 1},
		"auth_port":             Count{Min: 1, Max: 1},
		"auth_tls":              Count{Min: 1, Max: 1},
		"auth_bind_username":    Count{Min: 1, Max: 1},
		"auth_bind_password":    Count{Min: 1, Max: 1},
		"auth_search_base_user": Count{Min: 1, Max: 1},
		"auth_search_filter":    Count{Min: 1, Max: 1},
		"auth_attr_members":     Count{Min: 1, Max: 1},
		"auth_group_admin":      Count{Min: 1, Max: 100},
		"auth_group_user":       Count{Min: 1, Max: 100},
	}
}

func New() (driver.Auther, error) {
	return &Config{
		Keys: NewConfigKeys(),
	}, nil
}
