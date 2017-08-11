package sqlite

import (
	"fmt"

	"github.com/kassisol/tsa/api/types"
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

func (c *Config) GetConfig(key string) []types.ServerConfig {
	var serverconfigs []ServerConfig
	var result []types.ServerConfig

	c.DB.Where("key = ?", key).Find(&serverconfigs)

	for _, config := range serverconfigs {
		r := types.ServerConfig{
			Key:   config.Key,
			Value: config.Value,
		}

		result = append(result, r)
	}

	return result
}

func (c *Config) ListConfigs(prefix string) []types.ServerConfig {
	var serverconfigs []ServerConfig
	var result []types.ServerConfig

	sql := c.DB

	if prefix != "" {
		pre := fmt.Sprintf("%s_%s", prefix, "%")
		sql = sql.Where("key LIKE ?", pre).Order("key")
	}

	sql.Find(&serverconfigs)

	for _, config := range serverconfigs {
		r := types.ServerConfig{
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
