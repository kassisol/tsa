package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/juliengk/stack/client"
	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/api"
)

type Config struct {
	URL       *client.URL
	Directory api.Directory
}

func New(url string) (*Config, error) {
	u, err := client.ParseUrl(url)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "https" {
		return nil, fmt.Errorf("URL scheme should be https")
	}

	return &Config{URL: u}, nil
}

// Get TSA URL directories
func (c *Config) GetDirectory() error {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   "/",
	}

	req, _ := client.New(cc)

	result := req.Get()

	if result.StatusCode != 200 {
		return fmt.Errorf("Problem fetching directory")
	}

	var response jsonapi.Response
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return err
	}

	directory := GetDirectory(response.Data)

	if directory == (api.Directory{}) {
		return fmt.Errorf("Empty Directory")
	}

	c.Directory = directory

	return nil
}

// Authz
func (c *Config) GetToken(username, password string, ttl int) (string, error) {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   c.Directory.NewAuthz,
	}

	req, _ := client.New(cc)
	req.SetBasicAuth(username, password)
	req.ValueAdd("ttl", strconv.Itoa(ttl))

	result := req.Get()

	if result.StatusCode != 200 {
		return "", fmt.Errorf("Authorization denied")
	}

	var response jsonapi.Response
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return "", err
	}

	token := GetReflectStringValue(response.Data)

	return token, nil
}

// Get CA public Key
func (c *Config) GetCACertificate() ([]byte, error) {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   c.Directory.CAInfo,
	}

	req, _ := client.New(cc)

	result := req.Get()

	if result.StatusCode != 200 {
		return nil, fmt.Errorf("Could not fetch CA public key")
	}

	var response jsonapi.Response
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return nil, err
	}

	info := GetReflectStringValue(response.Data)

	return []byte(info), nil
}

// Get Certificate
func (c *Config) GetCertificate(token string, certType string, csr []byte, duration int) ([]byte, error) {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   c.Directory.NewApp,
	}

	req, _ := client.New(cc)
	req.HeaderAdd("Authorization", fmt.Sprintf("Bearer %s", token))

	newcert := api.NewCert{
		Type:     certType,
		CSR:      csr,
		Duration: duration,
	}

	data, err := json.Marshal(newcert)
	if err != nil {
		return nil, err
	}

	result := req.Post(bytes.NewBuffer(data))

	var response jsonapi.Response
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return nil, err
	}

	if result.StatusCode != 200 {
		return nil, fmt.Errorf(response.Errors.Message)
	}

	rc := GetReflectStringValue(response.Data)

	return []byte(rc), nil
}

// Revoke Certificate
func (c *Config) RevokeCertificate(token string, serialNumber int) error {
	revokecert := api.RevokeCert{
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
		Path:   c.Directory.RevokeCert,
	}

	req, _ := client.New(cc)
	req.HeaderAdd("Authorization", fmt.Sprintf("Bearer %s", token))

	result := req.Post(bytes.NewBuffer(data))

	var response jsonapi.Response
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return err
	}

	if result.StatusCode != 200 {
		return fmt.Errorf(response.Errors.Message)
	}

	return nil
}
