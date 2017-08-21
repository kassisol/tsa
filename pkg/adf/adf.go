// ADF (Application Directory and File)

package adf

func NewClient() *ClientConfig {
	return &ClientConfig{}
}

func NewDaemon() *DaemonConfig {
	return &DaemonConfig{}
}

func NewEngine() *EngineConfig {
	return &EngineConfig{}
}

func NewServer() *ServerConfig {
	return &ServerConfig{}
}
