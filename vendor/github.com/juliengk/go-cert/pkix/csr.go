package pkix

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
)

type CertificateRequest struct {
	Bytes []byte
	CR    *x509.CertificateRequest
}

func CreateRequestTemplate(pubkey *rsa.PublicKey, subject pkix.Name, altnames AltNames) (*x509.CertificateRequest, error) {
	template := &x509.CertificateRequest{
		PublicKey: pubkey,
		Subject:   subject,
	}

	if len(altnames.DNSNames) > 0 {
		template.DNSNames = altnames.DNSNames
	}

	if len(altnames.EmailAddresses) > 0 {
		template.EmailAddresses = altnames.EmailAddresses
	}

	if len(altnames.IPAddresses) > 0 {
		template.IPAddresses = altnames.IPAddresses
	}

	return template, nil
}

func NewCertificateRequest(key *Key, subject pkix.Name, altnames AltNames) (*CertificateRequest, error) {
	template, err := CreateRequestTemplate(key.Public, subject, altnames)
	if err != nil {
		return nil, err
	}

	derBytes, err := x509.CreateCertificateRequest(rand.Reader, template, key.Private)
	if err != nil {
		return nil, err
	}

	cr := &CertificateRequest{
		Bytes: derBytes,
		CR:    template,
	}

	return cr, nil
}

func NewCertificateRequestFromDER(data []byte) (*CertificateRequest, error) {
	cr, err := x509.ParseCertificateRequest(data)
	if err != nil {
		return nil, err
	}

	return &CertificateRequest{
		Bytes: data,
		CR:    cr,
	}, nil
}

func (cr *CertificateRequest) GetSubject() pkix.Name {
	return cr.CR.Subject
}

func (cr *CertificateRequest) SubjectToString() string {
	sbj := cr.CR.Subject

	return fmt.Sprintf("/C=%s ST=%s L=%s O=%s OU=%s CN=%s", sbj.Country[0], sbj.Province[0], sbj.Locality[0], sbj.Organization[0], sbj.OrganizationalUnit[0], sbj.CommonName)
}

func (cr *CertificateRequest) GetSubjectAltNames() AltNames {
	return AltNames{
		DNSNames:       cr.CR.DNSNames,
		EmailAddresses: cr.CR.EmailAddresses,
		IPAddresses:    cr.CR.IPAddresses,
	}
}

func (cr *CertificateRequest) GetPublicKey() interface{} {
	return cr.CR.PublicKey
}

func (cr *CertificateRequest) ToPEM() ([]byte, error) {
	block := &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: cr.Bytes,
	}

	pemBytes := pem.EncodeToMemory(block)
	if pemBytes == nil {
		return nil, errors.New(string(pemBytes))
	}

	return pemBytes, nil
}
