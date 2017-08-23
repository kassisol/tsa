package client

import (
	"encoding/json"
	"fmt"

	"github.com/juliengk/stack/client"
	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/api/types"
)

// Get TSA URL directories
func (c *Config) GetDirectory() error {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   "/",
	}

	req, err := client.New(cc)
	if err != nil {
		return err
	}

	req.HeaderAdd("Accept", "application/json")

	result := req.Get()

	var response jsonapi.Response
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return err
	}

	if result.Response.StatusCode != 200 {
		if response.Errors == (jsonapi.ResponseMessage{}) {
			return fmt.Errorf("Problem fetching directory")
		}

		return fmt.Errorf(response.Errors.Message)
	}

	r, err := json.Marshal(response.Data)
	if err != nil {
		return err
	}

	var directory types.Directory
	if err := json.Unmarshal(r, &directory); err != nil {
		return err
	}

	if directory == (types.Directory{}) {
		return fmt.Errorf("Empty Directory")
	}

	c.Directory = directory

	return nil
}
