package ldap

import (
	"fmt"
	"strings"

	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/pkg/adf"
)

func (c *Config) IsValidConfigKey(key string) error {
	if _, ok := c.Keys[key]; !ok {
		return fmt.Errorf("%s is not a valid config key", key)
	}

	return nil
}

func (c *Config) ValidConfigKeyCount(key string) error {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		return err
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		return err
	}
	defer s.End()

	count := c.Keys[key]

	if s.CountConfigKey(key) >= count.Max {
		return fmt.Errorf("%s is already used", key)
	}

	return nil
}

func (c *Config) IsConfigOK() error {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		return err
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		return err
	}
	defer s.End()

	missingKeys := []string{}

	for key, count := range c.Keys {
		if s.CountConfigKey(key) < count.Min {
			missingKeys = append(missingKeys, key)
		}
	}

	if len(missingKeys) > 0 {
		return fmt.Errorf("Authentication configuration keys MANDATORY: %s", strings.Join(missingKeys, ","))
	}

	return nil
}
