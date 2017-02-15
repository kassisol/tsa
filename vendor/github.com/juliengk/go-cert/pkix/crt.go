package pkix

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"time"
)

type Certificate struct {
	Bytes []byte
	Crt   *x509.Certificate
}

func NewCertificateFromDER(data []byte) (*Certificate, error) {
	crt, err := x509.ParseCertificate(data)
	if err != nil {
		return nil, err
	}

	return &Certificate{
		Bytes: data,
		Crt:   crt,
	}, nil
}

func NewCertificateFromPEM(pemBytes []byte) (*Certificate, error) {
	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil {
		return nil, errors.New("PEM decode failed")
	}

	return NewCertificateFromDER(pemBlock.Bytes)
}

func NewCertificateFromPEMFile(path string) (*Certificate, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return NewCertificateFromPEM(data)
}

func (c *Certificate) ToPEM() ([]byte, error) {
	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: c.Bytes,
	}

	pemBytes := pem.EncodeToMemory(block)
	if pemBytes == nil {
		return nil, errors.New(string(pemBytes))
	}

	return pemBytes, nil
}

func (c *Certificate) IsExpired() bool {
	now := time.Now()

	if now.After(c.Crt.NotAfter) {
		return true
	}

	return false
}
