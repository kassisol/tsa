package driver

type Storager interface {
	ListConfigs(prefix string) []ServerConfigResult
	AddConfig(key, value string)
	RemoveConfig(key, value string)
	GetConfig(key string) []ServerConfigResult
	CountConfigKey(key string) int

	End()
}
