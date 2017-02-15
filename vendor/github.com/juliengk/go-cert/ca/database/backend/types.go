package backend

type CertificateResult struct {
	StatusFlag        string `gorm:"status_flag"`
	ExpirationDate    string `gorm:"expiration_date"`
	RevocationDate    string `gorm:"revocation_date"`
	RevocationReason  string `gorm:"revocation_reason"`
	SerialNumber      int    `gorm:"unique; serial_number"`
	Filename          string
	DistinguishedName string `gorm:"dn"`
}
