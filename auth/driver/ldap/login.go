package ldap

import (
	"crypto/tls"
	"strconv"

	"github.com/juliengk/go-ldapc"
	"github.com/juliengk/go-utils"
	"github.com/kassisol/tsa/auth/driver"
	"github.com/kassisol/tsa/cli/command"
	"github.com/kassisol/tsa/storage"
	sdrv "github.com/kassisol/tsa/storage/driver"
)

func isMemberOf(userMembers []string, groups []sdrv.ServerConfigResult) bool {
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

func (c *Config) Login(username, password string) driver.LoginStatus {
	if err := c.IsConfigOK(); err != nil {
		return driver.Failed
	}

	s, err := storage.NewDriver("sqlite", command.DBFilePath)
	if err != nil {
		return driver.Failed
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
		return driver.Failed
	}

	userMembers := entry.GetAttributeValues(s.GetConfig("auth_attr_members")[0].Value)

	groups_admin := s.GetConfig("auth_group_admin")
	groups_user := s.GetConfig("auth_group_user")

	if isMemberOf(userMembers, groups_admin) {
		return driver.Admin
	}

	if isMemberOf(userMembers, groups_user) {
		return driver.User
	}

	return driver.None
}
