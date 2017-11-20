package ldap

import (
	"crypto/tls"
	"strconv"

	"github.com/juliengk/go-ldapc"
	"github.com/juliengk/go-utils"
	"github.com/kassisol/tsa/api/auth/driver"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/api/types"
	"github.com/kassisol/tsa/pkg/adf"
)

func isMemberOf(userMembers []string, groups []types.ServerConfig) bool {
	allowedGroups := []string{}
	for _, group := range groups {
		allowedGroups = append(allowedGroups, group.Value)
	}

	for _, group := range allowedGroups {
		if utils.StringInSlice(group, userMembers, true) {
			return true
		}
	}

	return false
}

func (c *Config) Login(username, password string) (driver.LoginStatus, error) {
	if err := c.IsConfigOK(); err != nil {
		return driver.Failed, err
	}

	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		return driver.Failed, err
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		return driver.Failed, err
	}
	defer s.End()

	port, _ := strconv.Atoi(s.GetConfig("auth_port")[0].Value)

	ldapclient := &ldapc.Client{
		Protocol: ldapc.LDAP,
		Host:     s.GetConfig("auth_host")[0].Value,
		Port:     port,
		Bind: &ldapc.AuthBind{
			BindDN:       s.GetConfig("auth_bind_username")[0].Value,
			BindPassword: s.GetConfig("auth_bind_password")[0].Value,
			BaseDN:       s.GetConfig("auth_search_base_user")[0].Value,
			Filter:       s.GetConfig("auth_search_filter")[0].Value,
		},
	}

	if s.GetConfig("auth_tls")[0].Value == "true" {
		ldapclient.Protocol = ldapc.LDAPS
		ldapclient.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	entry, err := ldapclient.Authenticate(username, password)
	if err != nil {
		return driver.Failed, err
	}

	userMembers := entry.GetAttributeValues(s.GetConfig("auth_attr_members")[0].Value)

	groupsAdmin := s.GetConfig("auth_group_admin")
	groupsUser := s.GetConfig("auth_group_user")

	if isMemberOf(userMembers, groupsAdmin) {
		return driver.Admin, nil
	}

	if isMemberOf(userMembers, groupsUser) {
		return driver.User, nil
	}

	return driver.None, nil
}
