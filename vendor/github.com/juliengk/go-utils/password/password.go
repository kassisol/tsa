package password

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func ValidatePassword(p string) bool {
	var passwordRegexp = regexp.MustCompile(`^[a-zA-Z0-9!@#$%\*\(\)\-_]{1,24}$`)

	result := passwordRegexp.MatchString(p)

	return result
}

func GeneratePassword(rawpassword string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(rawpassword), 10)

	return string(hashedPassword)
}

func ComparePassword(rawpassword, password []byte) bool {
	if err := bcrypt.CompareHashAndPassword(password, rawpassword); err != nil {
		return false
	}

	return true
}
