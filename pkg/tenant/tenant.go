package tenant

import (
	"fmt"
	"strconv"

	"github.com/juliengk/go-utils"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/api/types"
	"github.com/kassisol/tsa/pkg/adf"
)

type Config struct {
	Tenant      types.Tenant
	AdminGroups []string
	UserGroups  []string
}

func New() (*Config, error) {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		return nil, err
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		return nil, err
	}
	defer s.End()

	adminGroups := []string{"admin"}
	for _, group := range s.GetConfig("auth_group_admin") {
		adminGroups = append(adminGroups, group.Value)
	}

	config := Config{
		AdminGroups: adminGroups,
	}

	return &config, nil
}

func (c *Config) isMemberOf(groups []string) bool {
	for _, group := range groups {
		if utils.StringInSlice(group, c.UserGroups, true) {
			return true
		}
	}

	return false
}

func (c *Config) SetUserGroups(groups []string) {
	c.UserGroups = groups
}

func (c *Config) SetTenant(id int) error {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		return err
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		return err
	}
	defer s.End()

	filters := map[string]string{
		"id": strconv.Itoa(id),
	}
	tenants := s.ListTenants(filters)

	if len(tenants) == 0 {
		return fmt.Errorf("No Tenant with ID %d", id)
	}

	c.Tenant = tenants[0]

	return nil
}

func (c *Config) GetTenant() types.Tenant {
	return c.Tenant
}

func (c *Config) IsAdmin() bool {
	return c.isMemberOf(c.AdminGroups)
}

func (c *Config) IsMember() bool {
	groups := []string{}
	for _, group := range c.Tenant.AuthGroups {
		groups = append(groups, group.Name)
	}

	return c.isMemberOf(groups)
}
