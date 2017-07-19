package helpers

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/juliengk/go-cert/ca"
	"github.com/juliengk/go-cert/pkix"
	"github.com/juliengk/go-utils/validation"
)

func CreateKey(bits int, keyFile string) (*pkix.Key, error) {
	key, err := pkix.NewKey(bits)
	if err != nil {
		return &pkix.Key{}, err
	}

	keyBytes, err := key.ToPEM()
	if err != nil {
		return &pkix.Key{}, err
	}

	err = pkix.ToPEMFile(keyFile, keyBytes, 0400)
	if err != nil {
		return &pkix.Key{}, err
	}

	return key, nil
}

func CreateCSR(country, state, locality, org, ou, cn, email string, altnames []string, key *pkix.Key) (*pkix.CertificateRequest, error) {
	subject := pkix.NewSubject(country, state, locality, org, ou, cn)

	ndn := pkix.NewDNSNames()
	ne := pkix.NewEmails()
	nip := pkix.NewIPs()

	if len(email) > 0 {
		ne.AddEmail(email)
	}

	for _, an := range altnames {
		if err := validation.IsValidIP(an); err == nil {
			nip.AddIP(net.ParseIP(an))
		} else {
			ndn.AddDNS(an)
		}
	}

	ans := pkix.NewSubjectAltNames(*ndn, *ne, *nip)

	csr, err := pkix.NewCertificateRequest(key, subject, ans)
	if err != nil {
		return &pkix.CertificateRequest{}, err
	}

	return csr, nil
}

func CreateCrt(crt []byte, crtFile string) error {
	certificate, err := pkix.NewCertificateFromDER(crt)
	if err != nil {
		return err
	}

	crtBytes, err := certificate.ToPEM()
	if err != nil {
		return err
	}

	err = pkix.ToPEMFile(crtFile, crtBytes, 0444)
	if err != nil {
		return err
	}

	return nil
}

func IssueCrt(csr *pkix.CertificateRequest, duration int, caDir string) ([]byte, error) {
	newCA, err := ca.NewCA(caDir)
	if err != nil {
		return nil, err
	}

	caPubKey := csr.GetPublicKey()
	caSubject := csr.GetSubject()
	caSubjectAltNames := csr.GetSubjectAltNames()
	caDate := ca.CreateDate(duration)
	caSN, err := newCA.IncrementSerialNumber()
	if err != nil {
		return nil, err
	}

	template, err := ca.CreateTemplate(false, caSubject, caSubjectAltNames, caDate, caSN, "")
	if err != nil {
		return nil, err
	}

	crtDerBytes, err := ca.IssueCertificate(template, newCA.Certificate.Crt, caPubKey, newCA.Key.Private)
	if err != nil {
		return nil, err
	}

	err = newCA.WriteSerialNumber(caSN)
	if err != nil {
		return nil, err
	}

	return crtDerBytes, nil
}

func ExpireDateString(notafter time.Time) string {
	year := strconv.Itoa(notafter.Year())
	month := strconv.Itoa(int(notafter.Month()))
	day := strconv.Itoa(notafter.Day())

	if len(month) == 1 {
		month = fmt.Sprintf("0%s", month)
	}
	if len(day) == 1 {
		day = fmt.Sprintf("0%s", day)
	}

	return fmt.Sprintf("%s-%s-%s", year, month, day)
}
