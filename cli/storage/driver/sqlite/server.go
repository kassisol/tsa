package sqlite

import (
	"github.com/kassisol/tsa/cli/storage/driver"
)

func (c *Config) AddServer(name, tsaURL, description string) {
	c.DB.Create(&Server{
		Name:        name,
		TSAURL:      tsaURL,
		Description: description,
	})
}

func (c *Config) RemoveServer(name string) {
	c.DB.Where("name = ?", name).Delete(Server{})
}

func (c *Config) ListServers(filter map[string]string) []driver.ServerResult {
	var result []driver.ServerResult
	var servers []Server

	sql := c.DB

	if v, ok := filter["name"]; ok {
		sql = sql.Where("name = ?", v)
	}

	sql.Find(&servers)

	for _, server := range servers {
		r := driver.ServerResult{
			ID:          server.ID,
			Name:        server.Name,
			TSAURL:      server.TSAURL,
			Description: server.Description,
		}

		result = append(result, r)
	}

	return result
}
