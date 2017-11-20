package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Config struct {
	JWK            []byte
	ValidSignature bool
}

type MyCustomClaims struct {
	Admin bool `json:"admin"`
	jwt.StandardClaims
}

func New(jwk []byte, valid bool) *Config {
	return &Config{
		JWK:            jwk,
		ValidSignature: valid,
	}
}

func (c *Config) Create(audience, issuer string, admin bool, ttl int) (string, error) {
	now := time.Now()

	expireMinutes := 5
	if admin {
		expireMinutes = 1440
		if ttl > 0 {
			expireMinutes = ttl
		}
	}

	claims := MyCustomClaims{
		admin,
		jwt.StandardClaims{
			Audience:  audience,
			ExpiresAt: now.Add(time.Minute * time.Duration(expireMinutes)).Unix(),
			IssuedAt:  now.Unix(),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(c.JWK)
	if err != nil {
		return "", err
	}

	return ss, nil
}

func (c *Config) GetToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//		if _, ok := token.Method.(*jwt.SigningMethodHS256); !ok {
		//			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		//		}

		return c.JWK, nil
	})

	if c.ValidSignature {
		if err != nil {
			return nil, err
		}
	}

	return token, nil
}

func (c *Config) GetStandardClaims(tokenString string) (jwt.StandardClaims, error) {
	token, err := c.GetToken(tokenString)
	if err != nil {
		return jwt.StandardClaims{}, err
	}

	claims := token.Claims.(jwt.MapClaims)

	cl := jwt.StandardClaims{}

	if v, ok := claims["aud"]; ok {
		cl.Audience = v.(string)
	}
	if v, ok := claims["exp"]; ok {
		cl.ExpiresAt = int64(v.(float64))
	}
	if v, ok := claims["jti"]; ok {
		cl.Id = v.(string)
	}
	if v, ok := claims["iat"]; ok {
		cl.IssuedAt = int64(v.(float64))
	}
	if v, ok := claims["iss"]; ok {
		cl.Issuer = v.(string)
	}
	if v, ok := claims["nbf"]; ok {
		cl.NotBefore = int64(v.(float64))
	}
	if v, ok := claims["sub"]; ok {
		cl.Subject = v.(string)
	}

	return cl, nil
}

func (c *Config) GetCustomClaims(tokenString string) (MyCustomClaims, error) {
	token, err := c.GetToken(tokenString)
	if err != nil {
		return MyCustomClaims{}, err
	}

	claims := token.Claims.(jwt.MapClaims)

	admin := false
	if v, ok := claims["admin"]; ok {
		admin = v.(bool)
	}

	stdclaims, err := c.GetStandardClaims(tokenString)
	if err != nil {
		return MyCustomClaims{}, err
	}

	mcc := MyCustomClaims{admin, stdclaims}

	return mcc, nil
}
