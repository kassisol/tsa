package validation

import (
	"fmt"

	"github.com/juliengk/go-utils/validation"
)

var (
	ErrCountryCodeLength = fmt.Errorf("Country should be a 2 letters code")
)

func IsValidCountry(country string) error {
	if len(country) != 2 {
		return ErrCountryCodeLength
	}

	for _, c := range country {
		if err := validation.IsUpper(string(c)); err != nil {
			return fmt.Errorf("Country: %v", err)
		}
	}

	return nil
}
