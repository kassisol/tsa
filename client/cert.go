package client

import (
	"encoding/json"
	"fmt"

	"github.com/juliengk/go-cert/ca/database/backend"
	"github.com/juliengk/stack/client"
	"github.com/juliengk/stack/jsonapi"
)

func (c *Config) CertList(token string, filters map[string]string) ([]backend.CertificateResult, error) {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   "/system/cert",
	}

	req, err := client.New(cc)
	if err != nil {
		return nil, err
	}

	req.HeaderAdd("Accept", "application/json")
	req.HeaderAdd("Authorization", fmt.Sprintf("Bearer %s", token))

	for k, v := range filters {
		req.ValueAdd(k, v)
	}

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

	var certificates []backend.CertificateResult
	if err := json.Unmarshal(r, &certificates); err != nil {
		return nil, err
	}

	return certificates, nil
}

func (c *Config) CertRevoke(token string, serialNumber int) error {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   fmt.Sprintf("/system/cert/revoke/%d", serialNumber),
	}

	req, err := client.New(cc)
	if err != nil {
		return err
	}

	req.HeaderAdd("Accept", "application/json")
	req.HeaderAdd("Authorization", fmt.Sprintf("Bearer %s", token))

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
