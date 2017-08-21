package sqlite

import (
	"github.com/kassisol/tsa/cli/storage/driver"
)

func (c *Config) ListSessions(filter map[string]string) []driver.SessionResult {
	var sessions []Session
	var result []driver.SessionResult

	sql := c.DB.Preload("Server")

	if v, ok := filter["id"]; ok {
		sql = sql.Where("id = ?", v)
	}
	if v, ok := filter["active"]; ok {
		sql = sql.Where("active = ?", v)
	}
	if v, ok := filter["server"]; ok {
		var server Server
		c.DB.Where("name = ?", v).First(&server)
		sql = sql.Where("server_id = ?", server.ID)
	}

	sql.Find(&sessions)

	for _, session := range sessions {
		server := driver.ServerResult{
			ID:          session.Server.ID,
			Name:        session.Server.Name,
			TSAURL:      session.Server.TSAURL,
			Description: session.Server.Description,
		}

		r := driver.SessionResult{
			ID:     session.ID,
			Server: server,
			Active: session.Active,
			Token:  session.Token,
		}

		result = append(result, r)
	}

	return result
}
func (c *Config) AddSession(serverID uint, token string) {
	var server Server
	c.DB.Where("id = ?", serverID).First(&server)

	c.DB.Model(&Session{}).Update("active", false)
	c.DB.Create(&Session{
		Server: server,
		Active: true,
		Token:  token,
	})
}

func (c *Config) RemoveSession(id uint) {
	c.DB.Where("id = ?", id).Delete(Session{})
}

func (c *Config) RemoveAllSessions() {
	c.DB.Delete(Session{})
}

func (c *Config) ActivateSession(id uint, activate bool) {
	c.DB.Model(&Session{}).Update("active", false)

	c.DB.Model(&Session{}).Where("id = ?", id).Update("active", activate)
}
