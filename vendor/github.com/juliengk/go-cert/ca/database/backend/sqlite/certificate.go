package sqlite

import (
	"github.com/juliengk/go-cert/ca/database/backend"
)

func (c *Config) New(serialNumber int, expireDate, filename, dn string) {
	if len(filename) == 0 {
		filename = "unknown"
	}

	cert := Certificate{
		StatusFlag:        "V",
		ExpirationDate:    expireDate,
		SerialNumber:      serialNumber,
		Filename:          filename,
		DistinguishedName: dn,
	}

	c.DB.Create(&cert)
}

func (c *Config) UpdateStatus(serialNumber int, status string) {
	c.DB.Where("serial_number = ?", serialNumber).First(&Certificate{}).Update("status_flag", status)
}

func (c *Config) Revoke(serialNumber int, date string, reason int) {
	cert := Certificate{
		StatusFlag:       "R",
		RevocationDate:   date,
		RevocationReason: reason,
	}

	c.DB.Where("serial_number = ?", serialNumber).First(&Certificate{}).Updates(cert)
}

func (c *Config) List(filter map[string]string) []backend.CertificateResult {
	result := []backend.CertificateResult{}

	sql := c.DB.Table("certificates").Select("status_flag, expiration_date, revocation_date, revocation_reason, serial_number, filename, distinguished_name")

	if v, ok := filter["status"]; ok {
		sql = sql.Where("status_flag = ?", v)
	}

	if v, ok := filter["serial"]; ok {
		sql = sql.Where("serial_number = ?", v)
	}

	if v, ok := filter["cn"]; ok {
		sql = sql.Where("distinguished_name LIKE ?", "%CN="+v)
	}

	rows, _ := sql.Rows()
	defer rows.Close()

	for rows.Next() {
		var statusFlag string
		var expirationDate string
		var revocationDate string
		var revocationReason string
		var serialNumber int
		var filename string
		var dn string

		rows.Scan(&statusFlag, &expirationDate, &revocationDate, &revocationReason, &serialNumber, &filename, &dn)

		cr := backend.CertificateResult{
			StatusFlag:        statusFlag,
			ExpirationDate:    expirationDate,
			RevocationDate:    revocationDate,
			RevocationReason:  revocationReason,
			SerialNumber:      serialNumber,
			Filename:          filename,
			DistinguishedName: dn,
		}

		result = append(result, cr)
	}

	return result
}

func (c *Config) Count(status string) int {
	var count int64

	sql := c.DB.Table("certificates")

	if status == "A" {
		sql.Count(&count)
	} else {
		sql.Where("status_flag = ?", status).Count(&count)
	}

	return int(count)
}

func (c *Config) Exists(dn string) bool {
	var count int64

	c.DB.Table("certificates").Where("status_flag = 'V'").Where("distinguished_name = ?", dn).Count(&count)

	if count > 0 {
		return true
	}

	return false
}
