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

	End()
}
