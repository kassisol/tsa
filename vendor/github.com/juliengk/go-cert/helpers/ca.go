package helpers

import (
	"fmt"
)

func UpdateOrgUnitLabel(name string) string {
	return fmt.Sprintf("%s Certificate Authority", name)
}

func UpdateCommonNameLabel(ctype, name string) string {
	if ctype == "root" {
		return fmt.Sprintf("%s Root CA", name)
	}

	if ctype == "intermediate" {
		return fmt.Sprintf("%s Intermediate CA", name)
	}

	return name
}
