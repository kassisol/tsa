package api

type Directory struct {
	CAInfo     string
	NewApp     string
	NewAuthz   string
	RevokeCert string
}

type NewCert struct {
	Type     string
	CSR      []byte
	Duration int
}

type RevokeCert struct {
	SerialNumber int
}
