package ldap

import (
	"crypto/tls"
	"strconv"

	"github.com/juliengk/go-ldapc"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/pkg/adf"
)

func (c *Config) Login(username, password string) ([]string, error) {
	if err := c.IsConfigOK(); err != nil {
		return []string{}, err
	}

	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		return []string{}, err
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		return []string{}, err
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
		return []string{}, err
	}

	userMembers := entry.GetAttributeValues(s.GetConfig("auth_attr_members")[0].Value)

	return userMembers, nil
}
