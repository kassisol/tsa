package sqlite

import (
	"fmt"

	"github.com/kassisol/tsa/api/types"
)

func (c *Config) ListTenants(filters map[string]string) []types.Tenant {
	var result []types.Tenant
	var tenants []Tenant

	sql := c.DB

	if v, ok := filters["id"]; ok {
		sql = sql.Where("id = ?", v)
	}

	if v, ok := filters["name"]; ok {
		sql = sql.Where("name = ?", v)
	}

	sql.Find(&tenants)

	for _, tenant := range tenants {
		var groups []types.Group
		var groupsTemp []Group
		var caTemp CertificationAuthority

		c.DB.Model(&tenant).Association("AuthGroups").Find(&groupsTemp)

		for _, g := range groupsTemp {
			groups = append(groups, types.Group{ID: g.ID, CreatedAt: g.CreatedAt, Name: g.Name})
		}

		c.DB.Where("id = ?", tenant.CAID).First(&caTemp)
		ca := types.CertificationAuthority{
			Type:               caTemp.Type,
			Duration:           caTemp.Duration,
			Expire:             caTemp.Expire,
			Country:            caTemp.Country,
			State:              caTemp.State,
			Locality:           caTemp.Locality,
			Organization:       caTemp.Organization,
			OrganizationalUnit: caTemp.OrganizationalUnit,
			CommonName:         caTemp.CommonName,
		}

		tr := types.Tenant{
			ID:         tenant.ID,
			CreatedAt:  tenant.CreatedAt,
			Name:       tenant.Name,
			AuthGroups: groups,
			CA:         ca,
		}

		result = append(result, tr)
	}

	return result
}

func (c *Config) AddTenant(name string, groups []types.Group, caType string, caDuration int, caExpire, caCountry, caState, caLocality, caOrg, caOU, caCN string) error {
	tx := c.DB.Begin()

	grps := []Group{}

	for _, g := range groups {
		grps = append(grps, Group{Name: g.Name})
	}

	ca := CertificationAuthority{
		Type:               caType,
		Duration:           caDuration,
		Expire:             caExpire,
		Country:            caCountry,
		State:              caState,
		Locality:           caLocality,
		Organization:       caOrg,
		OrganizationalUnit: caOU,
		CommonName:         caCN,
	}

	if err := tx.Create(&ca).Error; err != nil {
		tx.Rollback()

		return err
	}

	caR := CertificationAuthority{}
	tx.Where("common_name= ?", caCN).First(&caR)

	tenant := Tenant{
		Name:       name,
		AuthGroups: grps,
		CAID:       caR.ID,
	}

	if err := tx.Create(&tenant).Error; err != nil {
		tx.Rollback()

		return err
	}

	tx.Commit()

	return nil
}

func (c *Config) RemoveTenant(id int) error {
	tx := c.DB.Begin()

	t := Tenant{}
	c.DB.Where("id = ?", id).Find(&t)

	if err := tx.Where("id = ?", t.CAID).Delete(CertificationAuthority{}).Error; err != nil {
		tx.Rollback()

		return err
	}

	if err := tx.Exec("DELETE FROM tenant_groups WHERE tenant_id = ?", id).Error; err != nil {
		tx.Rollback()

		return err
	}

	if err := tx.Where("id = ?", id).Delete(Tenant{}).Error; err != nil {
		tx.Rollback()

		return err
	}

	tx.Commit()

	return nil
}

func (c *Config) AddGroupToTenant(tenant, group string) error {
	t := Tenant{}
	g := Group{}

	tx := c.DB.Begin()

	c.DB.Where("name = ?", tenant).First(&t)

	/*
		if err := tx.Create(&Group{Name: group}).Error; err != nil {
			tx.Rollback()

			return err
		}
		if err := tx.Where("name = ?", group).First(&g).Error; err != nil {
			tx.Rollback()

			return err
		}
	*/
	if err := tx.FirstOrCreate(&g, Group{Name: group}).Error; err != nil {
		tx.Rollback()

		return err
	}

	if err := tx.Model(&t).Association("AuthGroups").Append(&g).Error; err != nil {
		tx.Rollback()

		return err
	}

	tx.Commit()

	return nil
}

func (c *Config) RemoveGroupFromTenant(tenant, group int) error {
	var ct int64
	var cg int64

	t := Tenant{}
	g := Group{}

	c.DB.Where("id = ?", tenant).First(&t).Count(&ct)
	c.DB.Where("id = ?", group).First(&g).Count(&cg)

	if ct == 0 {
		return fmt.Errorf("Tenant ID '%d' does not exist", tenant)
	}

	if cg == 0 {
		return fmt.Errorf("Group ID '%d' does not exist", group)
	}

	if err := c.DB.Model(&t).Association("AuthGroups").Delete(&g).Error; err != nil {
		return err
	}

	return nil
}
