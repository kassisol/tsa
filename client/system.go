package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/juliengk/stack/client"
	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/api/types"
	"github.com/kassisol/tsa/version"
)

func (c *Config) GetServerVersion() (*version.VersionInfo, error) {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   "/version",
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
			return nil, fmt.Errorf("Authorization denied")
		}

		return nil, fmt.Errorf(response.Errors.Message)
	}

	r, err := json.Marshal(response.Data)
	if err != nil {
		return nil, err
	}

	var ver version.VersionInfo
	if err := json.Unmarshal(r, &ver); err != nil {
		return nil, err
	}

	return &ver, nil
}

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
	if result.Error != nil {
		return nil, result.Error
	}

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

func (c *Config) AdminChangePassword(pold, pnew, pconfirm string) error {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   "/system/admin/password",
	}

	req, err := client.New(cc)
	if err != nil {
		return err
	}

	req.HeaderAdd("Content-Type", "application/json")

	p := types.ChangePassword{
		Old:     pold,
		New:     pnew,
		Confirm: pconfirm,
	}

	data, err := json.Marshal(p)
	if err != nil {
		return err
	}

	result := req.Put(bytes.NewBuffer(data))
	if result.Error != nil {
		return result.Error
	}

	var response jsonapi.Response
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return err
	}

	if result.Response.StatusCode != 200 {
		return fmt.Errorf(response.Errors.Message)
	}

	return nil
}
