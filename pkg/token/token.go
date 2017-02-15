package token

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kassisol/tsa/cli/command"
	"github.com/kassisol/tsa/errors"
	"github.com/kassisol/tsa/storage"
	"github.com/labstack/echo"
)

type MyCustomClaims struct {
	Admin bool `json:"admin"`
	jwt.StandardClaims
}

func New(jwk []byte, audience string, admin bool) (string, error) {
	now := time.Now()

	claims := MyCustomClaims{
		admin,
		jwt.StandardClaims{
			Audience:  audience,
			ExpiresAt: now.Add(time.Minute * 5).Unix(),
			IssuedAt:  now.Unix(),
			Issuer:    "harbormaster",
		},
	}
	/*claims := jwt.MapClaims{
		"aud": audience,
		"exp": now.Add(time.Minute * 5).Unix(),
		"iat": now.Unix(),
		"iss": "harbormaster",
		"admin": admin,
	}*/

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(jwk)
	if err != nil {
		return "", err
	}

	return ss, nil
}

func JWTFromHeader(c echo.Context, header string, authScheme string) (string, error) {
	auth := c.Request().Header.Get(header)
	l := len(authScheme)

	if len(auth) > l+1 && auth[:l] == authScheme {
		return auth[l+1:], nil
	}

	return "", fmt.Errorf("Missing or invalid jwt in the request header")
}

func GetSigningKey() ([]byte, error) {
	s, err := storage.NewDriver("sqlite", command.DBFilePath)
	if err != nil {
		e := errors.New(errors.DatabaseError, errors.ReadFailed)

		return []byte(""), fmt.Errorf(e.Message)
	}
	defer s.End()

	return []byte(s.GetConfig("jwk")[0].Value), nil
}

func GetToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		jwk, err := GetSigningKey()
		if err != nil {
			return "", err
		}

		return jwk, nil
	})

	if err != nil {
		return &jwt.Token{}, fmt.Errorf("Token not valid")
	}

	return token, nil
}

func GetStandardClaims(tokenString string) (jwt.StandardClaims, error) {
	token, err := GetToken(tokenString)
	if err != nil {
		return jwt.StandardClaims{}, err
	}

	claims := token.Claims.(jwt.MapClaims)

	c := jwt.StandardClaims{}

	if v, ok := claims["aud"]; ok {
		c.Audience = v.(string)
	}
	if v, ok := claims["exp"]; ok {
		c.ExpiresAt = int64(v.(float64))
	}
	if v, ok := claims["jti"]; ok {
		c.Id = v.(string)
	}
	if v, ok := claims["iat"]; ok {
		c.IssuedAt = int64(v.(float64))
	}
	if v, ok := claims["iss"]; ok {
		c.Issuer = v.(string)
	}
	if v, ok := claims["nbf"]; ok {
		c.NotBefore = int64(v.(float64))
	}
	if v, ok := claims["sub"]; ok {
		c.Subject = v.(string)
	}

	return c, nil
}

func GetCustomClaims(tokenString string) (MyCustomClaims, error) {
	token, err := GetToken(tokenString)
	if err != nil {
		return MyCustomClaims{}, err
	}

	claims := token.Claims.(jwt.MapClaims)

	admin := false
	if v, ok := claims["admin"]; ok {
		admin = v.(bool)
	}

	stdclaims, err := GetStandardClaims(tokenString)
	if err != nil {
		return MyCustomClaims{}, err
	}

	c := MyCustomClaims{admin, stdclaims}

	return c, nil
}
