package ca

import (
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"path"

	"github.com/juliengk/go-cert/pkix"
	"github.com/juliengk/go-utils/filedir"
)

type CA struct {
	RootDir     string
	Key         *pkix.Key
	Certificate *pkix.Certificate
}

func InitCA(rootDir string, template *x509.Certificate) (*CA, error) {
	caDir := path.Join(rootDir, "ca")
	privateDir := path.Join(caDir, "private")
	certsDir := path.Join(caDir, "certs")
	crlDir := path.Join(caDir, "crl")

	caKeyFile := path.Join(privateDir, "ca.key")
	caCrtFile := path.Join(certsDir, "ca.crt")

	if err := filedir.CreateDirIfNotExist(caDir, 0755); err != nil {
		return nil, err
	}

	if err := filedir.CreateDirIfNotExist(certsDir, 0755); err != nil {
		return nil, err
	}

	if err := filedir.CreateDirIfNotExist(crlDir, 0755); err != nil {
		return nil, err
	}

	if err := filedir.CreateDirIfNotExist(privateDir, 0755); err != nil {
		return nil, err
	}

	if !filedir.FileExists(caKeyFile) {
		newCA := &CA{
			RootDir: rootDir,
		}

		// generate private key
		key, err := pkix.NewKey(4096)
		if err != nil {
			return nil, err
		}

		keyBytes, err := key.ToPEM()
		if err != nil {
			return nil, err
		}

		err = pkix.ToPEMFile(caKeyFile, keyBytes, 0400)
		if err != nil {
			return nil, err
		}

		// generate self-signed certificate
		parent := template

		derBytes, err := IssueCertificate(template, parent, key.Public, key.Private)
		if err != nil {
			return nil, err
		}

		// create certificate PEM file
		certificate, err := pkix.NewCertificateFromDER(derBytes)
		if err != nil {
			return nil, err
		}

		crtBytes, err := certificate.ToPEM()
		if err != nil {
			return nil, err
		}

		err = pkix.ToPEMFile(caCrtFile, crtBytes, 0400)
		if err != nil {
			return nil, err
		}

		// create serial number file
		newCA.WriteSerialNumber(int(certificate.Crt.SerialNumber.Int64()))

		newCA.Key = key
		newCA.Certificate = certificate

		return newCA, nil
	}

	return nil, nil
}

func NewCA(rootDir string) (*CA, error) {
	caDir := path.Join(rootDir, "ca")
	privateDir := path.Join(caDir, "private")
	certsDir := path.Join(caDir, "certs")

	caKeyFile := path.Join(privateDir, "ca.key")
	caCrtFile := path.Join(certsDir, "ca.crt")

	if !filedir.FileExists(caKeyFile) && !filedir.FileExists(caCrtFile) {
		return nil, fmt.Errorf("CA key and/or certificate do not exist")
	}

	key, err := pkix.NewKeyFromPEMFile(caKeyFile)
	if err != nil {
		return nil, err
	}

	certificate, err := pkix.NewCertificateFromPEMFile(caCrtFile)
	if err != nil {
		return nil, err
	}

	return &CA{
		RootDir:     rootDir,
		Key:         key,
		Certificate: certificate,
	}, nil
}

func IssueCertificate(template, parent *x509.Certificate, publickey, privatekey interface{}) ([]byte, error) {
	cert, err := x509.CreateCertificate(rand.Reader, template, parent, publickey, privatekey)
	if err != nil {
		return nil, err
	}

	return cert, nil
}
