package client

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/juliengk/stack/client"
	"github.com/juliengk/stack/jsonapi"
)

// Authz
func (c *Config) GetToken(username, password string, ttl int) (string, error) {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   c.Directory.NewAuthz,
	}

	req, err := client.New(cc)
	if err != nil {
		return "", err
	}

	req.HeaderAdd("Accept", "application/json")
	req.SetBasicAuth(username, password)
	req.ValueAdd("ttl", strconv.Itoa(ttl))

	result := req.Get()

	var response jsonapi.Response
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return "", err
	}

	if result.Response.StatusCode != 200 {
		if response.Errors == (jsonapi.ResponseMessage{}) {
			return "", fmt.Errorf("Authorization denied")
		}

		return "", fmt.Errorf(response.Errors.Message)
	}

	token := GetReflectStringValue(response.Data)

	return token, nil
}
