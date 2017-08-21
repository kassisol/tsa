package driver

type ServerResult struct {
	ID          uint
	Name        string
	TSAURL      string
	Description string
}

type SessionResult struct {
	ID     uint
	Server ServerResult
	Active bool
	Token  string
}
