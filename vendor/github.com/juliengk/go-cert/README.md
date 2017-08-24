# Certificate

## Examples
### CA Init
```
package main

import (
	"log"

	"github.com/juliengk/go-cert/ca"
	"github.com/juliengk/go-cert/pkix"
)

func main() {
	rootDir := "/tmp/cert"
	duration := 12

	subject := pkix.NewSubject("Canada", "Quebec", "Montreal", "Example inc.", "IT Department", "ca.example.com")

	ne := pkix.NewEmails()
	ne.AddEmail("cert@example.com")

	ndn := pkix.NewDNSNames()

	altnames := pkix.NewSubjectAltNames(ne, ndn)

	date := ca.CreateDate(duration)
	sn := 1

	template, err := ca.CreateTemplate(true, subject, altnames, date, sn)
	if err != nil {
		log.Fatal(err)
	}

	err = ca.InitCA(rootDir, template)
	if err != nil {
		log.Fatal(err)
	}
}
```

### CA Server
```
package main

import (
	"log"
	"path"

	"github.com/juliengk/go-cert/ca"
	"github.com/juliengk/go-cert/pkix"
)

func main() {
	rootDir := "/tmp/cert"
	duration := 12

	newCA, err := ca.NewCA(rootDir)
	if err != nil {
		log.Fatal(err)
	}

	csr, err := pkix.NewCertificateRequestFromDER(csr.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	caPubKey := csr.GetPublicKey()
	caSubject := csr.GetSubject()
	caSubjectAltNames := csr.GetSubjectAltNames()
	caDate := ca.CreateDate(duration)
	caSN, err := newCA.IncrementSerialNumber()
	if err != nil {
		log.Fatal(err)
	}

	template, err := ca.CreateTemplate(false, caSubject, caSubjectAltNames, caDate, caSN)
	if err != nil {
		log.Fatal(err)
	}

	derBytes, err := ca.IssueCertificate(template, newCA.Certificate.Crt, caPubKey, newCA.Key.Private)
	if err != nil {
		log.Fatal(err)
	}

	// create certificate PEM file
	certificate, err := pkix.NewCertificateFromDER(derBytes)
	if err != nil {
		log.Fatal(err)
	}

	crtBytes, err := certificate.ToPEM()
	if err != nil {
		log.Fatal(err)
	}

	err = pkix.ToPEMFile(crtFile, crtBytes, 0400)
	if err != nil {
		log.Fatal(err)
	}

	// create serial number file
	newCA.WriteSerialNumber(caSN)
}
```

### Client
```
package main

import (
	"log"
	"path"

	"github.com/juliengk/go-cert/pkix"
)

func main() {
	certDir := path.Join("/tmp", "cert", "certs")
	keyFile := path.Join(certDir, "key.pem")
	csrFile := path.Join(certDir, "csr.pem")
	crtFile := path.Join(certDir, "crt.pem")

	// Key pair
	key, err := pkix.NewKey(2048)
	if err != nil {
		log.Fatal(err)
	}

	keyBytes, err := key.ToPEM()
	if err != nil {
		log.Fatal(err)
	}

	err = pkix.ToPEMFile(keyFile, keyBytes, 0400)
	if err != nil {
		log.Fatal(err)
	}

	// CSR
	subject := pkix.NewSubject("Canada", "Quebec", "Montreal", "Example inc.", "IT Department", "www.example.com")

	ne := pkix.NewEmails()
	ne.AddEmail("cert@example.com")

	ndn := pkix.NewDNSNames()

	altnames := pkix.NewSubjectAltNames(ne, ndn)

	csr, err := pkix.NewCertificateRequest(key, subject, altnames)
	if err != nil {
		log.Fatal(err)
	}

	csrBytes, err := csr.ToPEM()
	if err != nil {
		log.Fatal(err)
	}

	err = pkix.ToPEMFile(csrFile, csrBytes, 0400)
	if err != nil {
		log.Fatal(err)
	}
}
```
