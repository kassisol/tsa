package client

import (
	"github.com/juliengk/stack/client"
)

type Config struct {
	URL *client.URL
}

func New(url string) (*Config, error) {
	u, err := client.ParseUrl(url)
	if err != nil {
		return nil, err
	}

	return &Config{URL: u}, nil
}
