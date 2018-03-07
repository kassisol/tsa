package client

import (
	"encoding/json"
	"fmt"

	"github.com/juliengk/stack/client"
	"github.com/juliengk/stack/jsonapi"
)

// Get CA public Key
func (c *Config) GetCACertificate(tenantID int) ([]byte, error) {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   fmt.Sprintf("/tenant/%d/ca", tenantID),
	}

	req, err := client.New(cc)
	if err != nil {
		return nil, err
	}

	req.HeaderAdd("Accept", "application/json")

	result := req.Get()
	if result.Error != nil {
		return nil, result.Error
	}

	var response jsonapi.Response
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return nil, err
	}

	if result.Response.StatusCode != 200 {
		if response.Errors == (jsonapi.ResponseMessage{}) {
			return nil, fmt.Errorf("Could not fetch CA public key")
		}

		return nil, fmt.Errorf(response.Errors.Message)
	}

	info := GetReflectStringValue(response.Data)

	return []byte(info), nil
}
