package pkix

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/juliengk/go-utils/filedir"
)

type Key struct {
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
	Bytes   []byte
}

func NewKey(bits int) (*Key, error) {
	privatekey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return &Key{}, err
	}

	return &Key{
		Private: privatekey,
		//Public: &public,
		Public: &privatekey.PublicKey,
	}, nil
}

func NewEmptyKey() *Key {
	return &Key{}
}

func (k *Key) ToDER() ([]byte, error) {
	derBytes := x509.MarshalPKCS1PrivateKey(k.Private)
	if derBytes == nil {
		return nil, errors.New("Marshal RSA failed")
	}

	return derBytes, nil
}

func (k *Key) ToPEM() ([]byte, error) {
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(k.Private),
	}

	pemBytes := pem.EncodeToMemory(block)
	if pemBytes == nil {
		return nil, errors.New(string(pemBytes))
	}

	return pemBytes, nil
}

func NewKeyFromPEMFile(path string) (*Key, error) {
	if !filedir.FileExists(path) {
		return nil, errors.New(fmt.Sprintf("%s does not exist", path))
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("No PEM found")
	}

	if block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("PEM Type is not the one expected")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &Key{
		Private: privateKey,
		Bytes:   block.Bytes,
	}, nil
}
