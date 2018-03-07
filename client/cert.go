package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/juliengk/go-cert/ca/database/backend"
	"github.com/juliengk/stack/client"
	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/api/types"
)

func (c *Config) CertList(token string, tenantID int, filters map[string]string) ([]backend.CertificateResult, error) {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   fmt.Sprintf("/tenant/%d/cert", tenantID),
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

	var certificates []backend.CertificateResult
	if err := json.Unmarshal(r, &certificates); err != nil {
		return nil, err
	}

	return certificates, nil
}

// Get Certificate
func (c *Config) CertGet(token string, tenantID int, certType string, csr []byte, duration int) ([]byte, error) {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   fmt.Sprintf("/tenant/%d/cert", tenantID),
	}

	req, err := client.New(cc)
	if err != nil {
		return nil, err
	}

	req.HeaderAdd("Accept", "application/json")
	req.HeaderAdd("Content-Type", "application/json")
	req.HeaderAdd("Authorization", fmt.Sprintf("Bearer %s", token))

	newcert := types.NewCert{
		Type:     certType,
		CSR:      csr,
		Duration: duration,
	}

	data, err := json.Marshal(newcert)
	if err != nil {
		return nil, err
	}

	result := req.Post(bytes.NewBuffer(data))
	if result.Error != nil {
		return nil, result.Error
	}

	var response jsonapi.Response
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return nil, err
	}

	if result.Response.StatusCode != 200 {
		return nil, fmt.Errorf(response.Errors.Message)
	}

	rc := GetReflectStringValue(response.Data)

	return []byte(rc), nil
}

// Revoke Certificate
func (c *Config) CertRevoke(token string, tenantID int, serialNumber int) error {
	revokecert := types.RevokeCert{
		SerialNumber: serialNumber,
	}

	data, err := json.Marshal(revokecert)
	if err != nil {
		return err
	}

	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   fmt.Sprintf("/tenant/%d/cert/revoke", tenantID),
	}

	req, err := client.New(cc)
	if err != nil {
		return err
	}

	req.HeaderAdd("Accept", "application/json")
	req.HeaderAdd("Content-Type", "application/json")
	req.HeaderAdd("Authorization", fmt.Sprintf("Bearer %s", token))

	result := req.Post(bytes.NewBuffer(data))
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
