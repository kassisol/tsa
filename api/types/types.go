package types

type Directory struct {
	CAInfo     string `json:"ca_info"`
	NewApp     string `json:"new_app"`
	NewAuthz   string `json:"new_authz"`
	RevokeCert string `json:"revoke_cert"`
}

type NewCert struct {
	Type     string `json:"type"`
	CSR      []byte `json:"csr"`
	Duration int    `json:"duration"`
}

type RevokeCert struct {
	SerialNumber int `json:"serial_number"`
}

type ServerConfig struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SystemInfo struct {
	CA               CertificationAuthority `json:"certification_authority"`
	CertificateStats CertificateStats       `json:"certificate_stats"`
	API              API                    `json:"api"`
	Auth             Auth                   `json:"auth"`
	ServerVersion    string                 `json:"server_version"`
	StorageDriver    string                 `json:"storage_driver"`
	LoggingDriver    string                 `json:"logging_driver"`
	TSARootDir       string                 `json:"tsa_root_dir"`
}

type CertificationAuthority struct {
	Type               string `json:"type"`
	Duration           int    `json:"duration;omitempty"`
	Expire             string `json:"expire;omitempty"`
	Country            string `json:"country"`
	State              string `json:"state"`
	Locality           string `json:"locality"`
	Organization       string `json:"organization"`
	OrganizationalUnit string `json:"organizatinal_unit"`
	CommonName         string `json:"common_name;omitempty"`
	Email              string `json:"email"`
}

type CertificateStats struct {
	Certificate int `json:"certificate"`
	Valid       int `json:"valid"`
	Expired     int `json:"expired"`
	Revoked     int `json:"revoked"`
}

type API struct {
	FQDN        string `json:"fqdn"`
	BindAddress string `json:"bind_address"`
	BindPort    string `json:"bind_port"`
}

type Auth struct {
	Type string `json:"type"`
}

type ChangePassword struct {
	Old     string `json:"old"`
	New     string `json:"new"`
	Confirm string `json:"confirm"`
}
