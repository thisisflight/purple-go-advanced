package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	Secret string
}

type JWTData struct {
	PhoneNumber string
}

func NewJWT(secret string) *JWT {
	return &JWT{Secret: secret}
}

func (j *JWT) Create(phoneNumber string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phoneNumber": phoneNumber,
	})
	s, err := t.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return s, nil
}

func (j *JWT) Parse(token string) (bool, *JWTData) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		return false, nil
	}
	phoneNumber := t.Claims.(jwt.MapClaims)["phoneNumber"]
	return t.Valid, &JWTData{PhoneNumber: phoneNumber.(string)}
}
