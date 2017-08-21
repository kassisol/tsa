package httputils

import (
	"fmt"

	"github.com/kassisol/tsa/api/errors"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/pkg/adf"
)

func GetTokenSigningKey() ([]byte, error) {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		return nil, err
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		e := errors.New(errors.DatabaseError, errors.ReadFailed)

		return []byte(""), fmt.Errorf(e.Message)
	}
	defer s.End()

	return []byte(s.GetConfig("jwk")[0].Value), nil
}
