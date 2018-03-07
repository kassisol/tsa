package sqlite

import (
	"time"
)

type Model struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"created_at"`
}

type ServerConfig struct {
	Model
	Key   string
	Value string
}

type Tenant struct {
	Model
	Name       string  `gorm:"unique;"`
	AuthGroups []Group `gorm:"many2many:tenant_groups;"`
	CA         CertificationAuthority
	CAID       uint
}

type Group struct {
	Model
	Name string
}

type CertificationAuthority struct {
	Model
	Type               string
	Duration           int
	Expire             string
	Country            string
	State              string
	Locality           string
	Organization       string
	OrganizationalUnit string
	CommonName         string
}
