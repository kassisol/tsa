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
