package jwt

import (
	"time"

	"github.com/danivideda/satu-apotek-be/internal/env"
	"github.com/golang-jwt/jwt/v5"
)

var refreshTokenKey = env.GetString("JWT_REFRESH_SECRET", "refresh-token-secret")
var accessTokenKey = env.GetString("JWT_ACCESS_SECRET", "access-token-secret")

type myClaim struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func NewRefreshToken(id string) (string, error) {
	exp := time.Now().Add(24 * time.Hour)
	return generate(id, exp, refreshTokenKey)
}

func NewAccessToken(id string) (string, error) {
	exp := time.Now().Add(1 * time.Minute)
	return generate(id, exp, accessTokenKey)
}

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
