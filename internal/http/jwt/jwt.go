package jwt

import (
	"time"

	"github.com/danivideda/satu-apotek-be/internal/env"
	"github.com/golang-jwt/jwt/v5"
)

var (
	refreshTokenKey = env.GetString("JWT_REFRESH_SECRET", "refresh-token-secret")
	accessTokenKey  = env.GetString("JWT_ACCESS_SECRET", "access-token-secret")
	refreshTokenTTL = env.GetString("REFRESH_TTL", "168h")
	accessTokenTTL  = env.GetString("ACCESS_TTL", "5m")
)

type myClaim struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func NewRefreshToken(id string) (*time.Time, string, error) {
	ttl, err := time.ParseDuration(refreshTokenTTL)
	if err != nil {
		return nil, "", err
	}
	exp := time.Now().Add(ttl)
	signed, err := generate(id, exp, refreshTokenKey)
	return &exp, signed, err 
}

func NewAccessToken(id string) (string, error) {
	ttl, err := time.ParseDuration(accessTokenTTL)
	if err != nil {
		return "", err
	}
	exp := time.Now().Add(ttl)
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
	return t.SignedString([]byte(signKey))
}