package token

import (
	"github.com/juliengk/go-utils/random"
)

var defaultLetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()_-+="

func GenerateJWK(letterBytes string, length int) string {
	lb := letterBytes
	if len(letterBytes) == 0 {
		lb = defaultLetterBytes
	}

	return random.RandString(lb, length)
}
