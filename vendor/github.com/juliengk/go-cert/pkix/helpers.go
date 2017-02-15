package pkix

import (
	"crypto/x509/pkix"
	"io/ioutil"
	"net"
	"os"
)

type AltNames struct {
	DNSNames       CertDNSNames
	EmailAddresses CertEmails
	IPAddresses    CertIPs
}

type CertDNSNames []string
type CertEmails []string
type CertIPs []net.IP

func NewSubject(country, state, city, o, ou, cn string) pkix.Name {
	return pkix.Name{
		Country:            []string{country},
		Province:           []string{state},
		Locality:           []string{city},
		Organization:       []string{o},
		OrganizationalUnit: []string{ou},
		CommonName:         cn,
	}
}

func NewSubjectAltNames(dnsNames, emailAddresses []string, ipAddresses []net.IP) AltNames {
	return AltNames{
		DNSNames:       dnsNames,
		EmailAddresses: emailAddresses,
		IPAddresses:    ipAddresses,
	}
}

func NewDNSNames() CertDNSNames {
	return CertDNSNames{}
}

func (d CertDNSNames) AddDNS(dns string) {
	d = append(d, dns)
}

func NewEmails() CertEmails {
	return CertEmails{}
}

func (e CertEmails) AddEmail(email string) {
	e = append(e, email)
}

func NewIPs() CertIPs {
	return CertIPs{}
}

func (i CertIPs) AddIP(ip net.IP) {
	i = append(i, ip)
}

func ToPEMFile(path string, pemBytes []byte, mode os.FileMode) error {
	return ioutil.WriteFile(path, pemBytes, mode)
}
