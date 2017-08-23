package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/juliengk/stack/client"
	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/api/types"
)

func (c *Config) AuthList(token string) ([]types.ServerConfig, error) {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   "/system/auth",
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

	var configs []types.ServerConfig
	if err := json.Unmarshal(r, &configs); err != nil {
		return nil, err
	}

	return configs, nil
}

func (c *Config) AuthCreate(token, key, value string) (*types.ServerConfig, error) {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   "/system/auth",
	}

	req, err := client.New(cc)
	if err != nil {
		return nil, err
	}

	req.HeaderAdd("Accept", "application/json")
	req.HeaderAdd("Content-Type", "application/json")
	req.HeaderAdd("Authorization", fmt.Sprintf("Bearer %s", token))

	config := types.ServerConfig{
		Key:   key,
		Value: value,
	}

	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	result := req.Post(bytes.NewBuffer(data))

	var response jsonapi.Response
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return nil, err
	}

	if result.Response.StatusCode > 201 {
		if response.Errors == (jsonapi.ResponseMessage{}) {
			return nil, fmt.Errorf("Authorization denied")
		}

		return nil, fmt.Errorf(response.Errors.Message)
	}

	r, err := json.Marshal(response.Data)
	if err != nil {
		return nil, err
	}

	var info types.ServerConfig
	if err := json.Unmarshal(r, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

func (c *Config) AuthDelete(token, key, value string) error {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   fmt.Sprintf("/system/auth/%s", key),
	}

	if len(value) == 0 {
		value = "ALL"
	}

	req, err := client.New(cc)
	if err != nil {
		return err
	}

	req.HeaderAdd("Accept", "application/json")
	req.HeaderAdd("Authorization", fmt.Sprintf("Bearer %s", token))

	req.ValueAdd("value", value)

	result := req.Delete()

	if result.Response.StatusCode > 204 {
		var response jsonapi.Response
		if err := json.Unmarshal(result.Body, &response); err != nil {
			return err
		}

		return fmt.Errorf(response.Errors.Message)
	}

	return nil
}

func (c *Config) AuthEnable(token, atype string) error {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   fmt.Sprintf("/system/auth/enable/%s", atype),
	}

	req, err := client.New(cc)
	if err != nil {
		return err
	}

	req.HeaderAdd("Accept", "application/json")
	req.HeaderAdd("Content-Type", "application/json")
	req.HeaderAdd("Authorization", fmt.Sprintf("Bearer %s", token))

	config := types.ServerConfig{
		Key:   "type",
		Value: atype,
	}

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	result := req.Put(bytes.NewBuffer(data))

	if result.Response.StatusCode > 204 {
		var response jsonapi.Response
		if err := json.Unmarshal(result.Body, &response); err != nil {
			return err
		}

		return fmt.Errorf(response.Errors.Message)
	}

	return nil
}

func (c *Config) AuthDisable(token string) error {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   "/system/auth/disable",
	}

	req, err := client.New(cc)
	if err != nil {
		return err
	}

	req.HeaderAdd("Accept", "application/json")
	req.HeaderAdd("Content-Type", "application/json")
	req.HeaderAdd("Authorization", fmt.Sprintf("Bearer %s", token))

	config := types.ServerConfig{
		Key:   "type",
		Value: "none",
	}

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	result := req.Put(bytes.NewBuffer(data))

	if result.Response.StatusCode > 204 {
		var response jsonapi.Response
		if err := json.Unmarshal(result.Body, &response); err != nil {
			return err
		}

		return fmt.Errorf(response.Errors.Message)
	}

	return nil
}
