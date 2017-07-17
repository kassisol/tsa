package sqlite

import (
	"fmt"

	"github.com/kassisol/tsa/storage/driver"
)

func (c *Config) AddConfig(key, value string) {
	c.DB.Create(&ServerConfig{
		Key:   key,
		Value: value,
	})
}

func (c *Config) RemoveConfig(key, value string) {
	sql := c.DB.Where("key = ?", key)

	if value != "ALL" {
		sql = sql.Where("value = ?", value)
	}

	sql.Delete(ServerConfig{})
}

func (c *Config) GetConfig(key string) []driver.ServerConfigResult {
	var serverconfigs []ServerConfig
	var result []driver.ServerConfigResult

	c.DB.Where("key = ?", key).Find(&serverconfigs)

	for _, config := range serverconfigs {
		r := driver.ServerConfigResult{
			Key:   config.Key,
			Value: config.Value,
		}

		result = append(result, r)
	}

	return result
}

func (c *Config) ListConfigs(prefix string) []driver.ServerConfigResult {
	var serverconfigs []ServerConfig
	var result []driver.ServerConfigResult

	sql := c.DB

	if prefix != "" {
		pre := fmt.Sprintf("%s_%s", prefix, "%")
		sql = sql.Where("key LIKE ?", pre).Order("key")
	}

	sql.Find(&serverconfigs)

	for _, config := range serverconfigs {
		r := driver.ServerConfigResult{
			Key:   config.Key,
			Value: config.Value,
		}

		result = append(result, r)
	}

	return result
}

func (c *Config) CountConfigKey(key string) int {
	var count int64

	c.DB.Model(ServerConfig{}).Where("key = ?", key).Count(&count)

	return int(count)
}
