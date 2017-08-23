package client

import (
	"bytes"
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

	var response jsonapi.Response
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return err
	}

	if result.Response.StatusCode != 200 {
		return fmt.Errorf(response.Errors.Message)
	}

	return nil
}

func (c *Config) CAInit(token, ctype, country, state, locality, org, ou, email string, duration int) error {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   "/system/ca/init",
	}

	req, err := client.New(cc)
	if err != nil {
		return err
	}

	req.HeaderAdd("Content-Type", "application/json")
	req.HeaderAdd("Authorization", fmt.Sprintf("Bearer %s", token))

	ca := types.CertificationAuthority{
		Type:               ctype,
		Duration:           duration,
		Country:            country,
		State:              state,
		Locality:           locality,
		Organization:       org,
		OrganizationalUnit: ou,
		Email:              email,
	}

	data, err := json.Marshal(ca)
	if err != nil {
		return err
	}

	result := req.Post(bytes.NewBuffer(data))

	var response jsonapi.Response
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return err
	}

	if result.Response.StatusCode > 201 {
		if response.Errors == (jsonapi.ResponseMessage{}) {
			return fmt.Errorf("Authorization denied")
		}

		return fmt.Errorf(response.Errors.Message)
	}

	return nil
}
