package ldap

import (
	"fmt"

	"github.com/kassisol/tsa/cli/command"
	"github.com/kassisol/tsa/storage"
)

func (c *Config) AddConfig(key, value string) error {
	s, err := storage.NewDriver("sqlite", command.DBFilePath)
	if err != nil {
		return err
	}
	defer s.End()

	// Input validations
	// IV - Key
	if err = c.IsValidConfigKey(key); err != nil {
		return err
	}

	if err = c.ValidConfigKeyCount(key); err != nil {
		return err
	}

	// IV - Value
	values := s.GetConfig(key)
	for _, v := range values {
		if v.Value == value {
			return fmt.Errorf("Key \"%s\" already has value \"%s\"", key, value)
		}
	}

	// Add data to DB
	s.AddConfig(key, value)

	return nil
}
