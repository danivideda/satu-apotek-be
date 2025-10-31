package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type myClaim struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

const refreshTokenKey string = "refresh-token-secret-key"
const accessTokenKey string = "refresh-token-secret-key"

func generate(id string, exp time.Time, signKey string) (string, error) {
	claims := myClaim{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "satu-apotek-api",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(signKey)
}

func NewRefreshToken(id string) (string, error) {
	exp := time.Now().Add(24 * time.Hour)
	return generate(id, exp, refreshTokenKey)
}

func NewAccessToken(id string) (string, error) {
	exp := time.Now().Add(1 * time.Minute)
	return generate(id, exp, accessTokenKey)
}
