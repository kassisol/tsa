package validation

import (
	"fmt"
)

var (
	NotValidCAType = fmt.Errorf("ca type is not valid")
)

func IsValidCAType(catype string) error {
	if catype != "root" {
		return NotValidCAType
	}

	return nil
}
