package driver

type Storager interface {
	ListServers(filter map[string]string) []ServerResult
	AddServer(name, tsaURL, description string)
	RemoveServer(name string)

	ListSessions(filter map[string]string) []SessionResult
	AddSession(serverID uint, token string)
	RemoveSession(id uint)
	RemoveAllSessions()
	ActivateSession(id uint, activate bool)

	End()
}
