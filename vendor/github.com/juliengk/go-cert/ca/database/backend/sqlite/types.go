package sqlite

import (
	"time"
)

type Model struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"created_at"`
}

type Certificate struct {
	Model
	StatusFlag        string `gorm:"status_flag"`
	ExpirationDate    string `gorm:"expiration_date"`
	RevocationDate    string `gorm:"revocation_date"`
	RevocationReason  int    `gorm:"revocation_reason"`
	SerialNumber      int    `gorm:"unique; serial_number"`
	Filename          string
	DistinguishedName string `gorm:"dn"`
}
