package validation

import (
	"fmt"
)

var (
	ErrNotValidCAType = fmt.Errorf("ca type is not valid")
)

func IsValidCAType(catype string) error {
	if catype != "root" {
		return ErrNotValidCAType
	}

	return nil
}
