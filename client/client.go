package client

import (
	"github.com/juliengk/stack/client"
	"github.com/kassisol/tsa/api/types"
)

type Config struct {
	URL       *client.URL
	Directory types.Directory
}

func New(url string) (*Config, error) {
	u, err := client.ParseUrl(url)
	if err != nil {
		return nil, err
	}

	return &Config{URL: u}, nil
}
