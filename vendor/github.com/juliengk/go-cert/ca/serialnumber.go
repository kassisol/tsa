package ca

import (
	"io/ioutil"
	"path"
	"strconv"
	"strings"
)

func (ca *CA) ReadSerialNumber() (int, error) {
	caSrlFile := path.Join(ca.RootDir, "ca", "ca.srl")

	snStr, err := ioutil.ReadFile(caSrlFile)
	if err != nil {
		return 0, err
	}

	snInt, err := strconv.Atoi(strings.Trim(string(snStr), "\n"))
	if err != nil {
		return 0, err
	}

	return snInt, nil
}

func (ca *CA) IncrementSerialNumber() (int, error) {
	sn, err := ca.ReadSerialNumber()
	if err != nil {
		return 0, err
	}

	return sn + 1, nil
}

func (ca *CA) WriteSerialNumber(sn int) error {
	caSrlFile := path.Join(ca.RootDir, "ca", "ca.srl")

	snStr := strconv.Itoa(sn) + "\n"

	ioutil.WriteFile(caSrlFile, []byte(snStr), 0644)

	return nil
}
