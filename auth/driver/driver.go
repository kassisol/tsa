package driver

type LoginStatus int

const (
	Failed LoginStatus = -1 + iota
	None
	Admin
	User
)

type Auther interface {
	AddConfig(key, value string) error

	Login(username, password string) LoginStatus
}
