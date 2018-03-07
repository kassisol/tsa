package token

import (
	"fmt"

	"github.com/juliengk/go-utils/random"
)

var defaultLetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()_-+="

func Generate(letterBytes string, length int) string {
	lb := letterBytes
	if len(letterBytes) == 0 {
		lb = defaultLetterBytes
	}

	return random.RandString(lb, length)
}

func JWTFromHeader(authHeader string, authScheme string) (string, error) {
	l := len(authScheme)

	if len(authHeader) > l+1 && authHeader[:l] == authScheme {
		return authHeader[l+1:], nil
	}

	return "", fmt.Errorf("Missing or invalid jwt in the request header")
}
