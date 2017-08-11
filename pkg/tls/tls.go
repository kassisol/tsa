package tls

import (
	"net"
	"time"

	"github.com/juliengk/go-cert/ca"
	"github.com/juliengk/go-cert/helpers"
	"github.com/juliengk/go-cert/pkix"
	"github.com/juliengk/go-utils/filedir"
	"github.com/juliengk/go-utils/ip"
	"github.com/juliengk/go-utils/validation"
)

type Config struct {
	CN       string
	Duration int
	KeyFile  string
	CertFile string
}

func New(cn string, duration int, keyfile, certfile string) (*Config, error) {
	if err := validation.IsValidFQDN(cn); err != nil {
		return nil, err
	}

	return &Config{
		CN:       cn,
		Duration: duration,
		KeyFile:  keyfile,
		CertFile: certfile,
	}, nil
}

func (c *Config) CertificateExist() bool {
	result := 0

	if filedir.FileExists(c.KeyFile) {
		result++
	}
	if filedir.FileExists(c.CertFile) {
		result += 1
	}

	if result == 2 {
		return true
	}

	return false
}

func (c *Config) IsCertificateExpire() bool {
	certificate, err := pkix.NewCertificateFromPEMFile(c.CertFile)
	if err != nil {
		return true
	}

	now := time.Now()
	notAfter := certificate.Crt.NotAfter

	if now.After(notAfter) {
		return true
	}

	return false
}

func (c *Config) CreateSelfSignedCertificate() error {
	subject := pkix.NewSubject("", "", "", "", "", c.CN)

	ndn := pkix.NewDNSNames()
	ne := pkix.NewEmails()
	nip := pkix.NewIPs()

	ndn.AddDNS("localhost")

	ips := []string{}

	interfaces := ip.New()
	interfaces.Get()

	for _, intf := range interfaces {
		if len(intf.V4) > 0 {
			ips = append(ips, intf.V4[0])
		}
	}

	if len(ips) > 0 {
		for _, i := range ips {
			nip.AddIP(net.ParseIP(i))
		}
	}

	altnames := pkix.NewSubjectAltNames(*ndn, *ne, *nip)

	date := ca.CreateDate(c.Duration)

	if err := helpers.CreateSelfSignedCertificate(c.KeyFile, c.CertFile, 4096, subject, altnames, date); err != nil {
		return err
	}

	return nil
}
