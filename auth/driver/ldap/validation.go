package ldap

import (
	"fmt"
	"strings"

	"github.com/kassisol/tsa/cli/command"
	"github.com/kassisol/tsa/storage"
)

func (c *Config) IsValidConfigKey(key string) error {
	if _, ok := c.Keys[key]; !ok {
		return fmt.Errorf("%s is not a valid config key", key)
	}

	return nil
}

func (c *Config) ValidConfigKeyCount(key string) error {
	s, err := storage.NewDriver("sqlite", command.DBFilePath)
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
	s, err := storage.NewDriver("sqlite", command.DBFilePath)
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
		return fmt.Errorf("Authentication configuration keys MANDATORY: ", strings.Join(missingKeys, ","))
	}

	return nil
}
