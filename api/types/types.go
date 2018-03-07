package types

import (
	"time"
)

type ServerConfig struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type API struct {
	FQDN        string `json:"fqdn"`
	BindAddress string `json:"bind_address"`
	BindPort    string `json:"bind_port"`
}

type ChangePassword struct {
	Old     string `json:"old"`
	New     string `json:"new"`
	Confirm string `json:"confirm"`
}

type Auth struct {
	Type string `json:"type"`
}

type SystemInfo struct {
	API           API    `json:"api"`
	Auth          Auth   `json:"auth"`
	ServerVersion string `json:"server_version"`
	ID            string `json:"id"`
	StorageDriver string `json:"storage_driver"`
	LoggingDriver string `json:"logging_driver"`
	TSARootDir    string `json:"tsa_root_dir"`
}

type Tenant struct {
	ID         uint                   `json:"id,omitempty"`
	CreatedAt  time.Time              `json:"created_at,omitempty"`
	Name       string                 `json:"name"`
	AuthGroups []Group                `json:"auth_groups"`
	CA         CertificationAuthority `json:"ca"`
}

type Group struct {
	ID        uint      `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Name      string    `json:"name"`
}

type CertificationAuthority struct {
	Type               string `json:"type"`
	Duration           int    `json:"duration;omitempty"`
	Expire             string `json:"expire;omitempty"`
	Country            string `json:"country"`
	State              string `json:"state"`
	Locality           string `json:"locality"`
	Organization       string `json:"organization"`
	OrganizationalUnit string `json:"organizational_unit"`
	CommonName         string `json:"common_name;omitempty"`
}

type TenantGroup struct {
	Tenant string `json:"tenant"`
	Group  string `json:"group"`
}

type CertificateStats struct {
	Certificate int `json:"certificate"`
	Valid       int `json:"valid"`
	Expired     int `json:"expired"`
	Revoked     int `json:"revoked"`
}

type NewCert struct {
	Type     string `json:"type"`
	CSR      []byte `json:"csr"`
	Duration int    `json:"duration"`
}

type RevokeCert struct {
	SerialNumber int `json:"serial_number"`
}
