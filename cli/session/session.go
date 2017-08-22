package session

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/kassisol/tsa/cli/storage"
	"github.com/kassisol/tsa/cli/storage/driver"
	"github.com/kassisol/tsa/client"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/tsa/pkg/token"
)

var (
	ErrNoActiveSession      = errors.New("No active session")
	ErrActiveSessionExpired = errors.New("Active session is expired")
)

type Config struct {
	storage driver.Storager
}

func New() (*Config, error) {
	cfg := adf.NewServer()
	if err := cfg.Init(); err != nil {
		return nil, err
	}

	s, err := storage.NewDriver("sqlite", cfg.AppDir)
	if err != nil {
		return nil, err
	}

	return &Config{storage: s}, nil
}

func (c *Config) End() {
	c.storage.End()
}

func (c *Config) List() []driver.SessionResult {
	return c.storage.ListSessions(map[string]string{})
}

func (c *Config) Create(server, username, password string, ttl int) error {
	filter := map[string]string{
		"name": server,
	}
	srv := c.storage.ListServers(filter)[0]

	// Check if there is an active session
	filter2 := map[string]string{
		"server": server,
	}
	active := c.storage.ListSessions(filter2)

	// If active session, check if still valid
	if len(active) == 1 {
		exp := c.GetExpire(active[0].Token)
		now := time.Now()

		if exp.After(now) {
			return nil
		}
	}

	// otherwise create new session
	clt, err := client.New(srv.TSAURL)
	if err != nil {
		return err
	}

	// Get TSA URL directories
	err = clt.GetDirectory()
	if err != nil {
		return err
	}

	// Authz
	token, err := clt.GetToken(username, password, ttl)
	if err != nil {
		return err
	}

	c.storage.AddSession(srv.ID, token)

	return nil
}

func (c *Config) Remove(serverID uint) error {
	c.storage.RemoveSession(serverID)

	return nil
}

func (c *Config) Clear() error {
	c.storage.RemoveAllSessions()

	return nil
}

func (c *Config) Use(id uint) error {
	filter := map[string]string{
		"id": strconv.Itoa(int(id)),
	}

	if len(c.storage.ListSessions(filter)) == 0 {
		return fmt.Errorf("Session ID %d is not valid", id)
	}

	c.storage.ActivateSession(id, true)

	return nil
}

func (c *Config) Unuse(id uint) error {
	filter := map[string]string{
		"id": strconv.Itoa(int(id)),
	}

	if len(c.storage.ListSessions(filter)) == 0 {
		return fmt.Errorf("Session ID %d is not valid", id)
	}

	c.storage.ActivateSession(id, false)

	return nil
}

func (c *Config) Get() (driver.SessionResult, error) {
	session := c.getActive()

	if session == (driver.SessionResult{}) {
		return driver.SessionResult{}, ErrNoActiveSession
	}

	if c.Expired() {
		return driver.SessionResult{}, ErrActiveSessionExpired
	}

	return session, nil
}

func (c *Config) Expired() bool {
	session := c.getActive()

	if session == (driver.SessionResult{}) {
		return true
	}

	exp := c.GetExpire(session.Token)
	now := time.Now()

	if exp.Before(now) {
		return true
	}

	return false
}

func (c *Config) Status() string {
	session := c.getActive()

	if session == (driver.SessionResult{}) {
		return "No active session"
	}

	expiresAt := c.GetExpire(session.Token)

	return fmt.Sprintf("Active server is %s and the session expires at %s", session.Server.Name, expiresAt.String())
}

func (c *Config) GetExpire(jwt string) time.Time {
	t := token.New([]byte(""), false)
	claims, _ := t.GetStandardClaims(jwt)

	return time.Unix(claims.ExpiresAt, 0)
}

func (c *Config) GetServer(server string) (driver.ServerResult, error) {
	filter := map[string]string{
		"name": server,
	}
	srv := c.storage.ListServers(filter)

	if len(srv) == 0 {
		return driver.ServerResult{}, fmt.Errorf("Server name '%s' does not exist", server)
	}

	return srv[0], nil
}

func (c *Config) getActive() driver.SessionResult {
	filter := map[string]string{
		"active": "1",
	}

	session := c.storage.ListSessions(filter)

	if len(session) == 0 {
		return driver.SessionResult{}
	}

	return session[0]
}
