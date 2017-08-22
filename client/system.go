package client

import (
	"encoding/json"
	"fmt"

	"github.com/juliengk/stack/client"
	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/api/types"
)

func (c *Config) GetInfo(token string) (*types.SystemInfo, error) {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   "/system/info",
	}

	req, err := client.New(cc)
	if err != nil {
		return nil, err
	}

	req.HeaderAdd("Accept", "application/json")
	req.HeaderAdd("Authorization", fmt.Sprintf("Bearer %s", token))

	result := req.Get()

	var response jsonapi.Response
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return nil, err
	}

	if result.Response.StatusCode != 200 {
		if response.Errors == (jsonapi.ResponseMessage{}) {
			return nil, fmt.Errorf("Authorization denied")
		}

		return nil, fmt.Errorf(response.Errors.Message)
	}

	r, err := json.Marshal(response.Data)
	if err != nil {
		return nil, err
	}

	var info types.SystemInfo
	if err := json.Unmarshal(r, &info); err != nil {
		return nil, err
	}

	return &info, nil
}
