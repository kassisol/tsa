package driver

import (
	"github.com/kassisol/tsa/api/types"
)

type Storager interface {
	ListConfigs(prefix string) []types.ServerConfig
	AddConfig(key, value string)
	RemoveConfig(key, value string)
	GetConfig(key string) []types.ServerConfig
	CountConfigKey(key string) int

	ListTenants(filters map[string]string) []types.Tenant
	AddTenant(name string, groups []types.Group, caType string, caDuration int, caExpire, caCountry, caState, caLocality, caOrg, caOU, caCN string) error
	RemoveTenant(id int) error
	AddGroupToTenant(tenant, group string) error
	RemoveGroupFromTenant(tenant, group int) error

	End()
}
