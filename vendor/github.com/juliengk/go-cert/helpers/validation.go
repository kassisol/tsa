package helpers

import (
	"fmt"
	"strings"

	"github.com/juliengk/go-utils"
)

var (
	NotValidCACommonName = fmt.Errorf("common name cannot contains the words \"root\", \"intermediate\" and \"ca\"")
	NotValidCAType       = fmt.Errorf("ca type is not valid")
	NotValidCAOU         = fmt.Errorf("organizational unit cannot contains the words \"certificate\" and \"authority\"")
)

func IsValidCACommonName(cn string) error {
	blacklisted := []string{
		"Root",
		"Intermediate",
		"CA",
	}

	words := strings.Split(cn, " ")

	for _, word := range words {
		if utils.StringInSlice(word, blacklisted, true) {
			return NotValidCACommonName
		}
	}

	return nil
}

func IsValidCAType(catype string) error {
	if catype != "root" && catype != "intermediate" {
		return NotValidCAType
	}

	return nil
}

func IsValidCAOrgUnit(ou string) error {
	blacklisted := []string{
		"Certificate",
		"Authority",
	}

	words := strings.Split(ou, " ")

	for _, word := range words {
		if utils.StringInSlice(word, blacklisted, true) {
			return NotValidCAOU
		}
	}

	return nil
}
