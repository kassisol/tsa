package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/juliengk/stack/client"
	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/api/types"
)

func (c *Config) TenantList(token string, filters map[string]string) ([]types.Tenant, error) {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   "/tenant/",
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

	var tenants []types.Tenant
	if err := json.Unmarshal(r, &tenants); err != nil {
		return nil, err
	}

	return tenants, nil
}

func (c *Config) TenantCreate(token, name string, authGroups []string, ctype, country, state, locality, org, ou string, duration int) (types.Tenant, error) {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   "/tenant/",
	}

	req, err := client.New(cc)
	if err != nil {
		return types.Tenant{}, err
	}

	req.HeaderAdd("Accept", "application/json")
	req.HeaderAdd("Content-Type", "application/json")
	req.HeaderAdd("Authorization", fmt.Sprintf("Bearer %s", token))

	groups := []types.Group{}
	for _, g := range authGroups {
		groups = append(groups, types.Group{Name: g})
	}

	ca := types.CertificationAuthority{
		Type:               ctype,
		Duration:           duration,
		Country:            country,
		State:              state,
		Locality:           locality,
		Organization:       org,
		OrganizationalUnit: ou,
	}

	tenant := types.Tenant{
		Name:       name,
		AuthGroups: groups,
		CA:         ca,
	}

	data, err := json.Marshal(tenant)
	if err != nil {
		return types.Tenant{}, err
	}

	result := req.Post(bytes.NewBuffer(data))
	if result.Error != nil {
		return types.Tenant{}, result.Error
	}

	var response jsonapi.Response
	if err := json.Unmarshal(result.Body, &response); err != nil {
		return types.Tenant{}, err
	}

	if result.Response.StatusCode > 201 {
		if response.Errors == (jsonapi.ResponseMessage{}) {
			return types.Tenant{}, fmt.Errorf("Authorization denied")
		}

		return types.Tenant{}, fmt.Errorf(response.Errors.Message)
	}

	r, err := json.Marshal(response.Data)
	if err != nil {
		return types.Tenant{}, err
	}

	var tenantR types.Tenant
	if err := json.Unmarshal(r, &tenantR); err != nil {
		return types.Tenant{}, err
	}

	return tenantR, nil
}

func (c *Config) TenantDelete(token string, id int) error {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   fmt.Sprintf("/tenant/%d", id),
	}

	req, err := client.New(cc)
	if err != nil {
		return err
	}

	req.HeaderAdd("Accept", "application/json")
	req.HeaderAdd("Authorization", fmt.Sprintf("Bearer %s", token))

	result := req.Delete()
	if result.Error != nil {
		return result.Error
	}

	if result.Response.StatusCode > 204 {
		var response jsonapi.Response
		if err := json.Unmarshal(result.Body, &response); err != nil {
			return err
		}

		return fmt.Errorf(response.Errors.Message)
	}

	return nil
}

func (c *Config) AddGroupToTenant(token string, id int, tenant, group string) error {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   fmt.Sprintf("/tenant/%d", id),
	}

	req, err := client.New(cc)
	if err != nil {
		return err
	}

	req.HeaderAdd("Accept", "application/json")
	req.HeaderAdd("Content-Type", "application/json")
	req.HeaderAdd("Authorization", fmt.Sprintf("Bearer %s", token))

	ig := types.TenantGroup{
		Tenant: tenant,
		Group:  group,
	}

	data, err := json.Marshal(ig)
	if err != nil {
		return err
	}

	result := req.Post(bytes.NewBuffer(data))
	if result.Error != nil {
		return result.Error
	}

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

func (c *Config) RemoveGroupFromTenant(token string, tenant, group int) error {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   fmt.Sprintf("/tenant/%d/group/%d", tenant, group),
	}

	req, err := client.New(cc)
	if err != nil {
		return err
	}

	req.HeaderAdd("Accept", "application/json")
	req.HeaderAdd("Authorization", fmt.Sprintf("Bearer %s", token))

	result := req.Delete()
	if result.Error != nil {
		return result.Error
	}

	if result.Response.StatusCode > 204 {
		var response jsonapi.Response
		if err := json.Unmarshal(result.Body, &response); err != nil {
			return err
		}

		return fmt.Errorf(response.Errors.Message)
	}

	return nil
}
